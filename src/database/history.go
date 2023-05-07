package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sharovik/orm/clients"
	"github.com/sharovik/orm/dto"
	"github.com/sharovik/orm/query"
	"github.com/sharovik/toggl-jira/src/database/entities"
	"github.com/sharovik/toggl-jira/src/log"
)

func InsertTask(taskKey string, duration int64, added string) error {
	item := entities.NewHistoryItem()
	item.UpdateFieldValue("task_key", taskKey)
	item.UpdateFieldValue("duration", duration)
	item.UpdateFieldValue("added", added)

	//We insert new item into our table
	q := new(clients.Query).Insert(item)
	_, err := DB.Execute(q)
	if err != nil {
		return err
	}

	log.Logger().Debug().
		Str("task_key", taskKey).
		Int64("duration", duration).
		Str("added", added).
		Msg("Inserted the history row")

	return nil
}

func FindTask(taskKey string, duration int64, added string) (item dto.ModelInterface, err error) {
	q := new(clients.Query).
		Select([]interface{}{}).
		From(entities.NewHistoryItem()).
		Where(query.Where{
			First:    "task_key",
			Operator: "=",
			Second: query.Bind{
				Field: "task_key",
				Value: taskKey,
			},
		}).
		Where(query.Where{
			First:    "duration",
			Operator: "=",
			Second: query.Bind{
				Field: "duration",
				Value: duration,
			},
		}).
		Where(query.Where{
			First:    "added",
			Operator: "=",
			Second: query.Bind{
				Field: "added",
				Value: added,
			},
		}).
		OrderBy("id", query.OrderDirectionDesc).
		Limit(query.Limit{
			From: 0,
			To:   1,
		})

	res, err := DB.Execute(q)
	if err == sql.ErrNoRows {
		return item, nil
	} else if err != nil {
		return item, err
	}

	if len(res.Items()) == 0 {
		return item, nil
	}

	//We take first item and use it as the result
	return res.Items()[0], nil
}
