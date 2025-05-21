package toggl

import (
	"encoding/json"
	"fmt"
	"net/http"

	client "github.com/sharovik/toggl-jira/src/clients"
	"github.com/sharovik/toggl-jira/src/config"
	"github.com/sharovik/toggl-jira/src/dto"
	"github.com/sharovik/toggl-jira/src/log"
	"github.com/sharovik/toggl-jira/src/services/arguments"
)

var TS TogglServiceInterface

type TogglServiceInterface interface {
	GetReport(args arguments.OutputArgs) (response dto.TogglDetailsResponse, err error)
}

type TogglService struct {
	Client client.BaseHTTPClientInterface
}

func (s TogglService) GetReport(args arguments.OutputArgs) (response dto.TogglDetailsResponse, err error) {
	log.Logger().Info().Interface("args", args).Msg("Start the report GET")
	var (
		query = map[string]string{
			"workspace_id": args.WorkspaceID,
			"user_agent":   config.Cfg.AppVersion,
			"since":        args.DateFrom,
			"until":        args.DateTo,
		}
		nextPage = 2
	)

	response, err = s.getDataForQuery(query)
	if err != nil {
		return dto.TogglDetailsResponse{}, err
	}

	for {
		query["page"] = fmt.Sprintf("%d", nextPage)
		nextPageResponse, err := s.getDataForQuery(query)
		if err != nil {
			return dto.TogglDetailsResponse{}, err
		}

		if len(nextPageResponse.Data) == 0 {
			break
		}

		nextPage++
		response.Data = append(response.Data, nextPageResponse.Data...)
	}

	log.Logger().Info().
		Interface("args", args).
		Interface("response", response).
		Msg("Finish the report GET")

	return
}

func (s TogglService) getDataForQuery(query map[string]string) (response dto.TogglDetailsResponse, err error) {
	log.Logger().Info().Interface("query", query).Msg("Start the report GET")
	var headers = map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", s.Client.BasicAuth(config.Cfg.TogglApiToken, "api_token")),
	}

	result, statusCode, err := s.Client.Get("/reports/api/v2/details", query, headers)
	if statusCode != http.StatusOK {
		log.Logger().Error().
			Int("status_code", statusCode).
			Interface("response", result).
			Msg("Bad response status code received")
		return dto.TogglDetailsResponse{}, fmt.Errorf("bad response status code received")
	}

	if err := json.Unmarshal(result, &response); err != nil {
		log.Logger().AddError(err).Msg("Failed unmarshal the response body")
		return response, err
	}

	log.Logger().Info().
		Interface("query", query).
		Interface("response", response).
		Msg("Finish the report GET")

	return
}
