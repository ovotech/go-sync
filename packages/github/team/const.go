package team

import "github.com/ovotech/go-sync/packages/gosync"

/*
GitHubToken is the token used to authenticate with GitHub.
See package docs for more information on how to obtain this token.
*/
const GitHubToken gosync.ConfigKey = "github_token"

/*
GitHubOrg is the name of your GitHub organisation.

https://docs.github.com/en/organizations/collaborating-with-groups-in-organizations/about-organizations

For example:

	https://github.com/ovotech/go-sync

`ovotech` is the name of our organisation.
*/
const GitHubOrg gosync.ConfigKey = "github_org"

/*
TeamSlug is the name of your team slug within your organisation.

For example:

	https://github.com/orgs/ovotech/teams/foobar

`foobar` is the name of our team slug.
*/
const TeamSlug gosync.ConfigKey = "team_slug"

/*
DiscoveryMechanism for converting emails into GitHub users and vice versa. Supported values are:
  - [saml]
*/
const DiscoveryMechanism gosync.ConfigKey = "discovery_mechanism"

/*
SamlMuteUserNotFoundErr mutes the UserNotFoundErr if SAML discovery fails to discover a user from GitHub.
*/
const SamlMuteUserNotFoundErr gosync.ConfigKey = "saml_mute_user_not_found_err"
