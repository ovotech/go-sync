package sync_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/ovotech/go-sync/pkg/core/sync"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) Get() ([]string, error) {
	args := m.Called()

	return args.Get(0).([]string), args.Error(1) //nolint:wrapcheck
}

func (m *mockRepo) Add(i ...string) ([]string, []error, error) {
	args := m.Called(i)

	return args.Get(0).([]string), args.Get(1).([]error), args.Error(2) //nolint:wrapcheck
}

func (m *mockRepo) Remove(i ...string) ([]string, []error, error) {
	args := m.Called(i)

	return args.Get(0).([]string), args.Get(1).([]error), args.Error(2) //nolint:wrapcheck
}

func less(a, b string) bool { return a < b }

func hocArgumentsAnyOrder(b ...string) func([]string) bool {
	return func(a []string) bool {
		return cmp.Diff(a, b, cmpopts.SortSlices(less)) == ""
	}
}

func TestHocArgumentsAnyOrder(t *testing.T) {
	t.Parallel()

	a := []string{"a", "b", "c"}
	assert.True(t, hocArgumentsAnyOrder("c", "b", "a")(a))
	assert.False(t, hocArgumentsAnyOrder("d", "b", "a")(a))
}

func TestSync_SyncWith_Add(t *testing.T) {
	t.Parallel()

	sourceOfTruth := new(mockRepo)
	serviceToBeSynced := new(mockRepo)

	syncService := sync.New(sourceOfTruth)

	sourceOfTruth.On("Get").Once().Return([]string{"foo", "bar"}, nil)
	serviceToBeSynced.On("Get").Once().Return([]string{}, nil)
	serviceToBeSynced.On("Add", mock.MatchedBy(hocArgumentsAnyOrder("foo", "bar"))).Once().Return([]string{"foo", "bar"}, []error{}, nil) //nolint:lll

	success, failure, err := syncService.SyncWith(serviceToBeSynced)

	assert.ElementsMatch(t, []string{"foo", "bar"}, success)
	assert.Empty(t, failure)
	assert.NoError(t, err)
}

func TestSync_SyncWith_Add_Failure(t *testing.T) {
	t.Parallel()

	sourceOfTruth := new(mockRepo)
	serviceToBeSynced := new(mockRepo)

	syncService := sync.New(sourceOfTruth)

	sourceOfTruth.On("Get").Once().Return([]string{"foo", "bar"}, nil)
	serviceToBeSynced.On("Get").Once().Return([]string{}, nil)
	serviceToBeSynced.On("Add", mock.MatchedBy(hocArgumentsAnyOrder("foo", "bar"))).Once().Return([]string{"foo"}, []error{errors.New("bar")}, nil) //nolint:lll,goerr113

	success, failure, err := syncService.SyncWith(serviceToBeSynced)

	assert.ElementsMatch(t, []string{"foo"}, success)
	assert.ElementsMatch(t, []error{errors.New("bar")}, failure) //nolint:goerr113
	assert.NoError(t, err)
}

func TestSync_SyncWith_Add_Error(t *testing.T) {
	t.Parallel()

	sourceOfTruth := new(mockRepo)
	serviceToBeSynced := new(mockRepo)

	syncService := sync.New(sourceOfTruth)

	sourceOfTruth.On("Get").Once().Return([]string{"foo", "bar"}, nil)
	serviceToBeSynced.On("Get").Once().Return([]string{}, nil)
	serviceToBeSynced.On("Add", mock.MatchedBy(hocArgumentsAnyOrder("foo", "bar"))).Once().Return([]string{}, []error{}, errors.New("add")) //nolint:lll,goerr113

	_, _, err := syncService.SyncWith(serviceToBeSynced)

	assert.ErrorContains(t, err, "add", "This test must return the error from Add.")
}

func TestSync_SyncWith_Remove(t *testing.T) {
	t.Parallel()

	sourceOfTruth := new(mockRepo)
	serviceToBeSynced := new(mockRepo)

	syncService := sync.New(sourceOfTruth)

	sourceOfTruth.On("Get").Once().Return([]string{}, nil)
	serviceToBeSynced.On("Get").Once().Return([]string{"foo", "bar"}, nil)
	serviceToBeSynced.On("Remove", mock.MatchedBy(hocArgumentsAnyOrder("foo", "bar"))).Once().Return([]string{"foo", "bar"}, []error{}, nil) //nolint:lll

	success, failure, err := syncService.SyncWith(serviceToBeSynced)

	assert.ElementsMatch(t, []string{"foo", "bar"}, success)
	assert.Empty(t, failure)
	assert.NoError(t, err)
}

func TestSync_SyncWith_Remove_Failure(t *testing.T) {
	t.Parallel()

	sourceOfTruth := new(mockRepo)
	serviceToBeSynced := new(mockRepo)

	syncService := sync.New(sourceOfTruth)

	sourceOfTruth.On("Get").Once().Return([]string{}, nil)
	serviceToBeSynced.On("Get").Once().Return([]string{"foo", "bar"}, nil)
	serviceToBeSynced.On("Remove", mock.MatchedBy(hocArgumentsAnyOrder("foo", "bar"))).Once().Return([]string{"foo"}, []error{errors.New("bar")}, nil) //nolint:lll,goerr113

	success, failure, err := syncService.SyncWith(serviceToBeSynced)

	assert.ElementsMatch(t, []string{"foo"}, success)
	assert.ElementsMatch(t, []error{errors.New("bar")}, failure) //nolint:goerr113
	assert.NoError(t, err)
}

func TestSync_SyncWith_Remove_Error_Get(t *testing.T) {
	t.Parallel()

	sourceOfTruth := new(mockRepo)
	serviceToBeSynced := new(mockRepo)

	syncService := sync.New(sourceOfTruth)

	sourceOfTruth.On("Get").Once().Return([]string{}, nil)
	serviceToBeSynced.On("Get").Once().Return([]string{}, errors.New("get")) //nolint:goerr113

	_, _, err := syncService.SyncWith(serviceToBeSynced)

	assert.ErrorContains(t, err, "get", "This test must return the error from Get.")
}

func TestSync_SyncWith_Remove_Error_Remove(t *testing.T) {
	t.Parallel()

	sourceOfTruth := new(mockRepo)
	serviceToBeSynced := new(mockRepo)

	syncService := sync.New(sourceOfTruth)

	sourceOfTruth.On("Get").Once().Return([]string{}, nil)
	serviceToBeSynced.On("Get").Once().Return([]string{"foo", "bar"}, nil)
	serviceToBeSynced.On("Remove", mock.MatchedBy(hocArgumentsAnyOrder("foo", "bar"))).Once().Return([]string{}, []error{}, errors.New("add")) //nolint:lll,goerr113

	_, _, err := syncService.SyncWith(serviceToBeSynced)

	assert.ErrorContains(t, err, "add", "This test must return the error from Remove.")
}

func TestSync_SyncWith_Simultaneous(t *testing.T) {
	t.Parallel()

	sourceOfTruth := new(mockRepo)
	serviceToBeSynced := new(mockRepo)

	syncService := sync.New(sourceOfTruth)

	sourceOfTruth.On("Get").Once().Return([]string{"foo", "bar"}, nil)
	serviceToBeSynced.On("Get").Once().Return([]string{"fizz", "buzz"}, nil)
	serviceToBeSynced.On("Add", mock.MatchedBy(hocArgumentsAnyOrder("foo", "bar"))).Once().Return([]string{"foo", "bar"}, []error{}, nil)        //nolint:lll
	serviceToBeSynced.On("Remove", mock.MatchedBy(hocArgumentsAnyOrder("fizz", "buzz"))).Once().Return([]string{"fizz", "buzz"}, []error{}, nil) //nolint:lll

	success, failure, err := syncService.SyncWith(serviceToBeSynced)

	assert.ElementsMatch(t, []string{"fizz", "buzz", "foo", "bar"}, success)
	assert.Empty(t, failure)
	assert.NoError(t, err)
}

func TestSync_SyncWith_DryRun_Add(t *testing.T) {
	t.Parallel()

	sourceOfTruth := new(mockRepo)
	serviceToBeSynced := new(mockRepo)

	syncService := sync.New(sourceOfTruth)
	syncService.SetDryRun(true)

	sourceOfTruth.On("Get").Once().Return([]string{"foo", "bar"}, nil)
	serviceToBeSynced.On("Get").Once().Return([]string{}, nil)
	success, failure, err := syncService.SyncWith(serviceToBeSynced)

	assert.ElementsMatch(t, []string{"foo", "bar"}, success)
	assert.Empty(t, failure)
	assert.NoError(t, err)
}

func TestSync_SyncWith_DryRun_Remove(t *testing.T) {
	t.Parallel()

	sourceOfTruth := new(mockRepo)
	serviceToBeSynced := new(mockRepo)

	syncService := sync.New(sourceOfTruth)
	syncService.SetDryRun(true)

	sourceOfTruth.On("Get").Once().Return([]string{}, nil)
	serviceToBeSynced.On("Get").Once().Return([]string{"foo", "bar"}, nil)
	success, failure, err := syncService.SyncWith(serviceToBeSynced)

	assert.ElementsMatch(t, []string{"foo", "bar"}, success)
	assert.Empty(t, failure)
	assert.NoError(t, err)
}

func TestSync_SyncWith_AddRemove_DisableAdd(t *testing.T) {
	t.Parallel()

	sourceOfTruth := new(mockRepo)
	serviceToBeSynced := new(mockRepo)

	syncService := sync.New(sourceOfTruth)
	syncService.SetAddRemove(false, true)

	sourceOfTruth.On("Get").Once().Return([]string{"foo", "bar"}, nil)
	serviceToBeSynced.On("Get").Once().Return([]string{"fizz", "buzz"}, nil)
	serviceToBeSynced.On("Remove", mock.MatchedBy(hocArgumentsAnyOrder("fizz", "buzz"))).Once().Return([]string{"fizz", "buzz"}, []error{}, nil) //nolint:lll

	success, failure, err := syncService.SyncWith(serviceToBeSynced)

	assert.ElementsMatch(t, []string{"fizz", "buzz"}, success)
	assert.Empty(t, failure)
	assert.NoError(t, err)
}

func TestSync_SyncWith_AddRemove_DisableRemove(t *testing.T) {
	t.Parallel()

	sourceOfTruth := new(mockRepo)
	serviceToBeSynced := new(mockRepo)

	syncService := sync.New(sourceOfTruth)
	syncService.SetAddRemove(true, false)

	sourceOfTruth.On("Get").Once().Return([]string{"foo", "bar"}, nil)
	serviceToBeSynced.On("Get").Once().Return([]string{"fizz", "buzz"}, nil)
	serviceToBeSynced.On("Add", mock.MatchedBy(hocArgumentsAnyOrder("foo", "bar"))).Once().Return([]string{"foo", "bar"}, []error{}, nil) //nolint:lll

	success, failure, err := syncService.SyncWith(serviceToBeSynced)

	assert.ElementsMatch(t, []string{"foo", "bar"}, success)
	assert.Empty(t, failure)
	assert.NoError(t, err)
}
