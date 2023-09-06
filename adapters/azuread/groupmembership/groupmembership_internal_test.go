package groupmembership

import (
	"context"
	"fmt"
	"io"
	"log"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoft/kiota-abstractions-go/store"
	"github.com/microsoftgraph/msgraph-sdk-go/groups"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	gosync "github.com/ovotech/go-sync"
)

func TestInit(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{
			GroupName: "example",
		})

		assert.NoError(t, err)
		assert.IsType(t, &GroupMembership{}, adapter)
		assert.Equal(t, "example", adapter.(*GroupMembership).group)
	})

	t.Run("missing config", func(t *testing.T) {
		t.Parallel()

		t.Run("missing group name", func(t *testing.T) {
			t.Parallel()

			_, err := Init(ctx, map[gosync.ConfigKey]string{})

			assert.ErrorIs(t, err, gosync.ErrMissingConfig)
		})
	})
}

//nolint:dupl
func Test_resolveGroupID(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	group := "example"

	mockClientWithGroupIDResponse := func(ids ...string) iGroupClient {
		coll := make([]models.Groupable, 0, len(ids))

		for _, id := range ids {
			id := id
			g := models.NewGroup()
			g.SetId(&id)
			coll = append(coll, g)
		}

		mockResp := models.NewGroupCollectionResponse()
		mockResp.SetValue(coll)

		client := newMockIGroupClient(t)
		client.On("Get", ctx, mock.MatchedBy(
			func(req *groups.GroupsRequestBuilderGetRequestConfiguration) bool {
				return *req.QueryParameters.Filter == fmt.Sprintf("displayName eq '%s'", group)
			},
		)).Return(mockResp, nil)

		return client
	}

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		client := mockClientWithGroupIDResponse()

		gid, err := resolveGroupID(ctx, client, group)
		assert.Empty(t, gid)
		assert.ErrorIs(t, err, ErrNoResults)
	})

	t.Run("one found", func(t *testing.T) {
		t.Parallel()

		client := mockClientWithGroupIDResponse("00000001-0000-0000-0000-000123456789")

		gid, err := resolveGroupID(ctx, client, group)
		assert.Contains(t, gid, "00000001-0000-0000-0000-000123456789")
		assert.NoError(t, err)
	})

	t.Run("multiple found", func(t *testing.T) {
		t.Parallel()

		client := mockClientWithGroupIDResponse(
			"00000001-0000-0000-0000-000123456789",
			"00000002-0000-0000-0000-000123456789",
		)

		gid, err := resolveGroupID(ctx, client, group)
		assert.Empty(t, gid)
		assert.ErrorIs(t, err, ErrTooManyResults)
	})
}

//nolint:dupl
func Test_resolveUserID(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	email := "test.user@example.com"

	mockClientWithUserIDResponse := func(ids ...string) iUserClient {
		coll := make([]models.Userable, 0, len(ids))

		for _, id := range ids {
			id := id
			u := models.NewUser()
			u.SetId(&id)
			coll = append(coll, u)
		}

		mockResp := models.NewUserCollectionResponse()
		mockResp.SetValue(coll)

		client := newMockIUserClient(t)
		client.On("Get", ctx, mock.MatchedBy(
			func(req *users.UsersRequestBuilderGetRequestConfiguration) bool {
				return *req.QueryParameters.Filter == fmt.Sprintf("mail eq '%s'", email)
			},
		)).Return(mockResp, nil)

		return client
	}

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		client := mockClientWithUserIDResponse()

		gid, err := resolveUserID(ctx, client, email)
		assert.Empty(t, gid)
		assert.ErrorIs(t, err, ErrNoResults)
	})

	t.Run("one found", func(t *testing.T) {
		t.Parallel()

		client := mockClientWithUserIDResponse("00000001-1000-0000-0000-000123456789")

		gid, err := resolveUserID(ctx, client, email)
		assert.Contains(t, gid, "00000001-1000-0000-0000-000123456789")
		assert.NoError(t, err)
	})

	t.Run("multiple found", func(t *testing.T) {
		t.Parallel()

		client := mockClientWithUserIDResponse("00000001-1000-0000-0000-000123456789", "00000002-1000-0000-0000-000123456789")

		gid, err := resolveUserID(ctx, client, email)
		assert.Empty(t, gid)
		assert.ErrorIs(t, err, ErrTooManyResults)
	})
}

func TestGroupMembership_Get(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	expected := []string{"test.user1@example.com", "test.user2@example.com"}

	mockGetGroupMembers := func(
		_ context.Context,
		_ *groups.GroupItemRequestBuilder,
		_ *groups.ItemMembersRequestBuilderGetRequestConfiguration,
	) (models.DirectoryObjectCollectionResponseable, error) {
		coll := make([]models.DirectoryObjectable, 0, len(expected))

		for _, email := range expected {
			u := models.NewUser()
			u.SetMail(to.Ptr(email))
			coll = append(coll, u)
		}

		resp := models.NewDirectoryObjectCollectionResponse()
		resp.SetValue(coll)

		return resp, nil
	}

	mockResp := models.NewGroupCollectionResponse()
	mockResp.SetValue([]models.Groupable{models.NewGroup()})
	mockResp.GetValue()[0].SetId(to.Ptr("00000001-0000-0000-0000-000123456789"))

	groupClient := newMockIGroupClient(t)
	groupClient.On("Get", ctx, mock.MatchedBy(
		func(req *groups.GroupsRequestBuilderGetRequestConfiguration) bool {
			return *req.QueryParameters.Filter == "displayName eq 'example'"
		},
	)).Return(mockResp, nil)
	groupClient.On(
		"ByGroupId",
		"00000001-0000-0000-0000-000123456789",
	).Return(&groups.GroupItemRequestBuilder{
		BaseRequestBuilder: abstractions.BaseRequestBuilder{RequestAdapter: &MockRequestAdapter{}},
	})

	client := newMockIClient(t)
	client.On("GetAdapter").Return(&MockRequestAdapter{})

	adapter := &GroupMembership{
		Logger:          log.New(io.Discard, "", 0),
		client:          client,
		groupClient:     groupClient,
		group:           "example",
		getGroupMembers: mockGetGroupMembers,
	}

	out, err := adapter.Get(ctx)
	assert.NoError(t, err)
	assert.ElementsMatch(t, out, expected)
}

//nolint:funlen
func TestGroupMembership_Add(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	gid1 := "00000001-0000-0000-0000-000123456789"
	mockResp := models.NewGroupCollectionResponse()
	mockResp.SetValue([]models.Groupable{models.NewGroup()})
	mockResp.GetValue()[0].SetId(to.Ptr(gid1))

	groupName := "TestGroupMembership_Add"

	groupClient := newMockIGroupClient(t)
	groupClient.On("Get", ctx, mock.MatchedBy(
		func(req *groups.GroupsRequestBuilderGetRequestConfiguration) bool {
			return *req.QueryParameters.Filter == fmt.Sprintf("displayName eq '%s'", groupName)
		},
	)).Return(mockResp, nil)
	groupClient.On(
		"ByGroupId",
		gid1,
	).Return(&groups.GroupItemRequestBuilder{
		BaseRequestBuilder: abstractions.BaseRequestBuilder{RequestAdapter: &MockRequestAdapter{}},
	})

	uid1 := "00000001-1000-0000-0000-000123456789"
	userResp1 := models.NewUserCollectionResponse()
	userResp1.SetValue([]models.Userable{models.NewUser()})
	userResp1.GetValue()[0].SetId(to.Ptr(uid1))

	uid2 := "00000002-1000-0000-0000-000123456789"
	userResp2 := models.NewUserCollectionResponse()
	userResp2.SetValue([]models.Userable{models.NewUser()})
	userResp2.GetValue()[0].SetId(to.Ptr(uid2))

	userMail1 := "test.user1@example.com"
	userMail2 := "test.user2@example.com"

	userClient := newMockIUserClient(t)
	userClient.On("Get", ctx, mock.MatchedBy(
		func(req *users.UsersRequestBuilderGetRequestConfiguration) bool {
			return *req.QueryParameters.Filter == fmt.Sprintf("mail eq '%s'", userMail1)
		},
	)).Return(userResp1, nil)
	userClient.On("Get", ctx, mock.MatchedBy(
		func(req *users.UsersRequestBuilderGetRequestConfiguration) bool {
			return *req.QueryParameters.Filter == fmt.Sprintf("mail eq '%s'", userMail2)
		},
	)).Return(userResp2, nil)

	adapter := &GroupMembership{
		Logger:      log.New(io.Discard, "", 0),
		groupClient: groupClient,
		userClient:  userClient,
		group:       groupName,
		patchGroup: func(
			_ context.Context,
			_ *groups.GroupItemRequestBuilder,
			req models.Groupable,
			_ *groups.GroupItemRequestBuilderPatchRequestConfiguration,
		) (models.Groupable, error) {
			assert.Contains(
				t,
				req.GetAdditionalData()["members@odata.bind"],
				fmt.Sprintf("https://graph.microsoft.com/v1.0/directoryObjects/%s", uid1),
			)
			assert.Contains(
				t,
				req.GetAdditionalData()["members@odata.bind"],
				fmt.Sprintf("https://graph.microsoft.com/v1.0/directoryObjects/%s", uid2),
			)

			return req, nil
		},
	}

	expected := []string{userMail1, userMail2}
	err := adapter.Add(ctx, expected)
	assert.NoError(t, err)
}

//nolint:funlen
func TestGroupMembership_Remove(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	gid1 := "00000001-0000-0000-0000-000123456789"
	mockResp := models.NewGroupCollectionResponse()
	mockResp.SetValue([]models.Groupable{models.NewGroup()})
	mockResp.GetValue()[0].SetId(to.Ptr(gid1))

	groupName := "TestGroupMembership_Remove"

	groupClient := newMockIGroupClient(t)
	groupClient.On("Get", ctx, mock.MatchedBy(
		func(req *groups.GroupsRequestBuilderGetRequestConfiguration) bool {
			return *req.QueryParameters.Filter == fmt.Sprintf("displayName eq '%s'", groupName)
		},
	)).Return(mockResp, nil)
	groupClient.On(
		"ByGroupId",
		gid1,
	).Return(&groups.GroupItemRequestBuilder{
		BaseRequestBuilder: abstractions.BaseRequestBuilder{RequestAdapter: &MockRequestAdapter{}},
	})

	uid1 := "00000001-1000-0000-0000-000123456789"
	userResp1 := models.NewUserCollectionResponse()
	userResp1.SetValue([]models.Userable{models.NewUser()})
	userResp1.GetValue()[0].SetId(to.Ptr(uid1))

	uid2 := "00000002-1000-0000-0000-000123456789"
	userResp2 := models.NewUserCollectionResponse()
	userResp2.SetValue([]models.Userable{models.NewUser()})
	userResp2.GetValue()[0].SetId(to.Ptr(uid2))

	userMail1 := "test.user1@example.com"
	userMail2 := "test.user2@example.com"

	userClient := newMockIUserClient(t)
	userClient.On("Get", ctx, mock.MatchedBy(
		func(req *users.UsersRequestBuilderGetRequestConfiguration) bool {
			return *req.QueryParameters.Filter == fmt.Sprintf("mail eq '%s'", userMail1)
		},
	)).Return(userResp1, nil)
	userClient.On("Get", ctx, mock.MatchedBy(
		func(req *users.UsersRequestBuilderGetRequestConfiguration) bool {
			return *req.QueryParameters.Filter == fmt.Sprintf("mail eq '%s'", userMail2)
		},
	)).Return(userResp2, nil)

	adapter := &GroupMembership{
		Logger:      log.New(io.Discard, "", 0),
		groupClient: groupClient,
		userClient:  userClient,
		group:       groupName,
		removeGroupMember: func(
			ctx context.Context,
			builder *groups.GroupItemRequestBuilder,
			uid string,
			configuration *groups.ItemMembersItemRefRequestBuilderDeleteRequestConfiguration,
		) error {
			assert.Contains(t, []string{uid1, uid2}, uid)

			return nil
		},
	}

	err := adapter.Remove(ctx, []string{userMail1, userMail2})
	assert.NoError(t, err)
}

type MockRequestAdapter struct {
	SerializationWriterFactory serialization.SerializationWriterFactory
}

func (r *MockRequestAdapter) Send(
	_ context.Context,
	_ *abstractions.RequestInformation,
	_ serialization.ParsableFactory,
	_ abstractions.ErrorMappings,
) (serialization.Parsable, error) {
	return nil, nil
}

func (r *MockRequestAdapter) SendEnum(
	_ context.Context,
	_ *abstractions.RequestInformation,
	_ serialization.EnumFactory,
	_ abstractions.ErrorMappings,
) (any, error) {
	return nil, nil //nolint:nilnil
}

func (r *MockRequestAdapter) SendCollection(
	_ context.Context,
	_ *abstractions.RequestInformation,
	_ serialization.ParsableFactory,
	_ abstractions.ErrorMappings,
) ([]serialization.Parsable, error) {
	return nil, nil
}

func (r *MockRequestAdapter) SendEnumCollection(
	_ context.Context,
	_ *abstractions.RequestInformation,
	_ serialization.EnumFactory,
	_ abstractions.ErrorMappings,
) ([]any, error) {
	return nil, nil
}

func (r *MockRequestAdapter) SendPrimitive(
	_ context.Context,
	_ *abstractions.RequestInformation,
	_ string,
	_ abstractions.ErrorMappings,
) (any, error) {
	return nil, nil //nolint:nilnil
}

func (r *MockRequestAdapter) SendPrimitiveCollection(
	_ context.Context,
	_ *abstractions.RequestInformation,
	_ string,
	_ abstractions.ErrorMappings,
) ([]any, error) {
	return nil, nil
}

func (r *MockRequestAdapter) SendNoContent(
	_ context.Context,
	_ *abstractions.RequestInformation,
	_ abstractions.ErrorMappings,
) error {
	return nil
}

func (r *MockRequestAdapter) ConvertToNativeRequest(
	_ context.Context,
	_ *abstractions.RequestInformation,
) (any, error) {
	return nil, nil //nolint:nilnil
}

func (r *MockRequestAdapter) GetSerializationWriterFactory() serialization.SerializationWriterFactory {
	return r.SerializationWriterFactory
}

func (r *MockRequestAdapter) EnableBackingStore(_ store.BackingStoreFactory) {
}

//nolint:revive,stylecheck
func (r *MockRequestAdapter) SetBaseUrl(_ string) {
}

//nolint:revive,stylecheck
func (r *MockRequestAdapter) GetBaseUrl() string {
	return ""
}
