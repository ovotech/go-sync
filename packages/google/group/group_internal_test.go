package group

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	admin "google.golang.org/api/admin/directory/v1"

	"github.com/ovotech/go-sync/packages/gosync"
)

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

func TestNew(t *testing.T) {
	t.Parallel()

	group := New(&admin.Service{}, "test")

	assert.Equal(t, "test", group.name)
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

	group := New(&admin.Service{}, "test")
	group.membersService = mockMembersService
	group.callList = mockCall.callList

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

	group := New(&admin.Service{}, "test")
	group.membersService = mockMembersService
	group.callInsert = mockCall.callInsert

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

	group := New(&admin.Service{}, "test")
	group.membersService = mockMembersService
	group.callDelete = mockCall.callDelete

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

	group := New(&admin.Service{}, "test")
	group.Role = "test-role"

	group.membersService = mockMembersService
	group.callInsert = mockCall.callInsert

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

	group := New(&admin.Service{}, "test")
	group.DeliverySettings = "test-delivery-settings"

	group.membersService = mockMembersService
	group.callInsert = mockCall.callInsert

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
			GoogleAuthenticationMechanism: "_testing_",
			Name:                          "name",
		})

		assert.NoError(t, err)
		assert.IsType(t, &Group{}, adapter)
		assert.Equal(t, "name", adapter.(*Group).name)
		assert.Equal(t, "", adapter.(*Group).DeliverySettings)
		assert.Equal(t, "", adapter.(*Group).Role)
	})

	t.Run("missing config", func(t *testing.T) {
		t.Parallel()

		t.Run("missing authentication", func(t *testing.T) {
			t.Parallel()

			_, err := Init(ctx, map[gosync.ConfigKey]string{
				Name: "name",
			})

			assert.ErrorIs(t, err, gosync.ErrMissingConfig)
			assert.ErrorContains(t, err, GoogleAuthenticationMechanism)
		})

		t.Run("missing name", func(t *testing.T) {
			t.Parallel()

			_, err := Init(ctx, map[gosync.ConfigKey]string{
				GoogleAuthenticationMechanism: "default",
			})

			assert.ErrorIs(t, err, gosync.ErrMissingConfig)
			assert.ErrorContains(t, err, Name)
		})
	})

	t.Run("invalid config", func(t *testing.T) {
		t.Parallel()

		_, err := Init(ctx, map[gosync.ConfigKey]string{
			GoogleAuthenticationMechanism: "foo",
			Name:                          "name",
		})

		assert.ErrorIs(t, err, gosync.ErrInvalidConfig)
		assert.ErrorContains(t, err, "foo")
	})

	t.Run("role", func(t *testing.T) {
		t.Parallel()

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{
			GoogleAuthenticationMechanism: "_testing_",
			Name:                          "name",
			Role:                          "role",
		})

		assert.NoError(t, err)
		assert.Equal(t, "role", adapter.(*Group).Role)
	})

	t.Run("delivery settings", func(t *testing.T) {
		t.Parallel()

		adapter, err := Init(ctx, map[gosync.ConfigKey]string{
			GoogleAuthenticationMechanism: "_testing_",
			Name:                          "name",
			DeliverySettings:              "delivery",
		})

		assert.NoError(t, err)
		assert.Equal(t, "delivery", adapter.(*Group).DeliverySettings)
	})
}
