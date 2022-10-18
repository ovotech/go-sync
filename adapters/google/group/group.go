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

const (
	maxResults = 200
)

// Ensure the Google Group adapter type fully satisfies the gosync.Adapter interface.
var _ gosync.Adapter = &Group{}

// InitKey is the required keys to Init a new adapter.
type InitKey = string

const (
	/*
		Authentication mechanism for Google. More info: https://cloud.google.com/docs/authentication

		Supported options are:
			`default`	Default Google credentials.
	*/
	Authentication InitKey = "google_authentication"
	Name           InitKey = "google_group_name" // Google Group name.
)

// callList allows us to mock the returned struct from the List Google API call.
func callList(ctx context.Context, call *admin.MembersListCall, pageToken string) (*admin.Members, error) {
	return call.Context(ctx).PageToken(pageToken).MaxResults(maxResults).Do() //nolint:wrapcheck
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
	logger         *log.Logger

	// Custom configuration for adding emails.
	deliverySettings string
	role             string

	callList   func(ctx context.Context, call *admin.MembersListCall, pageToken string) (*admin.Members, error)
	callInsert func(ctx context.Context, call *admin.MembersInsertCall) (*admin.Member, error)
	callDelete func(ctx context.Context, call *admin.MembersDeleteCall) error
}

// WithLogger sets a custom logger.
func WithLogger(logger *log.Logger) func(*Group) {
	return func(group *Group) {
		group.logger = logger
	}
}

// WithRole sets a custom role for new emails being added.
func WithRole(role string) func(*Group) {
	return func(group *Group) {
		group.role = role
	}
}

// WithDeliverySettings sets custom deliver settings when adding emails.
func WithDeliverySettings(deliverySettings string) func(*Group) {
	return func(group *Group) {
		group.deliverySettings = deliverySettings
	}
}

// New instantiates a new Google Group adapter.
func New(client *admin.Service, name string, optsFn ...func(*Group)) *Group {
	group := &Group{
		membersService: client.Members,
		name:           name,
		logger:         log.New(os.Stderr, "[go-sync/google/group] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix),

		callList:   callList,
		callInsert: callInsert,
		callDelete: callDelete,
	}

	for _, fn := range optsFn {
		fn(group)
	}

	return group
}

// Ensure the Init function fully satisfies the gosync.InitFn type.
var _ gosync.InitFn = Init

// Init a new Google Group gosync.Adapter. All InitKey keys are required in config.
func Init(config map[InitKey]string) (gosync.Adapter, error) {
	ctx := context.Background()

	for _, key := range []InitKey{Authentication, Name} {
		if _, ok := config[key]; !ok {
			return nil, fmt.Errorf("google.group.init -> %w(%s)", gosync.ErrMissingConfig, key)
		}
	}

	var (
		client *admin.Service
		err    error
	)

	scopes := option.WithScopes(admin.AdminDirectoryGroupMemberScope)

	switch config[Authentication] {
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
		return nil, fmt.Errorf("google.group.init -> %w(%s)", gosync.ErrInvalidConfig, config[Authentication])
	}

	return New(client, config[Name]), nil
}

// Get emails of Google users in a group.
func (g *Group) Get(ctx context.Context) ([]string, error) {
	var (
		pageToken = ""
		emails    = make([]string, 0)
	)

	for {
		g.logger.Printf("Fetching accounts from Google Group %s", g.name)

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

	g.logger.Println("Fetched accounts successfully")

	return emails, nil
}

// Add emails to a Google Group.
func (g *Group) Add(ctx context.Context, emails []string) error {
	g.logger.Printf("Adding %s to Google Group %s", emails, g.name)

	for _, email := range emails {
		_, err := g.callInsert(ctx, g.membersService.Insert(g.name, &admin.Member{
			Email:            email,
			DeliverySettings: g.deliverySettings,
			Role:             g.role,
		}))
		if err != nil {
			return fmt.Errorf("google.group.add(%s, %s) -> %w", g.name, email, err)
		}
	}

	g.logger.Println("Finished adding accounts successfully")

	return nil
}

// Remove emails from a Google Group.
func (g *Group) Remove(ctx context.Context, emails []string) error {
	g.logger.Printf("Removing %s from Google Group %s", emails, g.name)

	for _, email := range emails {
		err := g.callDelete(ctx, g.membersService.Delete(g.name, email))
		if err != nil {
			return fmt.Errorf("google.group.remove(%s, %s) -> %w", g.name, email, err)
		}
	}

	g.logger.Println("Finished removing accounts successfully")

	return nil
}
