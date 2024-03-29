# Contributing to Go Sync

First of all, thank you for wanting to contribute to Go Sync ✨

## Preparation 🍳

We recommend `asdf`, it's our recommended way of managing our runtime CLIs:

1. [asdf](https://asdf-vm.com/)
2. [asdf-golang](https://github.com/kennyp/asdf-golang) (install via `asdf plugin-add golang`)

Alternatively install Go from the [official documentation](https://go.dev/doc/install).
The version of Go you want can be [found here](https://github.com/ovotech/go-sync/blob/main/go.mod#L3).

For common tasks, we recommend installing [Taskfile](https://taskfile.dev). We run linters to ensure that code being 
checked in matches our quality standards, and have included a Taskfile in this repo containing common commands to assist 
with this.  All tools necessary to action the various Taskfile tasks will be automatically installed on-demand under the
`.task/bin/` sub-directory within this project.

To see what tasks are available, run `task --list-all`. As a bare minimum, we recommend running `task` against your
changes before checking in.

## Developing an adapter 🔌

An adapter's basic functionality is to provide a common interface to a third party service. In order to keep
synchronisation simple, Sync only works with strings. For users, we recommend email addresses as these are usually
common between services.

[Adapter interface documentation](https://pkg.go.dev/github.com/ovotech/go-sync/pkg/ports#Adapter)

Following our specification, your adapter will be compatible with Sync.

We've built a command-line tool to automatically scaffold a new adapter: <https://github.com/ovotech/go-sync-adapter-gen>

<details>
<summary>Example adapter</summary>

```go
package myadapter

import (
	"context"
	"errors"
	"fmt"
	gosync "github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/pkg/ports"
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
func Init(_ context.Context, _ map[gosync.ConfigKey]string, _ ...gosync.ConfigFn[*MyAdapter]) (*MyAdapter, error) {
	return nil, fmt.Errorf("myadapter.init -> %w", gosync.ErrNotImplemented)
}
```

</details>

### Add/Remove

The slice of strings passed to the Add/Remove methods are the diff between the source and destination adapters. If your
service needs a list of users, cache the response from Get in your adapter, and combine the results in your Add/Remove
methods.

### Error handling

Go Sync's error handling convention is to wrap all errors:

```go
if err != nil {
	return fmt.Errorf("some.context.here -> %w", err)
}
```

### Testing

When writing tests, you can autogenerate the mocked clients from interfaces using [Mockery](https://github.com/vektra/mockery):

```sh
task generate
```
