# Go Sync Adapters - Slack
These adapters synchronise Slack users.

| Source code                    | Go Documentation                     | Type  | Summary                                               |
|--------------------------------|--------------------------------------|-------|-------------------------------------------------------|
| [conversation](./conversation) | [conversation](/doc.md#conversation) | Email | Synchronise emails with a Slack channel/conversation. |
| [usergroup](./usergroup)       | [usergroup](/doc.md#usergroup)       | Email | Synchronise emails with a Slack User Group.           |

## Requirements
In order to synchronise with Slack, you'll need to [create a Slack app](https://api.slack.com/authentication/basics) 
with the following OAuth permissions:

| Bot Token Scopes                                                  | Conversation | UserGroup |
|-------------------------------------------------------------------|--------------|-----------|
| [users:read](https://api.slack.com/scopes/users:read)             | ✔️           | ✔️        |
| [users:read.email](https://api.slack.com/scopes/users:read.email) | ✔️           | ✔️        |
| [channels:manage](https://api.slack.com/scopes/channels:manage)   | ✔️           |           |
| [channels:read](https://api.slack.com/scopes/channels:read)       | ✔️           |           |
| [groups:read](https://api.slack.com/scopes/groups:read)           | ✔️           |           |
| [groups:write](https://api.slack.com/scopes/groups:write)         | ✔️           |           |
| [im:write](https://api.slack.com/scopes/im:write)                 | ✔️           |           |
| [mpim:write](https://api.slack.com/scopes/mpim:write)             | ✔️           |           |
| [usergroups:read](https://api.slack.com/scopes/usergroups:read)   |              | ✔️        |
| [usergroups:write](https://api.slack.com/scopes/usergroups:write) |              | ✔️        |


Can't find an adapter you're looking for? [Why not contribute your own! ✨](/CONTRIBUTING.md)
