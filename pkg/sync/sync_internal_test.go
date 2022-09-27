package sync

import (
	"context"
	"errors"
	"testing"

	"github.com/ovotech/go-sync/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	adapter := mocks.NewAdapter(t)
	syncService := New(adapter)

	assert.Empty(t, syncService.cache)
	assert.True(t, syncService.Add)
	assert.True(t, syncService.Remove)
	assert.False(t, syncService.DryRun)
	assert.Zero(t, adapter.Calls)
}

func TestSync_SyncWith(t *testing.T) { //nolint:funlen,maintidx
	t.Parallel()

	ctx := context.TODO()

	t.Run("Add", func(t *testing.T) {
		t.Parallel()

		t.Run("Add successful", func(t *testing.T) {
			t.Parallel()

			source := mocks.NewAdapter(t)
			destination := mocks.NewAdapter(t)

			syncService := New(source)

			source.EXPECT().Get(ctx).Once().Return([]string{"foo", "bar"}, nil)
			destination.EXPECT().Get(ctx).Once().Return([]string{}, nil)
			destination.EXPECT().Add(ctx, []string{"foo", "bar"}).Maybe().Return(nil)
			destination.EXPECT().Add(ctx, []string{"bar", "foo"}).Maybe().Return(nil)

			err := syncService.SyncWith(ctx, destination)

			assert.NoError(t, err)
		})

		t.Run("Add failure", func(t *testing.T) {
			t.Parallel()

			source := mocks.NewAdapter(t)
			destination := mocks.NewAdapter(t)

			syncService := New(source)

			testErr := errors.New("foo") //nolint:goerr113

			source.EXPECT().Get(ctx).Once().Return([]string{"foo", "bar"}, nil)
			destination.EXPECT().Get(ctx).Once().Return([]string{}, nil)
			destination.EXPECT().Add(ctx, []string{"foo", "bar"}).Maybe().Return(testErr)
			destination.EXPECT().Add(ctx, []string{"bar", "foo"}).Maybe().Return(testErr)

			err := syncService.SyncWith(ctx, destination)

			assert.Error(t, err)
			assert.ErrorIs(t, err, testErr)
		})

		t.Run("Add error get", func(t *testing.T) {
			t.Parallel()

			source := mocks.NewAdapter(t)
			destination := mocks.NewAdapter(t)

			syncService := New(source)
			syncService.cache = map[string]bool{}

			testErr := errors.New("foo") //nolint:goerr113

			source.EXPECT().Get(ctx).Once().Return([]string{"foo", "bar"}, nil)
			destination.EXPECT().Get(ctx).Once().Return([]string{}, testErr)

			err := syncService.SyncWith(ctx, destination)

			assert.ErrorIs(t, err, testErr)
		})
	})

	t.Run("Remove", func(t *testing.T) {
		t.Parallel()

		t.Run("Remove successful", func(t *testing.T) {
			t.Parallel()

			source := mocks.NewAdapter(t)
			destination := mocks.NewAdapter(t)

			syncService := New(source)

			source.EXPECT().Get(ctx).Once().Return([]string{}, nil)
			destination.EXPECT().Get(ctx).Once().Return([]string{"foo", "bar"}, nil)
			destination.EXPECT().Remove(ctx, []string{"foo", "bar"}).Maybe().Return(nil)
			destination.EXPECT().Remove(ctx, []string{"bar", "foo"}).Maybe().Return(nil)

			err := syncService.SyncWith(ctx, destination)

			assert.NoError(t, err)
		})

		t.Run("Remove failure", func(t *testing.T) {
			t.Parallel()

			source := mocks.NewAdapter(t)
			destination := mocks.NewAdapter(t)

			syncService := New(source)

			testErr := errors.New("foo") //nolint:goerr113

			source.EXPECT().Get(ctx).Once().Return([]string{}, nil)
			destination.EXPECT().Get(ctx).Once().Return([]string{"foo", "bar"}, nil)
			destination.EXPECT().Remove(ctx, []string{"foo", "bar"}).Maybe().Return(testErr)
			destination.EXPECT().Remove(ctx, []string{"bar", "foo"}).Maybe().Return(testErr)

			err := syncService.SyncWith(ctx, destination)

			assert.Error(t, err)
			assert.ErrorIs(t, err, testErr)
		})

		t.Run("Remove error get", func(t *testing.T) {
			t.Parallel()

			source := mocks.NewAdapter(t)
			destination := mocks.NewAdapter(t)

			syncService := New(source)

			testErr := errors.New("foo") //nolint:goerr113

			source.EXPECT().Get(ctx).Once().Return([]string{}, nil)
			destination.EXPECT().Get(ctx).Once().Return([]string{}, testErr)

			err := syncService.SyncWith(ctx, destination)

			assert.ErrorIs(t, err, testErr)
		})

		t.Run("Remove error remove", func(t *testing.T) {
			t.Parallel()

			source := mocks.NewAdapter(t)
			destination := mocks.NewAdapter(t)

			syncService := New(source)

			testErr := errors.New("foo") //nolint:goerr113

			source.EXPECT().Get(ctx).Once().Return([]string{}, nil)
			destination.EXPECT().Get(ctx).Once().Return([]string{"foo", "bar"}, nil)
			destination.EXPECT().Remove(ctx, []string{"foo", "bar"}).Maybe().Return(testErr)
			destination.EXPECT().Remove(ctx, []string{"bar", "foo"}).Maybe().Return(testErr)

			err := syncService.SyncWith(ctx, destination)

			assert.ErrorIs(t, err, testErr)
		})
	})

	t.Run("Simultaneous", func(t *testing.T) {
		t.Parallel()

		source := mocks.NewAdapter(t)
		destination := mocks.NewAdapter(t)

		syncService := New(source)

		source.EXPECT().Get(ctx).Once().Return([]string{"foo", "bar"}, nil)
		destination.EXPECT().Get(ctx).Once().Return([]string{"fizz", "buzz"}, nil)
		destination.EXPECT().Add(ctx, []string{"foo", "bar"}).Maybe().Return(nil)
		destination.EXPECT().Add(ctx, []string{"bar", "foo"}).Maybe().Return(nil)
		destination.EXPECT().Remove(ctx, []string{"fizz", "buzz"}).Maybe().Return(nil)
		destination.EXPECT().Remove(ctx, []string{"buzz", "fizz"}).Maybe().Return(nil)

		err := syncService.SyncWith(ctx, destination)

		assert.NoError(t, err)
	})

	t.Run("DryRun", func(t *testing.T) {
		t.Parallel()

		t.Run("Add", func(t *testing.T) {
			t.Parallel()

			source := mocks.NewAdapter(t)
			destination := mocks.NewAdapter(t)

			syncService := New(source)
			syncService.DryRun = true

			source.EXPECT().Get(ctx).Once().Return([]string{"foo", "bar"}, nil)
			destination.EXPECT().Get(ctx).Once().Return([]string{}, nil)
			err := syncService.SyncWith(ctx, destination)

			assert.NoError(t, err)
		})

		t.Run("Remove", func(t *testing.T) {
			t.Parallel()

			source := mocks.NewAdapter(t)
			destination := mocks.NewAdapter(t)

			syncService := New(source)
			syncService.DryRun = true

			source.EXPECT().Get(ctx).Once().Return([]string{}, nil)
			destination.EXPECT().Get(ctx).Once().Return([]string{"foo", "bar"}, nil)
			err := syncService.SyncWith(ctx, destination)

			assert.NoError(t, err)
		})
	})

	t.Run("Equal", func(t *testing.T) {
		t.Parallel()

		source := mocks.NewAdapter(t)
		destination := mocks.NewAdapter(t)

		syncService := New(source)

		source.EXPECT().Get(ctx).Once().Return([]string{"foo"}, nil)
		destination.EXPECT().Get(ctx).Once().Return([]string{"foo"}, nil)

		err := syncService.SyncWith(ctx, destination)

		assert.NoError(t, err)
	})

	t.Run("AddRemove", func(t *testing.T) {
		t.Parallel()

		t.Run("Disable Add", func(t *testing.T) {
			t.Parallel()

			source := mocks.NewAdapter(t)
			destination := mocks.NewAdapter(t)

			syncService := New(source)
			syncService.Add = false
			syncService.Remove = true

			source.EXPECT().Get(ctx).Once().Return([]string{"foo", "bar"}, nil)
			destination.EXPECT().Get(ctx).Once().Return([]string{"fizz", "buzz"}, nil)
			destination.EXPECT().Remove(ctx, []string{"fizz", "buzz"}).Maybe().Return(nil)
			destination.EXPECT().Remove(ctx, []string{"buzz", "fizz"}).Maybe().Return(nil)

			err := syncService.SyncWith(ctx, destination)

			assert.NoError(t, err)
		})

		t.Run("Disable Remove", func(t *testing.T) {
			t.Parallel()

			source := mocks.NewAdapter(t)
			destination := mocks.NewAdapter(t)

			syncService := New(source)
			syncService.Add = true
			syncService.Remove = false

			source.EXPECT().Get(ctx).Once().Return([]string{"foo", "bar"}, nil)
			destination.EXPECT().Get(ctx).Once().Return([]string{"fizz", "buzz"}, nil)
			destination.EXPECT().Add(ctx, []string{"foo", "bar"}).Maybe().Return(nil)
			destination.EXPECT().Add(ctx, []string{"bar", "foo"}).Maybe().Return(nil)

			err := syncService.SyncWith(ctx, destination)

			assert.NoError(t, err)
		})
	})
}
