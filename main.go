package main

import (
	"fmt"
	"github.com/sharovik/toggl-jira/config"
	"github.com/sharovik/toggl-jira/database"
	"github.com/sharovik/toggl-jira/dto"
	"github.com/sharovik/toggl-jira/log"
	"github.com/sharovik/toggl-jira/services/arguments"
	"github.com/sharovik/toggl-jira/services/jira"
	"github.com/sharovik/toggl-jira/services/toggl"
	"time"
)

var Cfg config.Config

func init() {
	Cfg = config.Init()
	if err := log.Init(Cfg); err != nil {
		panic(err)
	}

	database.PrepareDatabase()
}

func main() {
	args := arguments.ParseArgs()
	if args.WorkspaceID == "" {
		args.WorkspaceID = config.Get().TogglWorkspaceID
	}

	report, err := toggl.GetReport(args)
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to retrieve the export.")
		return
	}

	if len(report.Data) == 0 {
		log.Logger().Info().Msg("No data in report.")
		return
	}

	for _, item := range report.Data {
		taskKey, err := jira.ParseTaskKey(item.Description)
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

		timeEntry, err := UpSertHistoryScenario(taskKey, item)
		if err != nil {
			log.Logger().AddError(err).Msg("Failed to execute InsertHistoryScenario")
			continue
		}

		spentMinutes := int64(timeEntry.Minutes())
		if spentMinutes == 0 {
			log.Logger().Info().
				Str("task_key", taskKey).
				Int64("spent_in_minutes", spentMinutes).
				Str("spent_total_time", timeEntry.String()).
				Msg("Nothing to track.")
			continue
		}

		if err := jira.SendTheTime(taskKey, fmt.Sprintf("%dm", spentMinutes), item.Start); err != nil {
			log.Logger().AddError(err).Msg("Failed to send the worklog")
		}

		time.Sleep(time.Duration(1) * time.Second)
	}
}

func UpSertHistoryScenario(taskKey string, item dto.DataItem) (timeEntry time.Duration, err error) {
	timeEntry = time.Duration(item.Dur) * time.Millisecond

	historyItem, err := database.FindTask(taskKey, item.Start.Format("2006-01-02"))
	if historyItem == (database.HistoryItem{}) {
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

	if item.Dur == historyItem.Duration {
		log.Logger().Info().
			Str("task_key", taskKey).
			Int64("duration", item.Dur).
			Str("added", item.Start.Format("2006-01-02")).
			Msg("This item was already processed. Ignoring.")

		return 0, nil
	}

	log.Logger().Info().
		Str("task_key", taskKey).
		Int64("duration", item.Dur).
		Str("added", item.Start.Format("2006-01-02")).
		Msg("Updating the existing history row")

	timeEntry = time.Duration(item.Dur) * time.Millisecond

	historyItem.Duration = historyItem.Duration + item.Dur
	if err := database.UpdateHistoryItem(historyItem); err != nil {
		log.Logger().AddError(err).Msg("Failed to update the history item!")
		return timeEntry, err
	}

	return timeEntry, nil
}
