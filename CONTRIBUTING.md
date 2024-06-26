# Contributing to Go Sync

First of all, thank you for wanting to contribute to Go Sync ✨

## Preparation 🍳

We recommend using [Mise](https://mise.jdx.dev/) to manage your runtime CLIs.

We run linters to ensure that code being checked in matches our quality standards, and we use Mise tasks to assist 
with this.  All tools necessary to action the various Taskfile tasks will be automatically installed by Mise.

To see what tasks are available, run `mise tasks`. As a bare minimum, we recommend running `mise run lint` against your
changes before checking in.

## Developing an adapter 🔌

An adapter's basic functionality is to provide a common interface to a third party service. In order to keep
synchronisation simple, Sync only works with strings. For users, we recommend email addresses as these are usually
common between services.

[Adapter interface documentation](https://pkg.go.dev/github.com/ovotech/go-sync#Adapter)

Following our specification, your adapter will be compatible with Sync.

<details>
<summary>Example adapter</summary>

```go
package myadapter

import (
	"context"
	"errors"
	"fmt"
	gosync "github.com/ovotech/go-sync"
)

var (
	// Ensure [myadapter.myAdapter] fully satisfies the [gosync.Adapter] interface.
	_ gosync.Adapter            = &MyAdapter{}
	// Ensure [myadapter.Init] fully satisfies the [gosync.InitFn] type.
	_ gosync.InitFn[*MyAdapter] = Init
)

type MyAdapter struct{
	example struct{}
}

// Get things in MyAdapter.
func (m *MyAdapter) Get(_ context.Context) ([]string, error) {
	return nil, fmt.Errorf("myadapter.get -> %w", gosync.ErrNotImplemented)
}

// Add things to MyAdapter.
func (m *MyAdapter) Add(_ context.Context, _ []string) error {
	return fmt.Errorf("myadapter.add -> %w", gosync.ErrNotImplemented)
}

// Remove things from MyAdapter.
func (m *MyAdapter) Remove(_ context.Context, _ []string) error {
	return fmt.Errorf("myadapter.remove -> %w", gosync.ErrNotImplemented)
}

// WithExample passes a custom example struct.
func WithExample(example struct{}) gosync.ConfigFn[*MyAdapter] {
	return func(m *MyAdapter) {
		m.example = example
    }
}

// Init a new [myadapter.MyAdapter].
func Init(_ context.Context, _ map[string]string, _ ...gosync.ConfigFn[*MyAdapter]) (*MyAdapter, error) {
	return nil, fmt.Errorf("myadapter.init -> %w", gosync.ErrNotImplemented)
}
```

</details>

### Add/Remove

The slice of strings passed to the Add/Remove methods are the diff between the source and destination adapters. If your
service needs a list of users, cache the response from Get in your adapter, and combine the results in your Add/Remove
methods.
