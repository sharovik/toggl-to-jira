package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	client "github.com/sharovik/toggl-jira/src/clients"
	"github.com/sharovik/toggl-jira/src/config"
	"github.com/sharovik/toggl-jira/src/dto"
	"github.com/sharovik/toggl-jira/src/log"
)

const (
	jiraTimeFormat = `2006-01-02T15:04:05.49Z0700`
	jiraTaskRegex  = `([A-Z0-9]+-(\d+))`
)

var JS JiraServiceInterface

type JiraServiceInterface interface {
	ParseTaskKey(text string) (key string, err error)
	SendTheTime(taskKey string, timeStr string, date time.Time) error
}

type JiraService struct {
	Client client.BaseHTTPClientInterface
}

func (s JiraService) ParseTaskKey(text string) (key string, err error) {
	re, err := regexp.Compile(jiraTaskRegex)
	if err != nil {
		return "", err
	}

	matches := re.FindStringSubmatch(text)

	if len(matches) == 0 || matches[1] == "" {
		return "", err
	}

	key = matches[1]
	return
}

func (s JiraService) SendTheTime(taskKey string, timeStr string, date time.Time) error {
	request := dto.JiraWorklogRequest{
		Timespent: timeStr,
		Comment: dto.Comment{
			Type:    "doc",
			Version: 1,
			Content: []string{},
		},
		Started: date.Format(jiraTimeFormat),
	}
	log.Logger().Info().
		Str("task_key", taskKey).
		Interface("request", request).
		Str("start", request.Started).
		Msg("Sending the worklog to Jira.")

	var (
		headers = map[string]string{
			"Authorization": fmt.Sprintf("Basic %s", s.Client.BasicAuth(config.Cfg.JiraEmail, config.Cfg.JiraAppToken)),
		}
	)

	byteRequest, err := json.Marshal(request)
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to marshal the request struct")
		return err
	}

	requestURL := fmt.Sprintf("/rest/internal/3/issue/%s/worklog?adjustEstimate=auto", taskKey)
	result, statusCode, err := s.Client.Post(requestURL, byteRequest, headers)
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to send the worklog")
		return err
	}

	if statusCode != http.StatusCreated {
		log.Logger().Error().
			Int("status_code", statusCode).
			Interface("response", result).
			Msg("Bad response status code received")
		return fmt.Errorf("bad response status code received")
	}

	response := dto.JiraWorkLogResponse{}
	if err := json.Unmarshal(result, &response); err != nil {
		log.Logger().AddError(err).Msg("Failed unmarshal the response body")
		return err
	}

	log.Logger().Info().
		Str("task_key", taskKey).
		Interface("response", response).
		Msg("Worklog was sent")

	return nil
}
