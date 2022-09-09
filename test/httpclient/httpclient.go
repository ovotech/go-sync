package httpclient

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/stretchr/testify/mock"
)

func matchSlackRequest(path string, body map[string]string) func(*http.Request) bool {
	return func(request *http.Request) bool {
		if request.URL.Path != path {
			return false
		}

		bodyString, err := io.ReadAll(request.Body)
		if err != nil {
			return false
		}

		query, err := url.ParseQuery(string(bodyString))
		if err != nil {
			return false
		}

		for k, v := range query {
			if body[k] != v[0] {
				return false
			}
		}

		return true
	}
}

func matchGitHubRequest(method string, path string) func(*http.Request) bool {
	return func(request *http.Request) bool {
		if request.URL.Path != path {
			return false
		}

		if request.Method != method {
			return false
		}

		return true
	}
}

type graphQLRequest struct {
	Query     string
	Variables map[string]interface{}
}

func matchGraphQLRequest(query string, variables map[string]string) func(*http.Request) bool {
	return func(request *http.Request) bool {
		body, err := request.GetBody()
		if err != nil {
			panic("could not get body")
		}

		var requestBody graphQLRequest

		if err := json.NewDecoder(body).Decode(&requestBody); err != nil {
			return false
		}

		if requestBody.Query != query {
			return false
		}

		for k, v := range variables {
			if requestBody.Variables[k] != v {
				return false
			}
		}

		return true
	}
}

type MockRoundTrip struct {
	mock.Mock
}

func (m *MockRoundTrip) RoundTrip(req *http.Request) (*http.Response, error) {
	args := m.Called(req)

	return args.Get(0).(*http.Response), args.Error(1) //nolint:wrapcheck
}

type MockSlackRequest struct{}

// AddSlackRequest adds a mocked request intended for use with the Slack SDK.
func (m *MockRoundTrip) AddSlackRequest(path string, body map[string]string, response string) {
	sr := strings.NewReader(response)
	src := io.NopCloser(sr)
	httpResponse := http.Response{ //nolint:exhaustruct
		StatusCode: http.StatusOK,
		Body:       src,
	}
	m.On("RoundTrip", mock.MatchedBy(matchSlackRequest(path, body))).Return(&httpResponse, nil)
}

// AddGitHubRequest adds a mocked request intended for use with the GitHub SDK.
func (m *MockRoundTrip) AddGitHubRequest(method string, path string, statusCode int, response string) {
	sr := strings.NewReader(response)
	src := io.NopCloser(sr)
	httpResponse := http.Response{ //nolint:exhaustruct
		StatusCode: statusCode,
		Body:       src,
	}
	m.On("RoundTrip", mock.MatchedBy(matchGitHubRequest(method, path))).Return(&httpResponse, nil)
}

func (m *MockRoundTrip) AddGraphQlRequest(query string, variables map[string]string, response string) {
	sr := strings.NewReader(response)
	src := io.NopCloser(sr)
	httpResponse := http.Response{ //nolint:exhaustruct
		StatusCode: http.StatusOK,
		Body:       src,
	}
	m.On("RoundTrip", mock.MatchedBy(matchGraphQLRequest(query, variables))).Return(&httpResponse, nil)
}

func New(r http.RoundTripper) *http.Client {
	return &http.Client{ //nolint:exhaustruct
		Transport: r,
	}
}
