# Contributing to Go Sync

First of all, thank you for wanting to contribute to Go Sync ✨

## Preparation 🍳

We recommend `asdf`, it's our recommended way of managing our runtime CLIs:

1. [asdf](https://asdf-vm.com/)
2. [asdf-golang](https://github.com/kennyp/asdf-golang) (install via `asdf plugin-add golang`)

Alternatively install Go from the [official documentation](https://go.dev/doc/install).
The version of Go you want can be [found here](https://github.com/ovotech/go-sync/blob/main/go.mod#L3).

We also run the following tooling to ensure code quality:

1. [golangci-lint](https://golangci-lint.run/) for code quality.
2. [mockery](https://github.com/vektra/mockery) generates mocks for easy testing.
3. [gci](https://github.com/daixiang0/gci) for consistent, deterministic imports.

We run linters to ensure that code being checked in matches our quality standards, and have included a Makefile in this
repo containing common commands to assist with this.
All tools necessary to action the various Makefile targets will be automatically installed on-demand under the
`hack/bin/` sub-directory within this project.

| Command         | Description                    |
| --------------- | ------------------------------ |
| `make generate` | Generate automated code.       |
| `make lint`     | Lint Go Sync.                  |
| `make lint-fix` | Fix some common linter errors. |

The above is a small subset of available Makefile targets.
Running `make` or `make help` will display a more complete list of available targets.

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
 "github.com/ovotech/go-sync/pkg/ports"
)

// Ensure the adapter type fully satisfies the ports.Adapter interface.
var _ ports.Adapter = &MyAdapter{}

var ErrNotImplemented = errors.New("not implemented")

type MyAdapter struct{}

func New() *MyAdapter {
 return &MyAdapter {}
}

func (m *MyAdapter) Get(_ context.Context) ([]string, error) {
 return nil, fmt.Errorf("myadapter.get -> %w", ErrNotImplemented)
}

func (m *MyAdapter) Add(_ context.Context, _ []string) error {
 return fmt.Errorf("myadapter.add -> %w", ErrNotImplemented)
}

func (m *MyAdapter) Remove(_ context.Context, _ []string) error {
 return fmt.Errorf("myadapter.remove -> %w", ErrNotImplemented)
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

When writing tests, you can autogenerate the mocked clients using [Mockery](#preparation-):

```sh
make generate
```
