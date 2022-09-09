package slack_test

import (
	"testing"

	slacksync "github.com/ovotech/go_sync/pkg/slack"
	"github.com/ovotech/go_sync/test/httpclient"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

func TestSlackUserGroup_Get(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	httpClient := httpclient.New(mockRoundTrip)
	mockRoundTrip.AddSlackRequest(
		"/api/usergroups.users.list",
		map[string]string{"token": "TOKEN_HERE", "usergroup": "test"},
		`
{
	"ok": true,
	"users": [
		"foo",
		"bar"
	]
}
`)
	mockRoundTrip.AddSlackRequest(
		"/api/users.info",
		map[string]string{"users": "foo,bar", "token": "TOKEN_HERE", "usergroup": "test", "include_locale": "true"},
		`
{
	"ok": true,
	"users": [
		{
			"id": "foo",
			"profile": {
				"email": "foo@test.test"
			}
		},
		{
			"id": "bar",
			"profile": {
				"email": "bar@test.test"
			}
		}
	]
}
`)

	slackClient := slack.New("TOKEN_HERE", slack.OptionHTTPClient(httpClient))
	repo := slacksync.NewUserGroupService(slackClient, "test")

	accounts, err := repo.Get()

	assert.NoError(t, err, "No error should occur.")
	assert.ElementsMatch(t, []string{"foo@test.test", "bar@test.test"}, accounts, "Must return two accounts, foo and bar.")
}

func TestSlackUserGroup_Get_Error_ChannelDoesNotExist(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	httpClient := httpclient.New(mockRoundTrip)
	mockRoundTrip.AddSlackRequest(
		"/api/usergroups.users.list",
		map[string]string{"token": "TOKEN_HERE", "usergroup": "test"},
		`
{
	"ok": false,
	"error": "no_such_subteam"
}
`)

	slackClient := slack.New("TOKEN_HERE", slack.OptionHTTPClient(httpClient))
	repo := slacksync.NewUserGroupService(slackClient, "test")

	_, err := repo.Get()

	assert.Error(t, err)
	assert.ErrorContains(t, err, "no_such_subteam")
}

func TestSlackUserGroup_Get_Error_UserDoesNotExist(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	httpClient := httpclient.New(mockRoundTrip)
	mockRoundTrip.AddSlackRequest(
		"/api/usergroups.users.list",
		map[string]string{"token": "TOKEN_HERE", "usergroup": "test"},
		`
{
	"ok": true,
	"users": [
		"foo"
	]
}
`)
	mockRoundTrip.AddSlackRequest(
		"/api/users.info",
		map[string]string{"users": "foo", "token": "TOKEN_HERE", "usergroup": "test", "include_locale": "true"},
		`
{
	"ok": false,
	"error": "user_not_found" 
}
`)

	slackClient := slack.New("TOKEN_HERE", slack.OptionHTTPClient(httpClient))
	repo := slacksync.NewUserGroupService(slackClient, "test")

	_, err := repo.Get()

	assert.Error(t, err)
	assert.ErrorContains(t, err, "user_not_found")
}

func TestSlackUserGroup_Add_NoCache(t *testing.T) { //nolint:funlen
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	httpClient := httpclient.New(mockRoundTrip)
	mockRoundTrip.AddSlackRequest(
		"/api/usergroups.users.list",
		map[string]string{"token": "TOKEN_HERE", "usergroup": "test"},
		`
{
	"ok": true,
	"users": [
		"fizz"
	]
}
`)
	mockRoundTrip.AddSlackRequest(
		"/api/users.info",
		map[string]string{"users": "fizz", "token": "TOKEN_HERE", "usergroup": "test", "include_locale": "true"},
		`
{
	"ok": true,
	"users": [
		{
			"id": "fizz",
			"profile": {
				"email": "fizz@test.test"
			}
		}
	]
}
`)
	mockRoundTrip.AddSlackRequest(
		"/api/users.lookupByEmail",
		map[string]string{"token": "TOKEN_HERE", "email": "foo@test.test"},
		`
{
	"ok": true,
	"user": {
		"id": "foo",
		"profile": {
			"email": "foo@test.test"
		}
	}
}
`)
	mockRoundTrip.AddSlackRequest(
		"/api/users.lookupByEmail",
		map[string]string{"token": "TOKEN_HERE", "email": "bar@test.test"},
		`
{
	"ok": true,
	"user": {
		"id": "bar",
		"profile": {
			"email": "bar@test.test"
		}
	}
}
`)
	mockRoundTrip.AddSlackRequest(
		"/api/usergroups.users.update",
		map[string]string{"token": "TOKEN_HERE", "usergroup": "test", "users": "fizz,foo,bar"},
		`
{
	"ok": true,
	"usergroup": {
		"date_delete": 0
	}
}
`)

	slackClient := slack.New("TOKEN_HERE", slack.OptionHTTPClient(httpClient))
	repo := slacksync.NewUserGroupService(slackClient, "test")

	success, failure, err := repo.Add("foo@test.test", "bar@test.test")
	assert.ElementsMatch(t, success, []string{"foo@test.test", "bar@test.test"})
	assert.Empty(t, failure)
	assert.NoError(t, err)
}

func TestSlackUserGroup_Remove_NoCache(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	httpClient := httpclient.New(mockRoundTrip)
	mockRoundTrip.AddSlackRequest(
		"/api/usergroups.users.list",
		map[string]string{"token": "TOKEN_HERE", "usergroup": "test"},
		`
{
	"ok": true,
	"users": [
		"foo",
		"bar"
	]
}
`)
	mockRoundTrip.AddSlackRequest(
		"/api/users.info",
		map[string]string{"token": "TOKEN_HERE", "users": "foo,bar", "include_locale": "true"},
		`
{
	"ok": true,
	"users": [
		{
			"id": "foo",
			"profile": {
				"email": "foo@test.test"
			}
		},
		{
			"id": "bar",
			"profile": {
				"email": "bar@test.test"
			}
		}
	]
}
`)
	mockRoundTrip.AddSlackRequest(
		"/api/usergroups.users.update",
		map[string]string{"token": "TOKEN_HERE", "usergroup": "test", "users": "bar"},
		`
{
	"ok": true
}
`)

	slackClient := slack.New("TOKEN_HERE", slack.OptionHTTPClient(httpClient))
	repo := slacksync.NewUserGroupService(slackClient, "test")

	success, failure, err := repo.Remove("foo@test.test")

	assert.ElementsMatch(t, []string{"foo@test.test"}, success)
	assert.Empty(t, failure)
	assert.NoError(t, err)
}
