/*
Package group synchronises email addresses with a Google Group.

# Requirements

In order to synchronise with Google, you'll need to credentials with the Admin SDK enabled on your account, and
credentials with the following scopes:
  - [admin.AdminDirectoryGroupMemberScope]

# Examples

See [New] and [Init].
*/
package group

import (
	"context"
	"fmt"
	"log"
	"os"

	gosync "github.com/ovotech/go-sync"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

/*
GoogleAuthenticationMechanism sets the authentication mechanism for Google.

Supported options are:
  - [default]

[default]: https://cloud.google.com/docs/authentication/application-default-credentials
*/
const GoogleAuthenticationMechanism gosync.ConfigKey = "google_authentication_mechanism"

// Name is the name of your Google group.
const Name gosync.ConfigKey = "name"

/*
Role sets the role for new users being added to a group.

Acceptable values:
  - MANAGER
  - MEMBER
  - OWNER

See `role` field in [reference] for more information.

[reference]: https://developers.google.com/admin-sdk/directory/reference/rest/v1/members#resource:-member
*/
const Role gosync.ConfigKey = "role"

/*
DeliverySettings sets the delivery settings.

Acceptable values:
  - ALL_MAIL
  - DAILY
  - DIGEST
  - DISABLED
  - NONE

See `delivery_settings` field in [reference] for more information.

[reference]: https://developers.google.com/admin-sdk/directory/reference/rest/v1/members#resource:-member
*/
const DeliverySettings gosync.ConfigKey = "delivery_settings"

var (
	_ gosync.Adapter = &Group{} // Ensure [group.Group] fully satisfies the [gosync.Adapter] interface.
	_ gosync.InitFn  = Init     // Ensure the [group.Init] function fully satisfies the [gosync.InitFn] type.
)

// callList allows us to mock the returned struct from the List Google API call.
func callList(ctx context.Context, call *admin.MembersListCall, pageToken string) (*admin.Members, error) {
	return call.Context(ctx).PageToken(pageToken).MaxResults(200).Do() //nolint:wrapcheck,gomnd
}

// callInsert allows us to mock the returned struct from the Insert Google API call.
func callInsert(ctx context.Context, call *admin.MembersInsertCall) (*admin.Member, error) {
	return call.Context(ctx).Do() //nolint:wrapcheck
}

// callDelete allows us to mock the returned struct from the Delete Google API call.
func callDelete(ctx context.Context, call *admin.MembersDeleteCall) error {
	return call.Context(ctx).Do() //nolint:wrapcheck
}

// iMembersService is a subset of the Google MembersService, and used to build mocks for easy testing.
type iMembersService interface {
	List(groupKey string) *admin.MembersListCall
	Insert(groupKey string, member *admin.Member) *admin.MembersInsertCall
	Delete(groupKey string, memberKey string) *admin.MembersDeleteCall
}

type Group struct {
	membersService iMembersService
	name           string
	Logger         *log.Logger

	DeliverySettings string // See [group.DeliverySettings].
	Role             string // See [group.Role].

	callList   func(ctx context.Context, call *admin.MembersListCall, pageToken string) (*admin.Members, error)
	callInsert func(ctx context.Context, call *admin.MembersInsertCall) (*admin.Member, error)
	callDelete func(ctx context.Context, call *admin.MembersDeleteCall) error
}

// Get email addresses in a Google Group.
func (g *Group) Get(ctx context.Context) ([]string, error) {
	var (
		pageToken = ""
		emails    = make([]string, 0)
	)

	for {
		g.Logger.Printf("Fetching accounts from Google Group %s", g.name)

		response, err := g.callList(ctx, g.membersService.List(g.name), pageToken)
		if err != nil {
			return nil, fmt.Errorf("google.group.get(%s).list -> %w", g.name, err)
		}

		for _, member := range response.Members {
			emails = append(emails, member.Email)
		}

		pageToken = response.NextPageToken

		if pageToken == "" {
			break
		}
	}

	g.Logger.Println("Fetched accounts successfully")

	return emails, nil
}

// Add email addresses to a Google Group.
func (g *Group) Add(ctx context.Context, emails []string) error {
	g.Logger.Printf("Adding %s to Google Group %s", emails, g.name)

	for _, email := range emails {
		_, err := g.callInsert(ctx, g.membersService.Insert(g.name, &admin.Member{
			Email:            email,
			DeliverySettings: g.DeliverySettings,
			Role:             g.Role,
		}))
		if err != nil {
			return fmt.Errorf("google.group.add(%s, %s) -> %w", g.name, email, err)
		}
	}

	g.Logger.Println("Finished adding accounts successfully")

	return nil
}

// Remove email addresses from a Google Group.
func (g *Group) Remove(ctx context.Context, emails []string) error {
	g.Logger.Printf("Removing %s from Google Group %s", emails, g.name)

	for _, email := range emails {
		err := g.callDelete(ctx, g.membersService.Delete(g.name, email))
		if err != nil {
			return fmt.Errorf("google.group.remove(%s, %s) -> %w", g.name, email, err)
		}
	}

	g.Logger.Println("Finished removing accounts successfully")

	return nil
}

/*
New Google Group [gosync.Adapter].

Recommended reading for parameters:
  - name: [group.Name]
*/
func New(client *admin.Service, name string, optsFn ...func(*Group)) *Group {
	group := &Group{
		membersService: client.Members,
		name:           name,
		Logger:         log.New(os.Stderr, "[go-sync/google/group] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix),

		callList:   callList,
		callInsert: callInsert,
		callDelete: callDelete,
	}

	for _, fn := range optsFn {
		fn(group)
	}

	return group
}

/*
Init a new Google Group [gosync.Adapter].

Required config:
  - [group.GoogleAuthenticationMechanism]
  - [group.Name]
*/
func Init(ctx context.Context, config map[gosync.ConfigKey]string) (gosync.Adapter, error) {
	for _, key := range []gosync.ConfigKey{GoogleAuthenticationMechanism, Name} {
		if _, ok := config[key]; !ok {
			return nil, fmt.Errorf("google.group.init -> %w(%s)", gosync.ErrMissingConfig, key)
		}
	}

	var (
		client *admin.Service
		err    error
	)

	scopes := option.WithScopes(admin.AdminDirectoryGroupMemberScope)

	switch config[GoogleAuthenticationMechanism] {
	case "_testing_":
		// Only for use in testing in order to prevent failure to fetch default credentials.
		client, err = admin.NewService(ctx, scopes, option.WithAPIKey("_testing_"))
		if err != nil {
			return nil, fmt.Errorf("google.group.init -> %w", err)
		}
	case "default":
		client, err = admin.NewService(ctx, scopes)
		if err != nil {
			return nil, fmt.Errorf("google.group.init -> %w", err)
		}
	default:
		return nil, fmt.Errorf("google.group.init -> %w(%s)", gosync.ErrInvalidConfig, config[GoogleAuthenticationMechanism])
	}

	adapter := New(client, config[Name])

	if val, ok := config[Role]; ok {
		adapter.Role = val
	}

	if val, ok := config[DeliverySettings]; ok {
		adapter.DeliverySettings = val
	}

	return adapter, nil
}
