package jira

import (
	"encoding/json"
	"fmt"
	client "github.com/sharovik/toggl-jira/clients"
	"github.com/sharovik/toggl-jira/config"
	"github.com/sharovik/toggl-jira/dto"
	"github.com/sharovik/toggl-jira/log"
	"net/http"
	"regexp"
	"time"
)

const jiraTaskRegex = `([A-Z]+-(\d+))`

func ParseTaskKey(text string) (key string, err error) {
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

func SendTheTime(taskKey string, timeStr string, date time.Time) error {
	request := dto.JiraWorklogRequest{
		Timespent: timeStr,
		Comment: dto.Comment{
			Type:    "doc",
			Version: 1,
			Content: []string{},
		},
		Started: date.Format("2006-01-02T15:04:05.49Z0700"),
	}
	log.Logger().Info().
		Str("task_key", taskKey).
		Interface("request", request).
		Str("start", request.Started).
		Msg("Sending the worklog to Jira.")

	var (
		httpClient = client.GetHttpClient()
		c          = client.HTTPClient{
			Client:  &httpClient,
			BaseURL: config.Get().JiraBaseURL,
		}
		headers = map[string]string{
			"Authorization": fmt.Sprintf("Basic %s", c.BasicAuth(config.Get().JiraEmail, config.Get().JiraAppToken)),
		}
	)

	byteRequest, err := json.Marshal(request)
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to marshal the request struct")
		return err
	}

	requestURL := fmt.Sprintf("/rest/internal/3/issue/%s/worklog?adjustEstimate=auto", taskKey)
	result, statusCode, err := c.Post(requestURL, byteRequest, headers)
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to send the worklog")
		return err
	}

	if statusCode != http.StatusCreated {
		log.Logger().Error().
			Int("status_code", statusCode).
			Interface("response", result).
			Msg("Bad response status code received")
		return fmt.Errorf("Bad response status code received ")
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
