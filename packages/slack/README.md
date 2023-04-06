# Go Sync Adapters - Slack

<div align="center">

[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/ovotech/go-sync?filename=packages/slack/go.mod&label=go&logo=go)](go.mod)
[![Go Report Card](https://goreportcard.com/badge/github.com/ovotech/go-sync/packages/slack?style=flat)](https://goreportcard.com/report/github.com/ovotech/go-sync/packages/slack)
[![Go Reference](https://pkg.go.dev/badge/github.com/ovotech/go-sync.svg)](https://pkg.go.dev/github.com/ovotech/go-sync/packages/slack)

</div>

These adapters synchronise Slack users with [Go Sync](https://github.com/ovotech/go-sync).

[Read the documentation on pkg.go.dev](https://pkg.go.dev/github.com/ovotech/go-sync/packages/slack)

## Installation
```shell
go get github.com/ovotech/go-sync/packages/slack@latest
```

## Adapters

| Adapter                                                                                   | Type  | Summary                                               |
|-------------------------------------------------------------------------------------------|-------|-------------------------------------------------------|
| [conversation](https://pkg.go.dev/github.com/ovotech/go-sync/packages/slack/conversation) | Email | Synchronise emails with a Slack channel/conversation. |
| [usergroup](https://pkg.go.dev/github.com/ovotech/go-sync/packages/slack/usergroup)       | Email | Synchronise emails with a Slack User Group.           |

Can't find an adapter you're looking for? [Why not write your own! ✨](/CONTRIBUTING.md)
