package gosync

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {
	t.Parallel()

	adapter := NewMockAdapter(t)
	syncService := New(adapter)

	assert.Empty(t, syncService.cache)
	assert.Equal(t, RemoveAdd, syncService.OperatingMode)
	assert.False(t, syncService.DryRun)
	assert.Zero(t, adapter.Calls)
}

//nolint:funlen
func TestSync_SyncWith(t *testing.T) { //nolint:maintidx
	t.Parallel()

	ctx := context.TODO()

	t.Run("Add", func(t *testing.T) {
		t.Parallel()

		t.Run("Add successful", func(t *testing.T) {
			t.Parallel()

			source := NewMockAdapter(t)
			destination := NewMockAdapter(t)

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

			source := NewMockAdapter(t)
			destination := NewMockAdapter(t)

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

			source := NewMockAdapter(t)
			destination := NewMockAdapter(t)

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

			source := NewMockAdapter(t)
			destination := NewMockAdapter(t)

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

			source := NewMockAdapter(t)
			destination := NewMockAdapter(t)

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

			source := NewMockAdapter(t)
			destination := NewMockAdapter(t)

			syncService := New(source)

			testErr := errors.New("foo") //nolint:goerr113

			source.EXPECT().Get(ctx).Once().Return([]string{}, nil)
			destination.EXPECT().Get(ctx).Once().Return([]string{}, testErr)

			err := syncService.SyncWith(ctx, destination)

			assert.ErrorIs(t, err, testErr)
		})

		t.Run("Remove error remove", func(t *testing.T) {
			t.Parallel()

			source := NewMockAdapter(t)
			destination := NewMockAdapter(t)

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

		source := NewMockAdapter(t)
		destination := NewMockAdapter(t)

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

			source := NewMockAdapter(t)
			destination := NewMockAdapter(t)

			syncService := New(source)
			syncService.DryRun = true

			source.EXPECT().Get(ctx).Once().Return([]string{"foo", "bar"}, nil)
			destination.EXPECT().Get(ctx).Once().Return([]string{}, nil)
			err := syncService.SyncWith(ctx, destination)

			assert.NoError(t, err)
		})

		t.Run("Remove", func(t *testing.T) {
			t.Parallel()

			source := NewMockAdapter(t)
			destination := NewMockAdapter(t)

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

		source := NewMockAdapter(t)
		destination := NewMockAdapter(t)

		syncService := New(source)

		source.EXPECT().Get(ctx).Once().Return([]string{"foo"}, nil)
		destination.EXPECT().Get(ctx).Once().Return([]string{"foo"}, nil)

		err := syncService.SyncWith(ctx, destination)

		assert.NoError(t, err)
	})

	t.Run("OperatingMode", func(t *testing.T) {
		t.Parallel()

		t.Run("AddOnly", func(t *testing.T) {
			t.Parallel()

			source := NewMockAdapter(t)
			destination := NewMockAdapter(t)

			syncService := New(source)
			syncService.OperatingMode = AddOnly

			source.EXPECT().Get(ctx).Once().Return([]string{"foo"}, nil)
			destination.EXPECT().Get(ctx).Once().Return([]string{"bar"}, nil)
			destination.EXPECT().Add(ctx, []string{"foo"}).Once().Return(nil)

			err := syncService.SyncWith(ctx, destination)

			assert.NoError(t, err)
		})

		t.Run("RemoveOnly", func(t *testing.T) {
			t.Parallel()

			source := NewMockAdapter(t)
			destination := NewMockAdapter(t)

			syncService := New(source)
			syncService.OperatingMode = RemoveOnly

			source.EXPECT().Get(ctx).Once().Return([]string{"foo"}, nil)
			destination.EXPECT().Get(ctx).Once().Return([]string{"bar"}, nil)
			destination.EXPECT().Remove(ctx, []string{"bar"}).Once().Return(nil)

			err := syncService.SyncWith(ctx, destination)

			assert.NoError(t, err)
		})

		t.Run("RemoveAdd", func(t *testing.T) { //nolint:dupl
			t.Parallel()

			source := NewMockAdapter(t)
			destination := NewMockAdapter(t)

			syncService := New(source)
			syncService.OperatingMode = RemoveAdd

			source.EXPECT().Get(ctx).Once().Return([]string{"foo"}, nil)
			destination.EXPECT().Get(ctx).Once().Return([]string{"bar"}, nil)
			destination.EXPECT().Add(ctx, []string{"foo"}).Once().Return(nil)
			destination.EXPECT().Remove(ctx, []string{"bar"}).Once().Return(nil)

			err := syncService.SyncWith(ctx, destination)

			assert.NoError(t, err)
			assert.Equal(t, "Get", destination.Calls[0].Method)
			assert.Equal(t, "Remove", destination.Calls[1].Method)
			assert.Equal(t, "Add", destination.Calls[2].Method)
		})

		t.Run("AddRemove", func(t *testing.T) { //nolint:dupl
			t.Parallel()

			source := NewMockAdapter(t)
			destination := NewMockAdapter(t)

			syncService := New(source)
			syncService.OperatingMode = AddRemove

			source.EXPECT().Get(ctx).Once().Return([]string{"foo"}, nil)
			destination.EXPECT().Get(ctx).Once().Return([]string{"bar"}, nil)
			destination.EXPECT().Add(ctx, []string{"foo"}).Once().Return(nil)
			destination.EXPECT().Remove(ctx, []string{"bar"}).Once().Return(nil)

			err := syncService.SyncWith(ctx, destination)

			assert.NoError(t, err)
			assert.Equal(t, "Get", destination.Calls[0].Method)
			assert.Equal(t, "Add", destination.Calls[1].Method)
			assert.Equal(t, "Remove", destination.Calls[2].Method)
		})
	})

	t.Run("CaseSensitive", func(t *testing.T) {
		t.Parallel()

		t.Run("default", func(t *testing.T) {
			t.Parallel()

			source := NewMockAdapter(t)
			syncService := New(source)

			assert.True(t, syncService.CaseSensitive)
		})

		t.Run("true", func(t *testing.T) {
			t.Parallel()

			source := NewMockAdapter(t)
			source.EXPECT().Get(ctx).Return([]string{"FOO"}, nil)

			destination := NewMockAdapter(t)
			destination.EXPECT().Get(ctx).Return([]string{"foo"}, nil)
			destination.EXPECT().Add(ctx, []string{"FOO"}).Return(nil)
			destination.EXPECT().Remove(ctx, []string{"foo"}).Return(nil)

			syncService := New(source)
			syncService.CaseSensitive = true

			err := syncService.SyncWith(ctx, destination)

			assert.NoError(t, err)
		})

		t.Run("false", func(t *testing.T) {
			t.Parallel()

			source := NewMockAdapter(t)
			source.EXPECT().Get(ctx).Return([]string{"FOO", "BAR"}, nil)

			destination := NewMockAdapter(t)
			destination.EXPECT().Get(ctx).Return([]string{"foo"}, nil)
			destination.EXPECT().Add(ctx, []string{"bar"}).Return(nil)

			syncService := New(source)
			syncService.CaseSensitive = false

			err := syncService.SyncWith(ctx, destination)

			assert.NoError(t, err)
		})
	})

	t.Run("MaximumChanges", func(t *testing.T) {
		t.Parallel()

		source := NewMockAdapter(t)
		destination := NewMockAdapter(t)

		source.EXPECT().Get(ctx).Return([]string{"foo", "bar"}, nil)
		destination.EXPECT().Get(ctx).Return([]string{"fizz"}, nil)
		destination.EXPECT().Add(ctx, mock.Anything).Maybe().Return(nil)
		destination.EXPECT().Remove(ctx, []string{"fizz"}).Return(nil)

		t.Run("0", func(t *testing.T) {
			t.Parallel()

			syncService := New(source)
			syncService.MaximumChanges = 0

			err := syncService.SyncWith(ctx, destination)

			assert.ErrorIs(t, err, ErrTooManyChanges)
		})

		t.Run("1", func(t *testing.T) {
			t.Parallel()

			syncService := New(source)
			syncService.MaximumChanges = 1

			err := syncService.SyncWith(ctx, destination)

			assert.ErrorIs(t, err, ErrTooManyChanges)

			// Set the operating mode to Remove only (only 1 addition), which should pass successfully.
			syncService.OperatingMode = RemoveOnly

			err = syncService.SyncWith(ctx, destination)

			assert.NoError(t, err)
		})

		t.Run("2", func(t *testing.T) {
			t.Parallel()

			syncService := New(source)
			syncService.MaximumChanges = 2

			err := syncService.SyncWith(ctx, destination)

			assert.NoError(t, err)
		})
	})
}
