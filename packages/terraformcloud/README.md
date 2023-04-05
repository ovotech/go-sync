# Go Sync Adapters - Terraform Cloud
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
