# Go Sync Adapters - Opsgenie

<div align="center">

[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/ovotech/go-sync?filename=packages/opsgenie/go.mod&label=go&logo=go)](go.mod)
[![Go Report Card](https://goreportcard.com/badge/github.com/ovotech/go-sync/packages/opsgenie?style=flat)](https://goreportcard.com/report/github.com/ovotech/go-sync/packages/opsgenie)
[![Go Reference](https://pkg.go.dev/badge/github.com/ovotech/go-sync.svg)](https://pkg.go.dev/github.com/ovotech/go-sync/packages/opsgenie)

</div>

These adapters synchronise Opsgenie users with [Go Sync](https://github.com/ovotech/go-sync).

[Read the documentation on pkg.go.dev](https://pkg.go.dev/github.com/ovotech/go-sync/packages/opsgenie)

## Installation
```shell
go get github.com/ovotech/go-sync/packages/opsgenie@latest
```

## Adapters

| Adapter                                                                              | Type  | Summary                                                                           |
|:-------------------------------------------------------------------------------------|:------|:----------------------------------------------------------------------------------|
| [oncall](https://pkg.go.dev/github.com/ovotech/go-sync/packages/opsgenie/oncall)     | Email | Synchronise other adapters with emails of those currently on-call for a schedule. |
| [schedule](https://pkg.go.dev/github.com/ovotech/go-sync/packages/opsgenie/schedule) | Email | Synchronises emails with an Opsgenie schedule.                                    |


Can't find an adapter you're looking for? [Why not write your own! ✨](/CONTRIBUTING.md)
