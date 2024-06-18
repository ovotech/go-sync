| **⚠️ Go Sync is under active development and subject to breaking changes.** |
|-----------------------------------------------------------------------------|

# Go Sync (all the things)

<div align="center">

[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/ovotech/go-sync?label=go&logo=go)](go.mod)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/ovotech/go-sync)](https://github.com/ovotech/go-sync/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/ovotech/go-sync?style=flat)](https://goreportcard.com/report/github.com/ovotech/go-sync)
[![Go Reference](https://pkg.go.dev/badge/github.com/ovotech/go-sync.svg)](https://pkg.go.dev/github.com/ovotech/go-sync)
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
go get github.com/ovotech/go-sync@latest

# Then get the adapters you need.
go get github.com/ovotech/go-sync-adapter-slack@latest
```

You're ready to Go Sync! 🎉

## Usage

Read the [documentation on `pkg.go.dev`](https://pkg.go.dev/github.com/ovotech/go-sync).

Go Sync consists of two fundamental parts:

1. [Sync](#sync-)
2. [Adapters](#adapters-)

As long as your adapters are compatible, you can synchronise anything.

```go
// Initialise an adapter.
destination, err := myAdapter.Init(map[gosync.ConfigKey]string{
	myAdapter.Token:     "some-token",
	myAdapter.Something: "some-value",
})
if err != nil {
	log.Fatal(err)
}

sync := gosync.New(source)

err := sync.SyncWith(context.Background(), destination)
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
Adapters must implement our [Adapter interface](https://pkg.go.dev/github.com/ovotech/go-sync#Adapter) and functionally
perform 3 things:

1. Get the things.
2. Add some things.
3. Remove some things.

These things can be anything, but we recommend email addresses. There's no point trying to sync a Slack User ID with a
GitHub user! 🙅

Can't find an adapter you're looking for? [Why not write your own! ✨](/CONTRIBUTING.md)

### Made with 💚 by OVO Energy's DevEx team

<div>

![DevEx](./assets/devex.png)
![Platforms](./assets/platforms.png)
![Tools](./assets/tools.png)
![Golden Paths](./assets/golden-paths.png)
![Guard Rails](./assets/guard-rails.png)
![For You](./assets/for-you.png)

</div>
