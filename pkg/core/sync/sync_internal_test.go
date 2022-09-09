package sync

import (
	"errors"
	"testing"

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

func TestNew(t *testing.T) {
	t.Parallel()

	repo := new(mockRepo)
	syncService := New(repo)

	assert.IsType(t, new(Sync), syncService)
	assert.Empty(t, syncService.cache)
	assert.Zero(t, repo.Calls)
}

func TestSync_SetAddRemove(t *testing.T) {
	t.Parallel()

	repo := new(mockRepo)
	syncService := New(repo)

	assert.True(t, syncService.add)
	assert.True(t, syncService.remove)

	syncService.SetAddRemove(true, false)
	assert.True(t, syncService.add)
	assert.False(t, syncService.remove)

	syncService.SetAddRemove(false, true)
	assert.False(t, syncService.add)
	assert.True(t, syncService.remove)

	syncService.SetAddRemove(false, false)
	assert.False(t, syncService.add)
	assert.False(t, syncService.remove)

	syncService.SetAddRemove(true, true)
	assert.True(t, syncService.add)
	assert.True(t, syncService.remove)

	assert.Empty(t, repo.Calls)
}

func TestSync_SetDryRun(t *testing.T) {
	t.Parallel()

	repo := new(mockRepo)
	syncService := New(repo)

	assert.False(t, syncService.dryRun, "Dry Run mode must be disabled by default.")

	syncService.SetDryRun(true)
	assert.True(t, syncService.dryRun)

	syncService.SetDryRun(false)
	assert.False(t, syncService.dryRun)

	assert.Zero(t, repo.Calls, "No calls must be made when setting Dry Run mode.")
}

func TestSync_SyncWith_Equal(t *testing.T) {
	t.Parallel()

	sourceOfTruth := new(mockRepo)
	serviceToBeSynced := new(mockRepo)

	syncService := New(sourceOfTruth)

	sourceOfTruth.On("Get").Once().Return([]string{"foo"}, nil)
	serviceToBeSynced.On("Get").Once().Return([]string{"foo"}, nil)

	success, failure, err := syncService.SyncWith(serviceToBeSynced)

	assert.Empty(t, success, "When accounts are equal, there should be no successes.")
	assert.Empty(t, failure)
	assert.NoError(t, err)
}

func TestSync_SyncWith_AddError_Get(t *testing.T) {
	t.Parallel()

	sourceOfTruth := new(mockRepo)
	serviceToBeSynced := new(mockRepo)

	syncService := New(sourceOfTruth)

	syncService.cache = map[string]bool{}

	sourceOfTruth.On("Get").Once().Return([]string{"foo", "bar"}, nil)
	serviceToBeSynced.On("Get").Once().Return([]string{}, errors.New("get")) //nolint: goerr113

	_, _, err := syncService.SyncWith(serviceToBeSynced)

	assert.ErrorContains(t, err, "get", "This test must return the error from Get.")
}
