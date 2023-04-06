# Go Sync Adapters - Terraform Cloud

<div align="center">

[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/ovotech/go-sync?filename=packages/terraformcloud/go.mod&label=go&logo=go)](go.mod)
[![Go Report Card](https://goreportcard.com/badge/github.com/ovotech/go-sync/packages/terraformcloud?style=flat)](https://goreportcard.com/report/github.com/ovotech/go-sync/packages/terraformcloud)
[![Go Reference](https://pkg.go.dev/badge/github.com/ovotech/go-sync.svg)](https://pkg.go.dev/github.com/ovotech/go-sync/packages/terraformcloud)

</div>

These adapters synchronise a Terraform Cloud organisation with [Go Sync](https://github.com/ovotech/go-sync).

[Read the documentation on pkg.go.dev](https://pkg.go.dev/github.com/ovotech/go-sync/packages/terraformcloud)

## Installation
```shell
go get github.com/ovotech/go-sync/packages/terraformcloud@latest
```

## Adapters

| Adapter                                                                            | Type  | Summary                                              |
|------------------------------------------------------------------------------------|-------|------------------------------------------------------|
| [team](https://pkg.go.dev/github.com/ovotech/go-sync/packages/terraformcloud/team) | Email | Synchronise teams in a Terraform Cloud organisation. |
| [user](https://pkg.go.dev/github.com/ovotech/go-sync/packages/terraformcloud/user) | Email | Synchronise users in a Terraform Cloud team.         |

Can't find an adapter you're looking for? [Why not write your own! ✨](/CONTRIBUTING.md)
