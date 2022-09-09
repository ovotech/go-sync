package slack_test

import (
	"testing"

	slacksync "github.com/ovotech/go_sync/pkg/slack"
	"github.com/ovotech/go_sync/test/httpclient"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

func TestSlackConversation_Get(t *testing.T) { //nolint:funlen
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	httpClient := httpclient.New(mockRoundTrip)
	mockRoundTrip.AddSlackRequest(
		"/api/conversations.members",
		map[string]string{"channel": "test", "limit": "50", "token": "TOKEN_HERE"},
		`
{
	"ok": true,
	"members": [
		"fizz",
		"buzz"
	],
	"response_metadata": {
		"next_cursor": "page"
	}
}
`)
	mockRoundTrip.AddSlackRequest(
		"/api/conversations.members",
		map[string]string{"channel": "test", "limit": "50", "token": "TOKEN_HERE", "cursor": "page"},
		`
{
	"ok": true,
	"members": [
		"foo",
		"bar"
	]
}
`)
	mockRoundTrip.AddSlackRequest(
		"/api/users.info",
		map[string]string{"users": "fizz,buzz,foo,bar", "token": "TOKEN_HERE", "include_locale": "true"},
		`
{
	"ok": true,
	"users": [
		{
			"id": "fizz",
			"profile": {
				"email": "fizz@test.test"
			}
		},
		{
			"id": "buzz",
			"profile": {
				"email": "buzz@test.test"
			}
		},
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
	repo := slacksync.NewConversationService(slackClient, "test")

	accounts, err := repo.Get()

	assert.ElementsMatch(t, []string{"fizz@test.test", "buzz@test.test", "foo@test.test", "bar@test.test"}, accounts)
	assert.NoError(t, err, "No error should occur.")
}

func TestSlackConversation_Get_ChannelError(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	httpClient := httpclient.New(mockRoundTrip)
	mockRoundTrip.AddSlackRequest(
		"/api/conversations.members",
		map[string]string{"channel": "test", "limit": "50", "token": "TOKEN_HERE"},
		`
{
	"ok": false,
	"error": "test members error"
}
`)

	slackClient := slack.New("TOKEN_HERE", slack.OptionHTTPClient(httpClient))
	repo := slacksync.NewConversationService(slackClient, "test")

	_, err := repo.Get()

	assert.Error(t, err)
	assert.ErrorContains(t, err, "test members error")
}

func TestSlackConversation_Get_UserError(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	httpClient := httpclient.New(mockRoundTrip)
	mockRoundTrip.AddSlackRequest(
		"/api/conversations.members",
		map[string]string{"channel": "test", "limit": "50", "token": "TOKEN_HERE"},
		`
{
	"ok": true,
	"members": [
		"foo",
		"bar"
	]
}
`)
	mockRoundTrip.AddSlackRequest(
		"/api/users.info",
		map[string]string{"include_locale": "true", "token": "TOKEN_HERE", "users": "foo,bar"},
		`
{
	"ok": false,
	"error": "test users error"
}
`)

	slackClient := slack.New("TOKEN_HERE", slack.OptionHTTPClient(httpClient))
	repo := slacksync.NewConversationService(slackClient, "test")

	_, err := repo.Get()

	assert.Error(t, err)
	assert.ErrorContains(t, err, "test users error")
}

func TestSlackConversation_Add(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	httpClient := httpclient.New(mockRoundTrip)
	mockRoundTrip.AddSlackRequest(
		"/api/users.lookupByEmail",
		map[string]string{"email": "foo@test.test", "token": "TOKEN_HERE", "channel": "test"},
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
		map[string]string{"email": "bar@test.test", "token": "TOKEN_HERE", "channel": "test"},
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
		"/api/conversations.invite",
		map[string]string{"channel": "test", "token": "TOKEN_HERE", "users": "foo,bar"},
		`
{
	"ok": true
}
`)

	slackClient := slack.New("TOKEN_HERE", slack.OptionHTTPClient(httpClient))
	repo := slacksync.NewConversationService(slackClient, "test")

	success, failure, err := repo.Add("foo@test.test", "bar@test.test")

	assert.ElementsMatch(t, []string{"foo@test.test", "bar@test.test"}, success)
	assert.Empty(t, failure)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(mockRoundTrip.Calls))
}

func TestSlackConversation_Add_UserNoExist(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	httpClient := httpclient.New(mockRoundTrip)
	mockRoundTrip.AddSlackRequest(
		"/api/users.lookupByEmail",
		map[string]string{"email": "no-exist@test.test", "token": "TOKEN_HERE", "channel": "test"},
		`
{
	"ok": false,
	"error": "lookup error"
}
`)

	slackClient := slack.New("TOKEN_HERE", slack.OptionHTTPClient(httpClient))
	repo := slacksync.NewConversationService(slackClient, "test")

	_, _, err := repo.Add("no-exist@test.test")

	assert.Error(t, err)
	assert.ErrorContains(t, err, "lookup error")
}

func TestSlackConversation_Add_ChannelNoExist(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	httpClient := httpclient.New(mockRoundTrip)
	mockRoundTrip.AddSlackRequest(
		"/api/users.lookupByEmail",
		map[string]string{"email": "foo@test.test", "token": "TOKEN_HERE", "channel": "test"},
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
		"/api/conversations.invite",
		map[string]string{"channel": "test", "token": "TOKEN_HERE", "users": "foo"},
		`
{
	"ok": false,
	"error": "invite fail"
}
`)

	slackClient := slack.New("TOKEN_HERE", slack.OptionHTTPClient(httpClient))
	repo := slacksync.NewConversationService(slackClient, "test")

	_, _, err := repo.Add("foo@test.test")

	assert.Error(t, err)
	assert.ErrorContains(t, err, "invite fail")
}

func TestSlackConversation_RemoveAccount_NoCache(t *testing.T) { //nolint:funlen
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	httpClient := httpclient.New(mockRoundTrip)
	mockRoundTrip.AddSlackRequest(
		"/api/conversations.members",
		map[string]string{"channel": "test", "limit": "50", "token": "TOKEN_HERE"},
		`
{
	"ok": true,
	"members": [
		"foo",
		"bar",
		"fizz",
		"buzz"
	]
}
`)
	mockRoundTrip.AddSlackRequest(
		"/api/users.info",
		map[string]string{"users": "foo,bar,fizz,buzz", "token": "TOKEN_HERE", "include_locale": "true"},
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
		},
		{
			"id": "fizz",
			"profile": {
				"email": "fizz@test.test"
			}
		},
		{
			"id": "buzz",
			"profile": {
				"email": "buzz@test.test"
			}
		}
	]
}
`)
	mockRoundTrip.AddSlackRequest(
		"/api/conversations.kick",
		map[string]string{"channel": "test", "token": "TOKEN_HERE", "user": "foo"},
		`
{
	"ok": true
}
`)
	mockRoundTrip.AddSlackRequest(
		"/api/conversations.kick",
		map[string]string{"channel": "test", "token": "TOKEN_HERE", "user": "bar"},
		`
{
	"ok": true
}
`)

	slackClient := slack.New("TOKEN_HERE", slack.OptionHTTPClient(httpClient))
	repo := slacksync.NewConversationService(slackClient, "test")

	success, failure, err := repo.Remove("foo@test.test", "bar@test.test")

	assert.ElementsMatch(t, []string{"foo@test.test", "bar@test.test"}, success)
	assert.Empty(t, failure)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(mockRoundTrip.Calls))
}
