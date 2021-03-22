package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sharovik/toggl-jira/log"
	"os"
)

const databaseHost = "./database.sqlite"

var DB *sql.DB

type HistoryItem struct {
	ID       int64  `db:"id"`
	TaskKey  string `db:"task_key"`
	Duration int64  `db:"duration"`
	Added    string `db:"added"`
}

func PrepareDatabase() {
	var err error
	_, err = os.Stat(databaseHost)
	if err == nil {
		log.Logger().Info().Msg("Database file already exists")

		DB, err = sql.Open("sqlite3", databaseHost)
		if err != nil {
			log.Logger().AddError(err).Msg("Failed to open connection")
			return
		}
		return
	}

	log.Logger().Info().Msg("Creating the database file")

	_, err = os.Create(databaseHost)
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to create database file")
		return
	}

	DB, err = sql.Open("sqlite3", databaseHost)
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to open connection")
		return
	}

	_, err = DB.Exec(`
	-- auto-generated definition
	create table history
	(
		id       integer
			constraint history_pk
				primary key autoincrement,
		task_key varchar not null,
		duration integer not null,
		added    varchar not null
	);`)
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to install the data")
		return
	}
}

func InsertTask(taskKey string, duration int64, added string) error {
	_, err := DB.Exec(`insert into history (task_key, duration, added) values ($1, $2, $3)`, taskKey, duration, added)
	if err != nil {
		return err
	}

	log.Logger().Info().
		Str("task_key", taskKey).
		Int64("duration", duration).
		Str("added", added).
		Msg("Inserted the history row")

	return nil
}

func UpdateHistoryItem(item HistoryItem) error {
	_, err := DB.Exec(`update history set duration = $1 where id = $2`, item.Duration, item.ID)
	if err != nil {
		return err
	}

	return nil
}

func FindTask(taskKey string, added string) (item HistoryItem, err error) {
	err = DB.QueryRow("select id, task_key, duration, added from history where task_key = $1 and added = $2 order by id desc limit 1", taskKey, added).
		Scan(&item.ID, &item.TaskKey, &item.Duration, &item.Added)
	if err == sql.ErrNoRows {
		return item, nil
	} else if err != nil {
		return item, err
	}

	return item, nil
}
