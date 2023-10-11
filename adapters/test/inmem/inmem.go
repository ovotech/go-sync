package inmem

import (
	"context"

	"github.com/ovotech/go-sync/pkg/types"
)

var (
	_ types.Adapter        = &InMem{}
	_ types.InitFn[*InMem] = Init
)

type InMem struct {
	store map[string]bool
}

func (i *InMem) Get(_ context.Context) ([]string, error) {
	things := make([]string, 0, len(i.store))

	for k := range i.store {
		things = append(things, k)
	}

	return things, nil
}

func (i *InMem) Add(_ context.Context, things []string) error {
	for _, thing := range things {
		i.store[thing] = true
	}

	return nil
}

func (i *InMem) Remove(_ context.Context, things []string) error {
	for _, thing := range things {
		delete(i.store, thing)
	}

	return nil
}

func Init(_ context.Context, _ map[types.ConfigKey]string, _ ...types.ConfigFn[*InMem]) (*InMem, error) {
	return &InMem{
		store: make(map[string]bool),
	}, nil
}
