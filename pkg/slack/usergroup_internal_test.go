package slack

import (
	"testing"

	"github.com/ovotech/go_sync/test/httpclient"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	httpClient := httpclient.New(mockRoundTrip)
	slackClient := slack.New("TOKEN_HERE", slack.OptionHTTPClient(httpClient))

	repo := NewUserGroupService(slackClient, "test")

	assert.IsType(t, new(UserGroup), repo)
	assert.Zero(t, mockRoundTrip.Calls)
}

func TestSlackUserGroup_Add_Cache(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	httpClient := httpclient.New(mockRoundTrip)
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
		map[string]string{"token": "TOKEN_HERE", "usergroup": "test", "users": "foo,bar"},
		`
{
	"ok": true
}
`)

	slackClient := slack.New("TOKEN_HERE", slack.OptionHTTPClient(httpClient))
	repo := NewUserGroupService(slackClient, "test")
	repo.cache = map[string]string{"foo@test.test": "foo"}

	success, failure, err := repo.Add("bar@test.test")

	assert.ElementsMatch(t, []string{"bar@test.test"}, success)
	assert.Empty(t, failure)
	assert.NoError(t, err)
}

func TestSlackUserGroup_Add_ReEnableGroup(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	httpClient := httpclient.New(mockRoundTrip)
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
		map[string]string{"token": "TOKEN_HERE", "usergroup": "test", "users": "foo,bar"},
		`
{
	"ok": true,
	"usergroup": {
		"date_delete": 1
	}
}
`)
	mockRoundTrip.AddSlackRequest(
		"/api/usergroups.enable",
		map[string]string{"token": "TOKEN_HERE", "usergroup": "test"},
		`
{
	"ok": true
}
`)

	slackClient := slack.New("TOKEN_HERE", slack.OptionHTTPClient(httpClient))
	repo := NewUserGroupService(slackClient, "test")

	// Create a fake cache as though Get has been called previously.
	repo.cache = map[string]string{"foo@test.bar": "foo"}

	success, failure, err := repo.Add("bar@test.test")

	assert.ElementsMatch(t, []string{"bar@test.test"}, success)
	assert.Empty(t, failure)
	assert.NoError(t, err)
}

func TestSlackUserGroup_Add_LookupEmailFailure(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	httpClient := httpclient.New(mockRoundTrip)
	mockRoundTrip.AddSlackRequest(
		"/api/users.lookupByEmail",
		map[string]string{"token": "TOKEN_HERE", "email": "bar@test.test"},
		`
{
	"ok": false,
	"error": "users_not_found"
}
`)

	slackClient := slack.New("TOKEN_HERE", slack.OptionHTTPClient(httpClient))
	repo := NewUserGroupService(slackClient, "test")

	// Create a fake cache as though Get has been called previously.
	repo.cache = map[string]string{"foo@test.bar": "foo"}

	_, _, err := repo.Add("bar@test.test")

	assert.Error(t, err)
	assert.NotNil(t, err, "An error should be passed up from the lookup call.")
	assert.ErrorContains(t, err, "users_not_found")
}

func TestSlackUserGroup_Add_UpdateUserGroupsFailure(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	httpClient := httpclient.New(mockRoundTrip)
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
		map[string]string{"token": "TOKEN_HERE", "usergroup": "test", "users": "foo,bar"},
		`
{
	"ok": false,
	"error": "invalid_users"
}
`)

	slackClient := slack.New("TOKEN_HERE", slack.OptionHTTPClient(httpClient))
	repo := NewUserGroupService(slackClient, "test")

	// Create a fake cache as though Get has been called previously.
	repo.cache = map[string]string{"foo@test.bar": "foo"}

	_, _, err := repo.Add("bar@test.test")

	assert.Error(t, err)
	assert.NotNil(t, err, "An error should be passed up from the lookup call.")
	assert.ErrorContains(t, err, "invalid_users")
}

func TestSlackUserGroup_Remove_WithCache(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	httpClient := httpclient.New(mockRoundTrip)
	mockRoundTrip.AddSlackRequest(
		"/api/usergroups.users.update",
		map[string]string{"token": "TOKEN_HERE", "usergroup": "test", "users": "bar"},
		`
{
	"ok": true
}
`)

	slackClient := slack.New("TOKEN_HERE", slack.OptionHTTPClient(httpClient))
	repo := NewUserGroupService(slackClient, "test")
	repo.cache = map[string]string{"foo@test.test": "foo", "bar@test.test": "bar"}

	success, failure, err := repo.Remove("foo@test.test")

	assert.ElementsMatch(t, []string{"foo@test.test"}, success)
	assert.Empty(t, failure)
	assert.NoError(t, err)
}

func TestSlackUserGroup_Remove_DisableGroup(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	httpClient := httpclient.New(mockRoundTrip)
	mockRoundTrip.AddSlackRequest(
		"/api/usergroups.disable",
		map[string]string{"token": "TOKEN_HERE", "usergroup": "test"},
		`
{
	"ok": true
}
`)

	slackClient := slack.New("TOKEN_HERE", slack.OptionHTTPClient(httpClient))
	repo := NewUserGroupService(slackClient, "test")
	repo.cache = map[string]string{"foo@test.test": "foo"}

	success, failure, err := repo.Remove("foo@test.test")

	assert.ElementsMatch(t, []string{"foo@test.test"}, success)
	assert.Empty(t, failure)
	assert.NoError(t, err)
}

func TestSlackUserGroup_Cache(t *testing.T) { //nolint:funlen
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
	"ok": true,
	"users": [
		{
			"id": "foo",
			"profile": {
				"email": "foo@test.test"
			}
		}
	]
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
		map[string]string{"token": "TOKEN_HERE", "usergroup": "test", "users": "foo,bar"},
		`
{
	"ok": true
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
	repo := NewUserGroupService(slackClient, "test")

	// Assert cache is empty.
	assert.Equal(t, map[string]string{}, repo.cache)

	var (
		success []string
		failure []error
		err     error
	)

	_, err = repo.Get()

	assert.NoError(t, err)

	// Assert cache contains foo@test.test.
	assert.Equal(t, map[string]string{"foo@test.test": "foo"}, repo.cache)

	success, failure, err = repo.Add("bar@test.test")

	assert.ElementsMatch(t, []string{"bar@test.test"}, success)
	assert.Empty(t, failure)
	assert.NoError(t, err)

	// Assert cache has been updated with new bar@test.test.
	assert.Equal(t, map[string]string{"foo@test.test": "foo", "bar@test.test": "bar"}, repo.cache)

	success, failure, err = repo.Remove("foo@test.test")

	assert.ElementsMatch(t, []string{"foo@test.test"}, success)
	assert.Empty(t, failure)
	assert.NoError(t, err)

	// Assert cache no longer has foo@test.test.
	assert.Equal(t, map[string]string{"bar@test.test": "bar"}, repo.cache)
}
