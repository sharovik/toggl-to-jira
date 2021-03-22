package arguments

import (
	"flag"
	"github.com/sharovik/toggl-jira/log"
	"time"
)

const (
	dateFromHelp    = "The starting date for the filter export. Please use next format: YYYY-MM-DD"
	dateToHelp      = "The ending date for the filter export. Please use next format: YYYY-MM-DD"
	workspaceIDHelp = "The workspace ID which should be used for the toggl.track data report generation. By default will be used the ID from TOGGL_DEFAULT_WORKSPACE_ID environment variable."
)

type OutputArgs struct {
	DateFrom    string
	DateTo      string
	WorkspaceID string
}

func ParseArgs() OutputArgs {
	now := time.Now()

	args := OutputArgs{}
	dateFrom := flag.String("date_from", now.Format("2006-01-02"), dateFromHelp)
	dateTo := flag.String("date_to", now.Format("2006-01-02"), dateToHelp)
	workspaceID := flag.String("workspace_id", "", workspaceIDHelp)

	flag.Parse()

	args.DateFrom = *dateFrom
	args.DateTo = *dateTo
	args.WorkspaceID = *workspaceID
	log.Logger().Info().Interface("args", args).Msg("Received args")
	return args
}
