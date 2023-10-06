package plugin

import (
	"context"
	"fmt"

	"github.com/ovotech/go-sync/internal/proto"
	"github.com/ovotech/go-sync/pkg/errors"
	"github.com/ovotech/go-sync/pkg/types"
)

var _ proto.AdapterServer = &Server{}

type Server struct {
	proto.UnimplementedAdapterServer

	InitFn  UntypedInitFn
	adapter types.Adapter
}

func (s *Server) Init(ctx context.Context, request *proto.InitRequest) (*proto.InitResponse, error) {
	adapter, err := s.InitFn(ctx, request.Config)
	if err != nil {
		return nil, fmt.Errorf("server.init -> %w", err)
	}

	s.adapter = adapter

	return &proto.InitResponse{}, nil
}

func (s *Server) Get(ctx context.Context, _ *proto.GetRequest) (*proto.GetResponse, error) {
	if s.adapter == nil {
		return nil, fmt.Errorf("server.get -> %w", errors.ErrNotInitialised)
	}

	things, err := s.adapter.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("server.get -> %w", err)
	}

	return &proto.GetResponse{Things: things}, nil
}

func (s *Server) Add(ctx context.Context, request *proto.AddRequest) (*proto.AddResponse, error) {
	if s.adapter == nil {
		return nil, fmt.Errorf("server.add -> %w", errors.ErrNotInitialised)
	}

	err := s.adapter.Add(ctx, request.Things)
	if err != nil {
		return nil, fmt.Errorf("server.add -> %w", err)
	}

	return &proto.AddResponse{}, nil
}

func (s *Server) Remove(ctx context.Context, request *proto.RemoveRequest) (*proto.RemoveResponse, error) {
	if s.adapter == nil {
		return nil, fmt.Errorf("server.remove -> %w", errors.ErrNotInitialised)
	}

	err := s.adapter.Remove(ctx, request.Things)
	if err != nil {
		return nil, fmt.Errorf("server.add -> %w", err)
	}

	return &proto.RemoveResponse{}, nil
}
