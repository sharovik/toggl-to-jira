package toggl

import (
	"encoding/json"
	"fmt"
	client "github.com/sharovik/toggl-jira/clients"
	"github.com/sharovik/toggl-jira/config"
	"github.com/sharovik/toggl-jira/dto"
	"github.com/sharovik/toggl-jira/log"
	"github.com/sharovik/toggl-jira/services/arguments"
	"net/http"
)

func GetReport(args arguments.OutputArgs) (response dto.TogglDetailsResponse, err error) {
	log.Logger().Info().Interface("args", args).Msg("Start the report GET")
	var (
		httpClient = client.GetHttpClient()
		c          = client.HTTPClient{
			Client:  &httpClient,
			BaseURL: config.Get().TogglApiURL,
		}
		query = map[string]string{
			"workspace_id": args.WorkspaceID,
			"user_agent":   config.Get().AppVersion,
			"since":        args.DateFrom,
			"until":        args.DateTo,
		}
		headers = map[string]string{
			"Authorization": c.BasicAuth(config.Get().TogglApiToken, "api_token"),
		}
	)

	result, statusCode, err := c.Get("/reports/api/v2/details", query, headers)
	if statusCode != http.StatusOK {
		log.Logger().Error().
			Int("status_code", statusCode).
			Interface("response", result).
			Msg("Bad response status code received")
		return dto.TogglDetailsResponse{}, fmt.Errorf("Bad response status code received ")
	}

	if err := json.Unmarshal(result, &response); err != nil {
		log.Logger().AddError(err).Msg("Failed unmarshal the response body")
		return response, err
	}

	log.Logger().Info().
		Interface("args", args).
		Interface("response", response).
		Msg("Finish the report GET")
	return
}
