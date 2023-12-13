/*
Package user synchronises email addresses from an Azure Active Directory
tenancy. It is a read-only adapter, so it does not provide any functionality to
create or delete entries in AD.

# Requirements

In order to synchronise with Azure AD, you'll need to ensure you have an
App Registration with the following API Permissions:
  - User.Read.All

You will also need credentials for the App Registration to be configured
appropriately in the environment using one of the default authentication
mechanisms.
*/
package user

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphsdkgocore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"

	gosync "github.com/ovotech/go-sync"
)

/*
Filter is an optional filter for filtering down the users to return
See the [reference] for more information on filter queries.

[reference]: https://learn.microsoft.com/en-us/graph/filter-query-parameter
*/
const Filter gosync.ConfigKey = "filter"

type iUser interface {
	Get(
		ctx context.Context,
		config *users.UsersRequestBuilderGetRequestConfiguration,
	) (models.UserCollectionResponseable, error)
}

type iClient interface {
	Users() *users.UsersRequestBuilder
	GetAdapter() abstractions.RequestAdapter
}

type User struct {
	users  iUser
	client iClient
	Logger *log.Logger

	filter string
}

func (u *User) Get(ctx context.Context) ([]string, error) {
	var request users.UsersRequestBuilderGetRequestConfiguration

	if isAdvancedQuery(u.filter) {
		headers := abstractions.NewRequestHeaders()
		headers.Add("ConsistencyLevel", "eventual")

		request = users.UsersRequestBuilderGetRequestConfiguration{
			Headers: headers,
			QueryParameters: &users.UsersRequestBuilderGetQueryParameters{
				Count:  to.Ptr(true),
				Filter: to.Ptr(u.filter),
			},
		}
	} else {
		request = users.UsersRequestBuilderGetRequestConfiguration{
			QueryParameters: &users.UsersRequestBuilderGetQueryParameters{
				Filter: to.Ptr(u.filter),
			},
		}
	}

	resp, err := u.users.Get(ctx, to.Ptr(request))
	if err != nil {
		return nil, fmt.Errorf("azuread.user.get.userget -> %w", err)
	}

	// Use PageIterator to iterate through all users
	pageIterator, err := msgraphsdkgocore.NewPageIterator[models.Userable](
		resp,
		u.client.GetAdapter(),
		models.CreateUserCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		return nil, fmt.Errorf("azuread.user.get.iterator -> %w", err)
	}

	emails := make([]string, 0)

	err = pageIterator.Iterate(ctx, func(user models.Userable) bool {
		if email := user.GetMail(); email != nil {
			emails = append(emails, *email)
		}

		return true
	})
	if err != nil {
		return nil, fmt.Errorf("azuread.user.get.iterate -> %w", err)
	}

	return emails, nil
}

//nolint:gochecknoglobals
var advancedQueries = []*regexp.Regexp{
	regexp.MustCompile(`(?i:\bendswith\b)`),
	regexp.MustCompile(`(?i:\bne\b)`),
	regexp.MustCompile(`(?i:\bnot\b)`),
	regexp.MustCompile(`\bcompanyName\b`),
	regexp.MustCompile(`\bemployeeOrgData/costCenter\b`),
	regexp.MustCompile(`\bemployeeOrgData/division\b`),
	regexp.MustCompile(`\bemployeeType\b`),
	regexp.MustCompile(`\bofficeLocation\b`),
	regexp.MustCompile(`\bonPremisesExtensionAttributes/extensionAttribute(?:1[0-5]?|[2-9])\b`),
}

func isAdvancedQuery(filter string) bool {
	for _, aq := range advancedQueries {
		if aq.MatchString(filter) {
			return true
		}
	}

	return false
}

func (u *User) Add(_ context.Context, _ []string) error {
	return fmt.Errorf("azuread.user.add -> %w", gosync.ErrReadOnly)
}

func (u *User) Remove(_ context.Context, _ []string) error {
	return fmt.Errorf("azuread.user.remove -> %w", gosync.ErrReadOnly)
}

var (
	_ gosync.Adapter       = &User{} // Ensure [user.User] fully satisfies the [gosync.Adapter] interface.
	_ gosync.InitFn[*User] = Init    // Ensure the [user.Init] function fully satisfies the [gosync.InitFn] type.
)

// WithFilter provides a mechanism to set the Microsoft Graph query filter
// when instantiating a new User adapter with the [user.New] method.
func WithFilter(f string) func(u *User) {
	return func(u *User) {
		u.filter = f
	}
}

// WithClient provides a mechanism to pass a custom client to the adapter.
func WithClient(client iClient) func(u *User) {
	return func(u *User) {
		u.client = client
	}
}

// Init creates a new Adapter. By default, an Azure Graph Service Client will
// be created using the default credentials in the environment.
func Init(
	_ context.Context,
	config map[gosync.ConfigKey]string,
	configFns ...gosync.ConfigFn[*User],
) (*User, error) {
	creds, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("azuread.user.init.creds -> %w", err)
	}

	client, err := msgraphsdkgo.NewGraphServiceClientWithCredentials(creds, []string{})
	if err != nil {
		return nil, fmt.Errorf("azuread.user.init.client -> %w", err)
	}

	user := &User{
		client: client,
		users:  client.Users(),
		Logger: log.New(os.Stderr, "[go-sync/azuread/user] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix),
	}

	for _, configFn := range configFns {
		configFn(user)
	}

	if filter, ok := config[Filter]; ok {
		user.filter = filter
	}

	return user, nil
}
