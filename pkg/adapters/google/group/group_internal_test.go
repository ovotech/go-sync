package group

import (
	"context"
	"testing"

	"github.com/ovotech/go-sync/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	admin "google.golang.org/api/admin/directory/v1"
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

	mockMembersService := mocks.NewIMembersService(t)
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
	group.listCall = mockCall.callList

	emails, err := group.Get(ctx)

	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"foo@email", "bar@email"}, emails)
}

func TestGroups_Add(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	mockMembersService := mocks.NewIMembersService(t)
	mockMembersService.EXPECT().Insert("test", &admin.Member{Email: "foo@email"}).Return(nil)
	mockMembersService.EXPECT().Insert("test", &admin.Member{Email: "bar@email"}).Return(nil)

	mockCall := new(mockCalls)
	mockCall.On("callInsert", ctx, mock.Anything).Return(&admin.Member{}, nil)

	group := New(&admin.Service{}, "test")
	group.membersService = mockMembersService
	group.insertCall = mockCall.callInsert

	err := group.Add(ctx, []string{"foo@email", "bar@email"})

	assert.NoError(t, err)
}

func TestGroups_Remove(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	mockMembersService := mocks.NewIMembersService(t)
	mockMembersService.EXPECT().Delete("test", "foo@email").Return(nil)
	mockMembersService.EXPECT().Delete("test", "bar@email").Return(nil)

	mockCall := new(mockCalls)
	mockCall.On("callDelete", ctx, mock.Anything).Twice().Return(nil)

	group := New(&admin.Service{}, "test")
	group.membersService = mockMembersService
	group.deleteCall = mockCall.callDelete

	err := group.Remove(ctx, []string{"foo@email", "bar@email"})

	assert.NoError(t, err)
}
