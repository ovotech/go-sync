package plugin

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ovotech/go-sync/internal/proto"
	"github.com/ovotech/go-sync/pkg/types"
)

func TestClient_Init(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		var (
			ctx               = context.TODO()
			mockAdapterClient = newMockAdapterClient(t)
			config            = map[types.ConfigKey]string{}
		)

		client := Client{
			AdapterClient: mockAdapterClient,
		}

		mockAdapterClient.EXPECT().Init(ctx, &proto.InitRequest{Config: config}).Return(&proto.InitResponse{}, nil)

		err := client.Init(ctx, config)

		assert.NoError(t, err)
	})

	t.Run("Failure", func(t *testing.T) {
		t.Parallel()

		var (
			ctx               = context.TODO()
			mockAdapterClient = newMockAdapterClient(t)
			config            = map[types.ConfigKey]string{}
			mockErr           = errors.New("test") //nolint: goerr113
		)

		client := Client{
			AdapterClient: mockAdapterClient,
		}

		mockAdapterClient.EXPECT().Init(ctx, &proto.InitRequest{Config: config}).Return(nil, mockErr)

		err := client.Init(ctx, config)

		assert.ErrorIs(t, err, mockErr)
	})
}

func TestClient_Get(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		var (
			ctx               = context.TODO()
			mockAdapterClient = newMockAdapterClient(t)
		)

		client := Client{
			AdapterClient: mockAdapterClient,
		}

		mockAdapterClient.EXPECT().Get(
			ctx,
			&proto.GetRequest{},
		).Return(
			&proto.GetResponse{Things: []string{"foo", "bar"}},
			nil,
		)

		things, err := client.Get(ctx)

		assert.NoError(t, err)
		assert.ElementsMatch(t, things, []string{"foo", "bar"})
	})

	t.Run("Failure", func(t *testing.T) {
		t.Parallel()

		var (
			ctx               = context.TODO()
			mockAdapterClient = newMockAdapterClient(t)
			mockErr           = errors.New("test") //nolint: goerr113
		)

		client := Client{
			AdapterClient: mockAdapterClient,
		}

		mockAdapterClient.EXPECT().Get(ctx, &proto.GetRequest{}).Return(nil, mockErr)

		_, err := client.Get(ctx)

		assert.ErrorIs(t, err, mockErr)
	})
}

func TestClient_Add(t *testing.T) { //nolint: dupl
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		var (
			ctx               = context.TODO()
			mockAdapterClient = newMockAdapterClient(t)
		)

		client := Client{
			AdapterClient: mockAdapterClient,
		}

		mockAdapterClient.EXPECT().Add(
			ctx,
			&proto.AddRequest{Things: []string{"foo", "bar"}},
		).Return(&proto.AddResponse{}, nil)

		err := client.Add(ctx, []string{"foo", "bar"})

		assert.NoError(t, err)
	})

	t.Run("Failure", func(t *testing.T) {
		t.Parallel()

		var (
			ctx               = context.TODO()
			mockAdapterClient = newMockAdapterClient(t)
			mockErr           = errors.New("test") //nolint: goerr113
		)

		client := Client{
			AdapterClient: mockAdapterClient,
		}

		mockAdapterClient.EXPECT().Add(ctx, &proto.AddRequest{Things: []string{"foo", "bar"}}).Return(nil, mockErr)

		err := client.Add(ctx, []string{"foo", "bar"})

		assert.ErrorIs(t, err, mockErr)
	})
}

func TestClient_Remove(t *testing.T) { //nolint: dupl
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		var (
			ctx               = context.TODO()
			mockAdapterClient = newMockAdapterClient(t)
		)

		client := Client{
			AdapterClient: mockAdapterClient,
		}

		mockAdapterClient.EXPECT().Remove(
			ctx,
			&proto.RemoveRequest{Things: []string{"foo", "bar"}},
		).Return(&proto.RemoveResponse{}, nil)

		err := client.Remove(ctx, []string{"foo", "bar"})

		assert.NoError(t, err)
	})

	t.Run("Failure", func(t *testing.T) {
		t.Parallel()

		var (
			ctx               = context.TODO()
			mockAdapterClient = newMockAdapterClient(t)
			mockErr           = errors.New("test") //nolint: goerr113
		)

		client := Client{
			AdapterClient: mockAdapterClient,
		}

		mockAdapterClient.EXPECT().Remove(ctx, &proto.RemoveRequest{Things: []string{"foo", "bar"}}).Return(nil, mockErr)

		err := client.Remove(ctx, []string{"foo", "bar"})

		assert.ErrorIs(t, err, mockErr)
	})
}
