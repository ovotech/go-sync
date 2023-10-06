package plugin

import (
	"context"
	"fmt"

	"github.com/ovotech/go-sync/internal/proto"
	"github.com/ovotech/go-sync/pkg/types"
)

// adapterClient is a duplicate to enable mocking.
//
//goland:noinspection GoUnusedType
type adapterClient interface { //nolint:unused
	proto.AdapterClient
}

// Ensure the Client struct matches the Go Sync adapter spec.
var _ types.Adapter = &Client{}

type Client struct {
	AdapterClient proto.AdapterClient
}

func (c *Client) Init(ctx context.Context, config map[types.ConfigKey]string) error {
	_, err := c.AdapterClient.Init(ctx, &proto.InitRequest{Config: config})
	if err != nil {
		return fmt.Errorf("client.init -> %w", err)
	}

	return nil
}

func (c *Client) Get(ctx context.Context) ([]string, error) {
	response, err := c.AdapterClient.Get(ctx, &proto.GetRequest{})
	if err != nil {
		return nil, fmt.Errorf("client.get -> %w", err)
	}

	return response.Things, nil
}

func (c *Client) Add(ctx context.Context, things []string) error {
	_, err := c.AdapterClient.Add(ctx, &proto.AddRequest{Things: things})
	if err != nil {
		return fmt.Errorf("client.add -> %w", err)
	}

	return nil
}

func (c *Client) Remove(ctx context.Context, things []string) error {
	_, err := c.AdapterClient.Remove(ctx, &proto.RemoveRequest{Things: things})
	if err != nil {
		return fmt.Errorf("client.remove -> %w", err)
	}

	return nil
}
