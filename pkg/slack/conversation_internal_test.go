package slack

import (
	"testing"

	"github.com/ovotech/go_sync/test/httpclient"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

func TestNewConversationService(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	httpClient := httpclient.New(mockRoundTrip)
	slackClient := slack.New("TOKEN_HERE", slack.OptionHTTPClient(httpClient))

	repo := NewConversationService(slackClient, "test")

	assert.IsType(t, new(Conversation), repo)
	assert.Equal(t, slackClient, repo.client)
	assert.Equal(t, "test", repo.conversationName)
	assert.Empty(t, repo.cache)
	assert.Zero(t, mockRoundTrip.Calls, "No calls must be made when only using a constructor.")
}

func TestSlackConversation_Remove_WithCache(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	httpClient := httpclient.New(mockRoundTrip)
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
	repo := NewConversationService(slackClient, "test")
	repo.cache = map[string]string{"foo@test.test": "foo", "bar@test.test": "bar"}

	success, failure, err := repo.Remove("foo@test.test", "bar@test.test")

	assert.ElementsMatch(t, []string{"foo@test.test", "bar@test.test"}, success)
	assert.Empty(t, failure)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(mockRoundTrip.Calls))
}

func TestSlackConversation_Remove_ChannelNoExist(t *testing.T) {
	t.Parallel()

	mockRoundTrip := new(httpclient.MockRoundTrip)
	httpClient := httpclient.New(mockRoundTrip)
	mockRoundTrip.AddSlackRequest(
		"/api/conversations.kick",
		map[string]string{"channel": "test", "token": "TOKEN_HERE", "user": "foo"},
		`
{
	"ok": false,
	"error": "kick fail"
}
`)

	slackClient := slack.New("TOKEN_HERE", slack.OptionHTTPClient(httpClient))
	repo := NewConversationService(slackClient, "test")
	repo.cache = map[string]string{"foo@test.test": "foo"}

	_, failure, err := repo.Remove("foo@test.test")

	assert.NotEmpty(t, failure)
	assert.Len(t, failure, 1)
	assert.ErrorContains(t, failure[0], "kick fail")
	assert.NoError(t, err)
}
