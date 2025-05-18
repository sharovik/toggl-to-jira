package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	client "github.com/sharovik/toggl-jira/src/clients"

	"github.com/sharovik/toggl-jira/src/config"
	"github.com/sharovik/toggl-jira/src/database"
	"github.com/sharovik/toggl-jira/src/dto"
	"github.com/sharovik/toggl-jira/src/log"
	"github.com/sharovik/toggl-jira/src/services/arguments"
	"github.com/sharovik/toggl-jira/src/services/jira"
	"github.com/sharovik/toggl-jira/src/services/toggl"
)

var Cfg config.Config

func init() {
	Cfg = config.Init()
	if err := log.Init(Cfg); err != nil {
		panic(err)
	}

	if err := database.InitDB(); err != nil {
		panic(err)
	}

	httpClient := client.GetHttpClient()
	jira.JS = jira.JiraService{
		Client: &client.HTTPClient{
			Client:  &httpClient,
			BaseURL: config.Cfg.JiraBaseURL,
		},
	}

	httpClient = client.GetHttpClient()
	toggl.TS = toggl.TogglService{
		Client: &client.HTTPClient{
			Client:  &httpClient,
			BaseURL: config.Cfg.TogglApiURL,
		},
	}
}

func main() {
	args := arguments.ParseArgs()

	if err := validateConfiguration(); err != nil {
		log.Logger().AddError(err).Msg("It looks like there are problems with config. Stop.")
		os.Exit(1)
	}

	if args.WorkspaceID == "" {
		args.WorkspaceID = config.Get().TogglWorkspaceID
	}

	report, err := toggl.TS.GetReport(args)
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to retrieve the export.")
		os.Exit(1)
	}

	if len(report.Data) == 0 {
		log.Logger().Info().Msg("No data in report.")
		return
	}

	for _, item := range report.Data {
		taskKey, err := jira.JS.ParseTaskKey(item.Description)
		if err != nil {
			log.Logger().Warn().Err(err).Msg("Failed to parse the task key.")
			continue
		}

		if taskKey == "" {
			log.Logger().Info().
				Str("description", item.Description).
				Int64("duration", item.Dur).
				Str("added", item.Start.Format("2006-01-02")).
				Msg("There is no task key specified in description.")
			continue
		}

		if err = database.BeginTransaction(); err != nil {
			log.Logger().AddError(err).Msg("Failed to begin transaction")

			return
		}

		timeEntry, err := InsertHistoryScenario(taskKey, item)
		if err != nil {
			if err.Error() != "item already processed" {
				log.Logger().AddError(err).Msg("Failed to insert time entry")
			}

			if e := database.RollBackTransaction(); e != nil {
				log.Logger().AddError(e).Msg("Failed to rollback transaction")

				return
			}
			continue
		}

		spentMinutes := int64(timeEntry.Minutes())
		if spentMinutes == 0 {
			log.Logger().Info().
				Str("task_key", taskKey).
				Int64("spent_in_minutes", spentMinutes).
				Str("spent_total_time", timeEntry.String()).
				Msg("Nothing to track.")

			if e := database.RollBackTransaction(); e != nil {
				log.Logger().AddError(e).Msg("Failed to rollback transaction")

				return
			}

			continue
		}

		if err = jira.JS.SendTheTime(taskKey, fmt.Sprintf("%dm", spentMinutes), item.Start); err != nil {
			log.Logger().AddError(err).Msg("Failed to send the worklog")
			if e := database.RollBackTransaction(); e != nil {
				log.Logger().AddError(e).Msg("Failed to rollback transaction")

				return
			}
		}

		if e := database.CommitTransaction(); e != nil {
			log.Logger().AddError(e).Msg("Failed to rollback transaction")

			return
		}

		time.Sleep(time.Duration(1) * time.Second)
	}
}

func InsertHistoryScenario(taskKey string, item dto.DataItem) (timeEntry time.Duration, err error) {
	timeEntry = time.Duration(item.Dur) * time.Millisecond

	historyItem, err := database.FindTask(taskKey, item.Dur, item.Start.Format("2006-01-02"))
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to retrieve information from history table")
		return 0, err
	}

	if historyItem != nil {
		log.Logger().Info().
			Str("task_key", taskKey).
			Int64("duration", item.Dur).
			Str("added", item.Start.Format("2006-01-02")).
			Msg("This item was already processed. Ignoring.")

		return 0, errors.New("item already processed")
	}

	log.Logger().Info().
		Str("task_key", taskKey).
		Int64("duration", item.Dur).
		Str("added", item.Start.Format("2006-01-02")).
		Msg("Inserting new history row")

	if err := database.InsertTask(taskKey, item.Dur, item.Start.Format("2006-01-02")); err != nil {
		log.Logger().AddError(err).Msg("Failed to insert the history item!")
		return timeEntry, err
	}

	return timeEntry, nil
}

func validateConfiguration() error {
	if config.Get().JiraBaseURL == "" ||
		config.Get().JiraAppToken == "" ||
		config.Get().JiraEmail == "" ||
		config.Get().TogglApiToken == "" ||
		config.Get().TogglApiURL == "" ||
		config.Get().TogglWorkspaceID == "" {
		return errors.New("One of the required items of configuration is missing. ")
	}

	return nil
}
