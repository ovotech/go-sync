package group

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"

	gosync "github.com/ovotech/go-sync"
)

func withMockAdminService(ctx context.Context, t *testing.T) gosync.ConfigFn[*Group] {
	t.Helper()

	client, err := admin.NewService(
		ctx,
		option.WithScopes(admin.AdminDirectoryGroupMemberScope),
		option.WithAPIKey("_testing_"),
	)
	assert.NoError(t, err)

	return func(g *Group) {
		g.membersService = client.Members
	}
}

type mockCalls struct {
	mock.Mock
}

func (m *mockCalls) callList(
	ctx context.Context,
	call *admin.MembersListCall,
	pageToken string,
) (*admin.Members, error) {
	args := m.Called(ctx, call, pageToken)

	return args.Get(0).(*admin.Members), args.Error(1) //nolint:wrapcheck
}

func (m *mockCalls) callInsert(ctx context.Context, call *admin.MembersInsertCall) (*admin.Member, error) {
	args := m.Called(ctx, call)

	return args.Get(0).(*admin.Member), args.Error(1) //nolint:wrapcheck
}

func (m *mockCalls) callDelete(ctx context.Context, call *admin.MembersDeleteCall) error {
	args := m.Called(ctx, call)

	return args.Error(0) //nolint:wrapcheck
}

func TestGroups_Get(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	mockMembersService := newMockIMembersService(t)
	mockMembersService.EXPECT().List("test").Return(nil)

	mockCall := new(mockCalls)
	mockCall.On("callList", ctx, mock.Anything, "").Return(&admin.Members{
		NextPageToken: "page-2",
		Members: []*admin.Member{
			{Email: "foo@email"},
		},
	}, nil)
	mockCall.On("callList", ctx, mock.Anything, "page-2").Return(&admin.Members{
		Members: []*admin.Member{
			{Email: "bar@email"},
		},
	}, nil)

	group := &Group{
		name:           "test",
		membersService: mockMembersService,
		Logger:         log.New(os.Stdout, "", log.LstdFlags),
		callList:       mockCall.callList,
		callInsert:     mockCall.callInsert,
		callDelete:     mockCall.callDelete,
	}

	emails, err := group.Get(ctx)

	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"foo@email", "bar@email"}, emails)
}

func TestGroups_Add(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	mockMembersService := newMockIMembersService(t)
	mockMembersService.EXPECT().Insert("test", &admin.Member{Email: "foo@email"}).Return(nil)
	mockMembersService.EXPECT().Insert("test", &admin.Member{Email: "bar@email"}).Return(nil)

	mockCall := new(mockCalls)
	mockCall.On("callInsert", ctx, mock.Anything).Return(&admin.Member{}, nil)

	group := &Group{
		name:           "test",
		membersService: mockMembersService,
		Logger:         log.New(os.Stdout, "", log.LstdFlags),
		callList:       mockCall.callList,
		callInsert:     mockCall.callInsert,
		callDelete:     mockCall.callDelete,
	}

	err := group.Add(ctx, []string{"foo@email", "bar@email"})

	assert.NoError(t, err)
}

func TestGroups_Remove(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	mockMembersService := newMockIMembersService(t)
	mockMembersService.EXPECT().Delete("test", "foo@email").Return(nil)
	mockMembersService.EXPECT().Delete("test", "bar@email").Return(nil)

	mockCall := new(mockCalls)
	mockCall.On("callDelete", ctx, mock.Anything).Twice().Return(nil)

	group := &Group{
		name:           "test",
		membersService: mockMembersService,
		Logger:         log.New(os.Stdout, "", log.LstdFlags),
		callList:       mockCall.callList,
		callInsert:     mockCall.callInsert,
		callDelete:     mockCall.callDelete,
	}

	err := group.Remove(ctx, []string{"foo@email", "bar@email"})

	assert.NoError(t, err)
}

func TestRole(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	mockMembersService := newMockIMembersService(t)
	mockMembersService.EXPECT().Insert("test", &admin.Member{
		Email: "foo@email",
		Role:  "test-role",
	}).Return(nil)

	mockCall := new(mockCalls)
	mockCall.On("callInsert", ctx, mock.Anything).Return(&admin.Member{}, nil)

	group := &Group{
		name:           "test",
		Role:           "test-role",
		membersService: mockMembersService,
		Logger:         log.New(os.Stdout, "", log.LstdFlags),
		callList:       mockCall.callList,
		callInsert:     mockCall.callInsert,
		callDelete:     mockCall.callDelete,
	}

	err := group.Add(ctx, []string{"foo@email"})

	assert.NoError(t, err)
}

func TestDeliverySettings(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	mockMembersService := newMockIMembersService(t)
	mockMembersService.EXPECT().Insert("test", &admin.Member{
		Email:            "foo@email",
		DeliverySettings: "test-delivery-settings",
	}).Return(nil)

	mockCall := new(mockCalls)
	mockCall.On("callInsert", ctx, mock.Anything).Return(&admin.Member{}, nil)

	group := &Group{
		name:             "test",
		membersService:   mockMembersService,
		DeliverySettings: "test-delivery-settings",
		Logger:           log.New(os.Stdout, "", log.LstdFlags),
		callList:         mockCall.callList,
		callInsert:       mockCall.callInsert,
		callDelete:       mockCall.callDelete,
	}

	err := group.Add(ctx, []string{"foo@email"})

	assert.NoError(t, err)
}

//nolint:funlen
func TestInit(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{
			Name: "name",
		}, withMockAdminService(ctx, t))

		assert.NoError(t, err)
		assert.IsType(t, &Group{}, adapter)
		assert.Equal(t, "name", adapter.name)
		assert.Equal(t, "", adapter.DeliverySettings)
		assert.Equal(t, "", adapter.Role)
	})

	t.Run("missing config", func(t *testing.T) {
		t.Parallel()

		t.Run("missing name", func(t *testing.T) {
			t.Parallel()

			_, err := Init(ctx, map[gosync.ConfigKey]string{}, withMockAdminService(ctx, t))

			assert.ErrorIs(t, err, gosync.ErrMissingConfig)
			assert.ErrorContains(t, err, Name)
		})
	})

	t.Run("role", func(t *testing.T) {
		t.Parallel()

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{
			Name: "name",
			Role: "role",
		}, withMockAdminService(ctx, t))

		assert.NoError(t, err)
		assert.Equal(t, "role", adapter.Role)
	})

	t.Run("delivery settings", func(t *testing.T) {
		t.Parallel()

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{
			Name:             "name",
			DeliverySettings: "delivery",
		}, withMockAdminService(ctx, t))

		assert.NoError(t, err)
		assert.Equal(t, "delivery", adapter.DeliverySettings)
	})

	t.Run("with logger", func(t *testing.T) {
		t.Parallel()

		logger := log.New(os.Stderr, "custom logger", log.LstdFlags)

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{
			Name: "name",
		}, withMockAdminService(ctx, t), WithLogger(logger))

		assert.NoError(t, err)
		assert.Equal(t, logger, adapter.Logger)
	})

	t.Run("default", func(t *testing.T) {
		t.Parallel()

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{
			Name: "name",
		}, withMockAdminService(ctx, t))

		assert.NoError(t, err)
		assert.NotNil(t, adapter.membersService)
	})
}
