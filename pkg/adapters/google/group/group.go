package group

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ovotech/go-sync/internal/types"
	"github.com/ovotech/go-sync/pkg/ports"
	admin "google.golang.org/api/admin/directory/v1"
)

// Ensure the adapter type fully satisfies the ports.Adapter interface.
var _ ports.Adapter = &Group{}

// callList allows us to mock the returned struct from the List Google API call.
func callList(ctx context.Context, call *admin.MembersListCall, pageToken string) (*admin.Members, error) {
	return call.Context(ctx).PageToken(pageToken).Do() //nolint:wrapcheck
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
	logger         types.Logger

	listCall   func(ctx context.Context, call *admin.MembersListCall, pageToken string) (*admin.Members, error)
	insertCall func(ctx context.Context, call *admin.MembersInsertCall) (*admin.Member, error)
	deleteCall func(ctx context.Context, call *admin.MembersDeleteCall) error
}

// OptionLogger can be used to set a custom logger.
func OptionLogger(logger types.Logger) func(*Group) {
	return func(group *Group) {
		group.logger = logger
	}
}

// New instantiates a new Google Group adapter.
func New(client *admin.Service, name string, optsFn ...func(*Group)) *Group {
	group := &Group{
		membersService: client.Members,
		name:           name,
		logger:         log.New(os.Stderr, "[go-sync/google/group] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix),

		listCall:   callList,
		insertCall: callInsert,
		deleteCall: callDelete,
	}

	for _, fn := range optsFn {
		fn(group)
	}

	return group
}

// Get emails of Google users in a group.
func (g *Group) Get(ctx context.Context) ([]string, error) {
	var (
		pageToken = ""
		emails    = make([]string, 0)
	)

	for {
		g.logger.Printf("Fetching accounts from Google Group %s", g.name)

		response, err := g.listCall(ctx, g.membersService.List(g.name), pageToken)
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
		_, err := g.insertCall(ctx, g.membersService.Insert(g.name, &admin.Member{Email: email}))
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
		err := g.deleteCall(ctx, g.membersService.Delete(g.name, email))
		if err != nil {
			return fmt.Errorf("google.group.remove(%s, %s) -> %w", g.name, email, err)
		}
	}

	g.logger.Println("Finished removing accounts successfully")

	return nil
}
