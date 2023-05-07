package database

import (
	"os"

	"github.com/pkg/errors"
	"github.com/sharovik/orm/clients"
	"github.com/sharovik/toggl-jira/src/config"
	"github.com/sharovik/toggl-jira/src/database/entities"
	"github.com/sharovik/toggl-jira/src/log"
)

var DB clients.BaseClientInterface

func createSQLiteDatabaseFile() error {
	_, err := os.Stat(config.Cfg.Database.Host)
	if err == nil {
		log.Logger().Debug().Msg("Database file already exists")
		return nil
	}

	log.Logger().Info().Msg("Creating the database file")

	_, err = os.Create(config.Cfg.Database.Host)
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to create database file")
		return err
	}

	return nil
}

func InitDB() (err error) {
	if err = createSQLiteDatabaseFile(); err != nil {
		return err
	}

	DB, err = clients.InitClient(config.Cfg.Database)
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to initialise the database client")

		return errors.Wrap(err, "Failed to initialise the database client")
	}

	log.Logger().Debug().Msg("Database client initialised")

	return initTables()
}

func initTables() error {
	q := new(clients.Query).Create(entities.NewHistoryItem()).IfNotExists()
	_, err := DB.Execute(q)
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to create table")

		return err
	}

	return nil
}
