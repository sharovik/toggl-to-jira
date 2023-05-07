package jira

import (
	"encoding/json"
	"testing"
	"time"

	mocks "github.com/sharovik/toggl-jira/mocks/clients"
	"github.com/sharovik/toggl-jira/src/config"
	"github.com/sharovik/toggl-jira/src/dto"
	"github.com/sharovik/toggl-jira/src/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	config.Cfg = config.Config{
		Env:          "testing",
		JiraBaseURL:  "http://test.dev/",
		JiraEmail:    "some@email.dv",
		JiraAppToken: "token",
	}

	if err := log.Init(config.Cfg); err != nil {
		panic(err)
	}
}

func TestJiraService_SendTheTime(t *testing.T) {
	m := new(mocks.BaseHTTPClientInterface)

	date := time.Now()
	expectedEndpoint := "/rest/internal/3/issue/test/worklog?adjustEstimate=auto"
	expectedHeaders := map[string]string{
		"Authorization": "Basic Og==",
	}

	body := dto.JiraWorklogRequest{
		Timespent: "1m",
		Comment: dto.Comment{
			Type:    "doc",
			Version: 1,
			Content: []string{},
		},
		Started: date.Format(jiraTimeFormat),
	}

	expectedRequest, err := json.Marshal(body)
	assert.NoError(t, err)

	responseBody := new(dto.JiraWorkLogResponse)

	expectedResponse, err := json.Marshal(responseBody)
	assert.NoError(t, err)

	m.On("BasicAuth", config.Cfg.JiraEmail, config.Cfg.JiraAppToken).Return("Og==")

	m.On("Post", expectedEndpoint, expectedRequest, expectedHeaders).
		Once().
		Return(expectedResponse, 201, nil)

	JS = JiraService{Client: m}

	err = JS.SendTheTime("test", "1m", date)
	assert.NoError(t, err)
}

func TestJiraService_ParseTaskKey(t *testing.T) {
	cases := map[string]string{
		"SOME-323 test text": "SOME-323",
		"test test":          "",
	}

	JS = JiraService{}
	for text, expected := range cases {
		actual, err := JS.ParseTaskKey(text)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	}
}
