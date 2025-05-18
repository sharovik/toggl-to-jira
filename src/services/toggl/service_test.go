package toggl

import (
	"encoding/json"
	"github.com/stretchr/testify/mock"
	"testing"

	mocks "github.com/sharovik/toggl-jira/mocks/clients"
	"github.com/sharovik/toggl-jira/src/config"
	"github.com/sharovik/toggl-jira/src/dto"
	"github.com/sharovik/toggl-jira/src/log"
	"github.com/sharovik/toggl-jira/src/services/arguments"
	"github.com/stretchr/testify/assert"
)

func init() {
	config.Cfg = config.Config{
		Env:        "testing",
		AppVersion: "1",
	}

	if err := log.Init(config.Cfg); err != nil {
		panic(err)
	}
}

func TestTogglService_GetReport(t *testing.T) {
	m := new(mocks.BaseHTTPClientInterface)
	expectedEndpoint := "/reports/api/v2/details"
	expectedHeaders := map[string]string{
		"Authorization": "Basic Og==",
	}

	args := arguments.OutputArgs{
		WorkspaceID: "1",
		DateFrom:    "2023-01-01",
		DateTo:      "2023-03-01",
	}

	m.On("BasicAuth", config.Cfg.TogglApiToken, "api_token").Return("Og==")

	responseBody := new(dto.TogglDetailsResponse)
	expectedResponse, err := json.Marshal(responseBody)
	assert.NoError(t, err)

	m.On("Get", expectedEndpoint, mock.Anything, expectedHeaders).
		Twice().
		Return(expectedResponse, 200, nil)

	TS = TogglService{Client: m}

	actualResponse, err := TS.GetReport(args)
	assert.NoError(t, err)

	assert.IsType(t, dto.TogglDetailsResponse{}, actualResponse)
}
