| **⚠️ Go Sync is under heavy development.** |
|--------------------------------------------|

# Go Sync (all the things)

<div align="center">

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/ovotech/go-sync?label=go&logo=go)
[![Go Doc](https://img.shields.io/static/v1?label=gomarkdoc&message=doc.md&color=blue)](doc.md)
[![Go Report Card](https://goreportcard.com/badge/github.com/ovotech/go-sync?style=flat)](https://goreportcard.com/report/github.com/ovotech/go-sync)
[![Go Reference](https://pkg.go.dev/badge/ovotech/go-sync.svg)](https://pkg.go.dev/ovotech/go-sync)

[![Test Status](https://github.com/ovotech/go-sync/actions/workflows/test.yml/badge.svg)](https://github.com/ovotech/go-sync/actions/workflows/test.yml)
[![GitHub issues](https://img.shields.io/github/issues/ovotech/go-sync?style=flat)](https://github.com/ovotech/go-sync/issues)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/ovotech/go-sync?label=pull+requests&style=flat)](https://github.com/ovotech/go-sync/pull-requests)
[![License](https://img.shields.io/github/license/ovotech/go-sync?style=flat)](/LICENSE)

</div>

![Summary of Go-Sync](assets/sync-architecture.png)

You have people*. You have places you want them to be. Go Sync makes it happen.

_* Doesn't have to be people._

## Installation

```shell
go get github.com/ovotech/go-sync
```

You're ready to Go Sync 🎉

## Usage
[Read the Go documentation here.](doc.md)

Go Sync consists of two fundamental parts:
1. [Sync](#sync-)
2. [Adapters ](#adapters-)

As long as your adapters are compatible, you can synchronise anything.

```go
var (
    source      = mySourceAdapter.New("some-token")
    destination = myDestinationAdapter.New(&myDestinationAdapter.Input{})
)

syncSvc := sync.New(source)

err := syncSvc.SyncWith(context.Background(), destination)
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
Adapters provide a common interface to services. Adapters must implement our [Adapter interface](pkg/ports/adapter.go)
and functionally perform 3 things:

1. Get the things.
2. Add some things.
3. Remove some things.

These things can be anything, but we recommend email addresses. There's no point trying to sync a Slack User ID with a
GitHub user! 🙅

Read about our [built-in adapters here](pkg/adapters), or [build your own](CONTRIBUTING.md).

| *Made with 💚 by OVO's DevEx team* |
|------------------------------------|
