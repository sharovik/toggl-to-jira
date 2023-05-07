package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sharovik/orm/clients"
)

const (
	defaultEnvFilePath = "./.env"

	// database environment variables definition
	envDatabaseHost       = "DATABASE_HOST"
	envDatabaseUsername   = "DATABASE_USERNAME"
	envDatabasePassword   = "DATABASE_PASSWORD"
	envDatabaseName       = "DATABASE_NAME"
	envDatabaseConnection = "DATABASE_CONNECTION"

	defaultDatabaseHost       = "database.sqlite"
	defaultDatabaseConnection = "sqlite"
)

type Config struct {
	Initialized      bool
	Env              string
	TogglApiToken    string
	TogglApiURL      string
	TogglWorkspaceID string
	JiraAppToken     string
	JiraEmail        string
	JiraBaseURL      string
	AppVersion       string

	Database clients.DatabaseConfig
}

var Cfg Config

func Get() Config {
	if !Cfg.Initialized {
		return Init()
	}

	return Cfg
}

func (c Config) GetAppEnv() string {
	return c.Env
}

func Init() Config {
	envPath := defaultEnvFilePath
	if _, err := os.Stat(envPath); err != nil {
		panic(err)
	}

	if err := godotenv.Load(envPath); err != nil {
		panic(err)
	}

	Cfg = Config{
		Env:              os.Getenv("APP_ENV"),
		TogglApiToken:    os.Getenv("TOGGL_API_TOKEN"),
		TogglApiURL:      os.Getenv("TOGGL_API_URL"),
		TogglWorkspaceID: os.Getenv("TOGGL_DEFAULT_WORKSPACE_ID"),
		AppVersion:       os.Getenv("APP_VERSION"),
		JiraAppToken:     os.Getenv("JIRA_APP_TOKEN"),
		JiraBaseURL:      os.Getenv("JIRA_BASE_URL"),
		JiraEmail:        os.Getenv("JIRA_EMAIL"),
		Initialized:      true,
		Database:         initDBCfg(),
	}

	return Cfg
}

func initDBCfg() clients.DatabaseConfig {
	dbConnection := clients.DatabaseTypeSqlite
	if os.Getenv(envDatabaseConnection) != "" {
		dbConnection = os.Getenv(envDatabaseConnection)
	}

	host := defaultDatabaseHost
	if os.Getenv(envDatabaseHost) != "" {
		host = os.Getenv(envDatabaseHost)
	}

	return clients.DatabaseConfig{
		Type:     dbConnection,
		Host:     host,
		Username: os.Getenv(envDatabaseUsername),
		Password: os.Getenv(envDatabasePassword),
		Database: os.Getenv(envDatabaseName),
	}
}
