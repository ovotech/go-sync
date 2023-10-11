package plugin

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ovotech/go-sync/internal/proto"
	gosyncerrors "github.com/ovotech/go-sync/pkg/errors"
	"github.com/ovotech/go-sync/pkg/types"
)

type initFn struct {
	mock.Mock
}

func (i *initFn) InitFn(ctx context.Context, config map[types.ConfigKey]string) (types.Adapter, error) {
	args := i.Called(ctx, config)

	return args.Get(0).(types.Adapter), args.Error(1) //nolint:wrapcheck
}

// Ensure the mock adapter matches the Go Sync adapter.
var _ types.Adapter = &adapter{}

type adapter struct {
	mock.Mock
}

func (a *adapter) Get(ctx context.Context) ([]string, error) {
	args := a.Called(ctx)

	return args.Get(0).([]string), args.Error(1) //nolint:wrapcheck
}

func (a *adapter) Add(ctx context.Context, things []string) error {
	args := a.Called(ctx, things)

	return args.Error(0) //nolint:wrapcheck
}

func (a *adapter) Remove(ctx context.Context, things []string) error {
	args := a.Called(ctx, things)

	return args.Error(0) //nolint:wrapcheck
}

func TestServer_Init(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		var (
			mockInit    = new(initFn)
			mockAdapter = new(adapter)
			config      = map[types.ConfigKey]string{}
		)

		mockInit.On("InitFn", ctx, config).Return(mockAdapter, nil)

		srv := Server{
			InitFn: mockInit.InitFn,
		}

		response, err := srv.Init(ctx, &proto.InitRequest{Config: config})

		assert.NoError(t, err)
		assert.Equal(t, mockAdapter, srv.adapter)
		assert.Equal(t, &proto.InitResponse{}, response)
	})

	t.Run("Failure", func(t *testing.T) {
		t.Parallel()

		var (
			mockInit    = new(initFn)
			mockAdapter = new(adapter)
			config      = map[types.ConfigKey]string{}
			mockErr     = errors.New("test") //nolint: goerr113
		)

		mockInit.On("InitFn", ctx, config).Return(mockAdapter, mockErr)

		srv := Server{
			InitFn: mockInit.InitFn,
		}

		_, err := srv.Init(ctx, &proto.InitRequest{Config: config})

		assert.ErrorIs(t, err, mockErr)
	})
}

func TestServer_Get(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockAdapter := new(adapter)
		mockAdapter.On("Get", ctx).Return([]string{"foo", "bar"}, nil)

		srv := Server{
			adapter: mockAdapter,
		}

		response, err := srv.Get(ctx, &proto.GetRequest{})

		assert.NoError(t, err)
		assert.ElementsMatch(t, []string{"foo", "bar"}, response.Things)
	})

	t.Run("Not initialised", func(t *testing.T) {
		t.Parallel()

		srv := Server{}

		_, err := srv.Get(ctx, &proto.GetRequest{})

		assert.ErrorIs(t, err, gosyncerrors.ErrNotInitialised)
	})

	t.Run("Failure", func(t *testing.T) {
		t.Parallel()

		mockErr := errors.New("test") //nolint: goerr113

		mockAdapter := new(adapter)
		mockAdapter.On("Get", ctx).Return([]string(nil), mockErr)

		srv := Server{
			adapter: mockAdapter,
		}

		_, err := srv.Get(ctx, &proto.GetRequest{})

		assert.ErrorIs(t, err, mockErr)
	})
}

func TestServer_Add(t *testing.T) { //nolint: dupl
	t.Parallel()

	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockAdapter := new(adapter)
		mockAdapter.On("Add", ctx, []string{"foo", "bar"}).Return(nil)

		srv := Server{
			adapter: mockAdapter,
		}

		_, err := srv.Add(ctx, &proto.AddRequest{Things: []string{"foo", "bar"}})

		assert.NoError(t, err)
	})

	t.Run("Not initialised", func(t *testing.T) {
		t.Parallel()

		srv := Server{}

		_, err := srv.Add(ctx, &proto.AddRequest{})

		assert.ErrorIs(t, err, gosyncerrors.ErrNotInitialised)
	})

	t.Run("Failure", func(t *testing.T) {
		t.Parallel()

		mockErr := errors.New("test") //nolint: goerr113

		mockAdapter := new(adapter)
		mockAdapter.On("Add", ctx, []string{"foo", "bar"}).Return(mockErr)

		srv := Server{
			adapter: mockAdapter,
		}

		_, err := srv.Add(ctx, &proto.AddRequest{Things: []string{"foo", "bar"}})

		assert.ErrorIs(t, err, mockErr)
	})
}

func TestServer_Remove(t *testing.T) { //nolint: dupl
	t.Parallel()

	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockAdapter := new(adapter)
		mockAdapter.On("Remove", ctx, []string{"foo", "bar"}).Return(nil)

		srv := Server{
			adapter: mockAdapter,
		}

		_, err := srv.Remove(ctx, &proto.RemoveRequest{Things: []string{"foo", "bar"}})

		assert.NoError(t, err)
	})

	t.Run("Not initialised", func(t *testing.T) {
		t.Parallel()

		srv := Server{}

		_, err := srv.Remove(ctx, &proto.RemoveRequest{})

		assert.ErrorIs(t, err, gosyncerrors.ErrNotInitialised)
	})

	t.Run("Failure", func(t *testing.T) {
		t.Parallel()

		mockErr := errors.New("test") //nolint: goerr113

		mockAdapter := new(adapter)
		mockAdapter.On("Remove", ctx, []string{"foo", "bar"}).Return(mockErr)

		srv := Server{
			adapter: mockAdapter,
		}

		_, err := srv.Remove(ctx, &proto.RemoveRequest{Things: []string{"foo", "bar"}})

		assert.ErrorIs(t, err, mockErr)
	})
}
