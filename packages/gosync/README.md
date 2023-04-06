# Go Sync

<div align="center">

[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/ovotech/go-sync?filename=packages/gosync/go.mod&label=go&logo=go)](go.mod)
[![Go Report Card](https://goreportcard.com/badge/github.com/ovotech/go-sync/packages/gosync?style=flat)](https://goreportcard.com/report/github.com/ovotech/go-sync/packages/gosync)
[![Go Reference](https://pkg.go.dev/badge/github.com/ovotech/go-sync.svg)](https://pkg.go.dev/github.com/ovotech/go-sync/packages/gosync)

</div>

[Read the documentation on pkg.go.dev](https://pkg.go.dev/github.com/ovotech/go-sync/packages/gosync)

## Usage

Go Sync consists of two fundamental parts:

1. [Sync](#sync-)
2. [Adapters](#adapters-)

As long as your adapters are compatible, you can synchronise anything.

```go
ctx := context.Background()

// Initialise a new source adapter.
source, err := src.Init(ctx, map[gosync.ConfigKey]string{
    src.SomeSecret: "something",
})
if err != nil {
    log.Fatal(err)
}

// Initialise a destination adapter.
destination, err := dest.Init(ctx, map[gosync.ConfigKey]string{
    dest.AnotherSecret: "something-else",
})
if err != nil {
    log.Fatal(err)
}

sync := gosync.New(source)

// Synchronise the users in the destination with the source.
err := sync.SyncWith(ctx, destination)
if err != nil {
    log.Fatal(err)
}
```

## Sync 🔄

Sync is the logic that powers the automation. It accepts a source adapter, and synchronises it with destination
adapters.

Sync is only uni-directional by design. You know where your things are, and where you want them to be. It works by:

1. Get a list of things in your source service.
    1. Cache it, so you're not calling your source service more than you have to.
2. Get a list of things in your destination service.
3. Add the things that are missing.
4. Remove the things that shouldn't be there.
5. Repeat from 2 for further adapters.

## Adapters 🔌

Adapters provide a common interface to services.
Adapters must implement our [Adapter interface](https://pkg.go.dev/github.com/ovotech/go-sync/packages/gosync#Adapter) and functionally
perform 3 things:

1. Get the things.
2. Add some things.
3. Remove some things.

These things can be anything, but we recommend email addresses. There's no point trying to sync a Slack User ID with a
GitHub user! 🙅
