/*
Package groupmembership synchronises email addresses to groups within an Azure AD
tenancy.

# Requirements

In order to synchronise group membership within Azure AD, you'll need to ensure
you have an App Registration with the following API Permissions:
  - GroupMember.ReadWrite.All
  - User.ReadWrite.All

You will also need credentials for the App Registration to be configured
appropriately in the environment using one of the default authentication
mechanisms.
*/
package groupmembership

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphsdkgocore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/groups"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	"github.com/microsoftgraph/msgraph-sdk-go/users"

	gosync "github.com/ovotech/go-sync"
)

// GroupName is the name of your group within Azure AD.
const GroupName gosync.ConfigKey = "group_name"

type iClient interface {
	GetAdapter() abstractions.RequestAdapter
}

var _ iClient = &msgraphsdkgo.GraphServiceClient{}

type iGroupClient interface {
	Get(
		ctx context.Context,
		config *groups.GroupsRequestBuilderGetRequestConfiguration,
	) (models.GroupCollectionResponseable, error)
	ByGroupId(id string) *groups.GroupItemRequestBuilder
}

var _ iGroupClient = &groups.GroupsRequestBuilder{}

type iUserClient interface {
	Get(
		ctx context.Context,
		config *users.UsersRequestBuilderGetRequestConfiguration,
	) (models.UserCollectionResponseable, error)
}

var _ iUserClient = &users.UsersRequestBuilder{}

type GroupMembership struct {
	client      iClient
	groupClient iGroupClient
	userClient  iUserClient

	Logger *log.Logger

	group string

	getGroupMembers func(
		context.Context,
		*groups.GroupItemRequestBuilder,
		*groups.ItemMembersRequestBuilderGetRequestConfiguration,
	) (models.DirectoryObjectCollectionResponseable, error)
	patchGroup func(
		context.Context,
		*groups.GroupItemRequestBuilder,
		models.Groupable,
		*groups.GroupItemRequestBuilderPatchRequestConfiguration,
	) (models.Groupable, error)
	removeGroupMember func(
		context.Context,
		*groups.GroupItemRequestBuilder,
		string,
		*groups.ItemMembersItemRefRequestBuilderDeleteRequestConfiguration,
	) error
}

// Get will return the membership of the group.
func (g *GroupMembership) Get(ctx context.Context) ([]string, error) {
	gid, err := resolveGroupID(ctx, g.groupClient, g.group)
	if err != nil {
		return nil, fmt.Errorf("azuread.groupmembership(%s).get -> %w", g.group, resolveOdataError(err))
	}

	req := groups.ItemMembersRequestBuilderGetRequestConfiguration{
		QueryParameters: &groups.ItemMembersRequestBuilderGetQueryParameters{
			Select: []string{"mail"},
		},
	}

	resp, err := g.getGroupMembers(ctx, g.groupClient.ByGroupId(gid), to.Ptr(req))
	// resp, err := g.groupClient.ByGroupId(gid).Members().Get(ctx, to.Ptr(req))
	if err != nil {
		return nil, fmt.Errorf("azuread.groupmembership(%s).get -> %w", g.group, resolveOdataError(err))
	}

	pageIterator, err := msgraphsdkgocore.NewPageIterator[models.Userable](
		resp,
		g.client.GetAdapter(),
		models.CreateGroupCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		return nil, fmt.Errorf("azuread.groupmembership(%s).get -> %w", g.group, resolveOdataError(err))
	}

	emails := make([]string, 0)

	err = pageIterator.Iterate(ctx, func(user models.Userable) bool {
		if *user.GetOdataType() != "#microsoft.graph.user" {
			return true
		}

		if email := user.GetMail(); email != nil {
			emails = append(emails, *email)
		}

		return true
	})
	if err != nil {
		return nil, fmt.Errorf("azuread.groupmembership(%s).get -> %w", g.group, resolveOdataError(err))
	}

	return emails, nil
}

// The graph API supports adding a maximum of 20 members to a group in a single
// operation.
// https://learn.microsoft.com/en-gb/graph/api/group-post-members
const addBatchSize = 20

// Add will add the given members to the group's member list.
func (g *GroupMembership) Add(ctx context.Context, members []string) error {
	gid, err := resolveGroupID(ctx, g.groupClient, g.group)
	if err != nil {
		return fmt.Errorf(
			"azuread.groupmembership(%s).add.resolveGroupID -> %w",
			g.group,
			resolveOdataError(err),
		)
	}

	var payloads [][]string

	end := 0
	for start := 0; start < len(members); start += addBatchSize {
		end += addBatchSize
		if end > len(members) {
			end = len(members)
		}

		payload := make([]string, 0, end-start)

		for _, member := range members[start:end] {
			uid, err := resolveUserID(ctx, g.userClient, member)
			if err != nil {
				return fmt.Errorf(
					"azuread.groupmembership(%s).add(%s).resolveUserID -> %w",
					g.group,
					member,
					resolveOdataError(err),
				)
			}

			payload = append(payload, "https://graph.microsoft.com/v1.0/directoryObjects/"+uid)
		}

		payloads = append(payloads, payload)
	}

	for _, payload := range payloads {
		req := models.NewGroup()
		req.SetAdditionalData(map[string]interface{}{
			"members@odata.bind": payload,
		})

		_, err := g.patchGroup(ctx, g.groupClient.ByGroupId(gid), req, nil)
		// _, err := g.groupClient.ByGroupId(gid).Patch(ctx, req, nil)
		if err != nil {
			return fmt.Errorf("azuread.groupmembership(%s).add.patch -> %w", g.group, resolveOdataError(err))
		}
	}

	return nil
}

// Remove will remove the given email addresses from the group's member list.
func (g *GroupMembership) Remove(ctx context.Context, members []string) error {
	gid, err := resolveGroupID(ctx, g.groupClient, g.group)
	if err != nil {
		return fmt.Errorf(
			"azuread.groupmembership(%s).remove.resolveGroupID -> %w",
			g.group,
			resolveOdataError(err),
		)
	}

	for _, member := range members {
		uid, err := resolveUserID(ctx, g.userClient, member)
		if err != nil {
			return fmt.Errorf(
				"azuread.groupmembership(%s).remove(%s).resolveUserID -> %w",
				g.group,
				member,
				resolveOdataError(err),
			)
		}

		err = g.removeGroupMember(ctx, g.groupClient.ByGroupId(gid), uid, nil)
		// err = g.groupClient.ByGroupId(gid).Members().ByDirectoryObjectId(uid).Ref().Delete(ctx, nil)
		if err != nil {
			return fmt.Errorf(
				"azuread.groupmembership(%s).remove(%s).delete -> %w",
				g.group,
				member,
				resolveOdataError(err),
			)
		}
	}

	return nil
}

// ErrTooManyResults is returned when the query resulted in more results than
// expected.
var ErrTooManyResults = errors.New("too many results in response")

// ErrNoResults is returned when the query resulted in no results and some
// were expected.
var ErrNoResults = errors.New("no results found")

func resolveGroupID(ctx context.Context, client iGroupClient, group string) (string, error) {
	req := groups.GroupsRequestBuilderGetRequestConfiguration{
		QueryParameters: to.Ptr(groups.GroupsRequestBuilderGetQueryParameters{
			Select: []string{"id"},
			Filter: to.Ptr(fmt.Sprintf("displayName eq '%s'", group)),
		}),
	}

	resp, err := client.Get(ctx, to.Ptr(req))
	if err != nil {
		return "", fmt.Errorf("azuread.groupmembership.resolvegroupid -> %w", resolveOdataError(err))
	}

	count := len(resp.GetValue())
	if count == 0 {
		return "", fmt.Errorf("azuread.groupmembership.resolvegroupid -> %w", ErrNoResults)
	} else if count > 1 {
		return "", fmt.Errorf("azuread.groupmembership.resolvegroupid -> %w", ErrTooManyResults)
	}

	return *resp.GetValue()[0].GetId(), nil
}

func resolveUserID(ctx context.Context, client iUserClient, mail string) (string, error) {
	req := users.UsersRequestBuilderGetRequestConfiguration{
		QueryParameters: to.Ptr(users.UsersRequestBuilderGetQueryParameters{
			Select: []string{"id"},
			Filter: to.Ptr(fmt.Sprintf("mail eq '%s'", mail)),
		}),
	}

	resp, err := client.Get(ctx, to.Ptr(req))
	if err != nil {
		return "", fmt.Errorf("azuread.groupmembership.resolveuserid(%s) -> %w", mail, resolveOdataError(err))
	}

	count := len(resp.GetValue())
	if count == 0 {
		return "", fmt.Errorf("azuread.groupmembership.resolveuserid(%s) -> %w", mail, ErrNoResults)
	} else if count > 1 {
		return "", fmt.Errorf("azuread.groupmembership.resolveuserid(%s) -> %w", mail, ErrTooManyResults)
	}

	return *resp.GetValue()[0].GetId(), nil
}

func resolveOdataError(err error) error {
	var odataError *odataerrors.ODataError
	if errors.As(err, &odataError) {
		return fmt.Errorf(
			"%w -> (%s) %s",
			err,
			*odataError.GetErrorEscaped().GetCode(),
			*odataError.GetErrorEscaped().GetMessage(),
		)
	}

	return err
}

func getGroupMembers( //nolint:ireturn
	ctx context.Context,
	builder *groups.GroupItemRequestBuilder,
	cfg *groups.ItemMembersRequestBuilderGetRequestConfiguration,
) (models.DirectoryObjectCollectionResponseable, error) {
	out, err := builder.Members().Get(ctx, cfg)
	if err != nil {
		return out, fmt.Errorf("azuread.groupmembership.getGroupMembers -> %w", err)
	}

	return out, nil
}

func patchGroup( //nolint:ireturn
	ctx context.Context,
	builder *groups.GroupItemRequestBuilder,
	group models.Groupable,
	cfg *groups.GroupItemRequestBuilderPatchRequestConfiguration,
) (models.Groupable, error) {
	out, err := builder.Patch(ctx, group, cfg)
	if err != nil {
		return out, fmt.Errorf("azuread.groupmembership.patchGroup -> %w", err)
	}

	return out, nil
}

func removeGroupMember(
	ctx context.Context,
	builder *groups.GroupItemRequestBuilder,
	uid string,
	cfg *groups.ItemMembersItemRefRequestBuilderDeleteRequestConfiguration,
) error {
	if err := builder.Members().ByDirectoryObjectId(uid).Ref().Delete(ctx, cfg); err != nil {
		return fmt.Errorf("azuread.groupmembership.removeGroupMember(%s) -> %w", uid, err)
	}

	return nil
}

var (
	_ gosync.Adapter                  = &GroupMembership{}
	_ gosync.InitFn[*GroupMembership] = Init
)

// WithClient provides a mechanism to pass a custom client to the adapter.
func WithClient(client iClient) func(u *GroupMembership) {
	return func(u *GroupMembership) {
		u.client = client
	}
}

// Init creates a new adapter. It expects a single configuration entry.
// Required config:
//   - groupmembership.GroupName: the name of the AD group to sync members to.
func Init(
	_ context.Context,
	config map[gosync.ConfigKey]string,
	configFns ...gosync.ConfigFn[*GroupMembership],
) (*GroupMembership, error) {
	for _, key := range []gosync.ConfigKey{GroupName} {
		if _, ok := config[key]; !ok {
			return nil, fmt.Errorf("azuread.groupmembership.init -> %w(%s)", gosync.ErrMissingConfig, key)
		}
	}

	creds, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("azuread.groupmembership.init.creds -> %w", err)
	}

	client, err := msgraphsdkgo.NewGraphServiceClientWithCredentials(creds, []string{})
	if err != nil {
		return nil, fmt.Errorf("azuread.groupmembership.init.client -> %w", err)
	}

	adapter := &GroupMembership{
		client:      client,
		groupClient: client.Groups(),
		userClient:  client.Users(),

		group: config[GroupName],

		getGroupMembers:   getGroupMembers,
		patchGroup:        patchGroup,
		removeGroupMember: removeGroupMember,

		Logger: log.New(
			os.Stderr,
			"[go-sync/azuread/groupmembership] ",
			log.LstdFlags|log.Lshortfile|log.Lmsgprefix,
		),
	}

	for _, configFn := range configFns {
		configFn(adapter)
	}

	return adapter, nil
}
