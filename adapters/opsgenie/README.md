# Go Sync Adapters - Opsgenie
These adapters synchronise Opsgenie users with [Go Sync](https://github.com/ovotech/go-sync).

[Read the documentation on pkg.go.dev](https://pkg.go.dev/github.com/ovotech/go-sync/adapters/opsgenie)

## Installation
```shell
go get github.com/ovotech/go-sync/adapters/opsgenie@latest
```

## Adapters

| Adapter                                                                              | Type  | Summary                                                                           |
|:-------------------------------------------------------------------------------------|:------|:----------------------------------------------------------------------------------|
| [oncall](https://pkg.go.dev/github.com/ovotech/go-sync/adapters/opsgenie/oncall)     | Email | Synchronise other adapters with emails of those currently on-call for a schedule. |
| [schedule](https://pkg.go.dev/github.com/ovotech/go-sync/adapters/opsgenie/schedule) | Email | Synchronises emails with an Opsgenie schedule.                                    |


Can't find an adapter you're looking for? [Why not write your own! ✨](/CONTRIBUTING.md)
