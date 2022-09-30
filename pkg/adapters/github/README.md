# Go Sync Adapters - GitHub
These adapters synchronise GitHub users.

| Source code    | Go Documentation     | Type  | Summary                                |
|----------------|----------------------|-------|----------------------------------------|
| [team](./team) | [team](/doc.md#team) | Email | Synchronise emails with a GitHub team. |

## Requirements
In order to synchronise with GitHub, you'll need to create a [Personal Access Token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)
with the following permissions:

| Scopes                            | Team |
|-----------------------------------|------|
| repo                              |      |
| repo:status                       |      |
| repo_deployment Access deployment |      |
| public_repo                       |      |
| repo:invite                       |      |
| security_events                   |      |
| workflow                          |      |
| write:packages                    |      |
| read:packages                     |      |
| delete:packages                   |      |
| admin:org                         | ✔️   |
| write:org                         | ✔️   |
| read:org                          | ✔️   |
| manage_runners:org                |      |
| admin:public_key                  |      |
| write:public_key                  |      |
| read:public_key                   |      |
| admin:repo_hook                   |      |
| write:repo_hook                   |      |
| read:repo_hook                    |      |
| admin:org_hook                    |      |
| gist                              |      |
| notifications                     |      |
| user                              |      |
| read:user                         |      |
| user:email                        |      |
| user:follow                       |      |
| delete_repo                       |      |
| write:discussion                  |      |
| read:discussion                   |      |
| admin:enterprise                  |      |
| manage_runners:enterprise         |      |
| manage_billing:enterprise         |      |
| read:enterprise                   |      |
| project                           |      |
| read:project                      |      |
| admin:gpg_key                     |      |
| write:gpg_key                     |      |
| read:gpg_key                      |      |
| admin:ssh_signing_key             |      |
| write:ssh_signing_key             |      |
| read:ssh_signing_key              |      |




Can't find an adapter you're looking for? [Why not contribute your own! ✨](/CONTRIBUTING.md)
