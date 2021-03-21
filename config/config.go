package config

import (
	"github.com/joho/godotenv"
	"os"
)

const defaultEnvFilePath = "./.env"

type Config struct {
	Initialized bool
	Env string
	TogglApiToken string
	TogglApiURL string
	TogglWorkspaceID string
	JiraAppToken string
	JiraEmail string
	JiraBaseURL string
	AppVersion string
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
		Env: os.Getenv("APP_ENV"),
		TogglApiToken: os.Getenv("TOGGL_API_TOKEN"),
		TogglApiURL: os.Getenv("TOGGL_API_URL"),
		TogglWorkspaceID: os.Getenv("TOGGL_DEFAULT_WORKSPACE_ID"),
		AppVersion: os.Getenv("APP_VERSION"),
		JiraAppToken: os.Getenv("JIRA_APP_TOKEN"),
		JiraBaseURL: os.Getenv("JIRA_BASE_URL"),
		JiraEmail: os.Getenv("JIRA_EMAIL"),
		Initialized: true,
	}

	return Cfg
}
