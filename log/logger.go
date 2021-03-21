package log

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var loggerInstance LoggerInstance

// Config configuration required by Logger
type Config interface {
	GetAppEnv() string
}

// LoggerInstance shared state
type LoggerInstance struct {
	env         string
	context     map[string]interface{}
	initialized bool
}

//default constants
const (
	FieldContext      = "context"
	FieldLevelName    = "level_name"
	FieldErrorMessage = "error_message"

	appEnvDevelopment = "development"
	appEnvTesting     = "testing"
)

// Init initializes the logger (required before use)
func Init(config Config) error {
	loggerInstance = LoggerInstance{
		env:         config.GetAppEnv(),
		context:     make(map[string]interface{}),
		initialized: true,
	}

	return nil
}

//Refresh refreshes the logger instance
func Refresh() {
	loggerInstance = LoggerInstance{}
}

//IsInitialized function retrieves current status of logger instance
func IsInitialized() bool {
	return loggerInstance.initialized
}

//Logger returns a pointer to the singleton Logger loggerInstance
func Logger() *LoggerInstance {
	if !loggerInstance.initialized {
		panic("logger not initialized")
	}
	return &loggerInstance
}

//AppendGlobalContext for setting global context
func (l *LoggerInstance) AppendGlobalContext(context map[string]interface{}) {
	if l.context == nil {
		l.context = context
	}

	for field, value := range context {
		l.context[field] = value
	}

	l.Debug().Interface("context_changes", context).Msg("Append new global context")
}

//GlobalContext method retrieve the GlobalContext variable
func (l *LoggerInstance) GlobalContext() map[string]interface{} {
	return l.context
}

//DestroyGlobalContext method for global context destroy
func (l *LoggerInstance) DestroyGlobalContext() {
	l.context = make(map[string]interface{})
}

//AddError for correct error messages parse
func (l *LoggerInstance) AddError(err error) *zerolog.Event {
	err = errors.Wrap(err, err.Error())
	return l.Error().Stack().Err(err)
}

//DefaultContext method which returns Logger with default context
func (l *LoggerInstance) DefaultContext() *zerolog.Logger {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	var context = zerolog.Context{}
	if l.env == appEnvTesting {
		context = log.Output(ioutil.Discard).With()
	} else if l.env == appEnvDevelopment {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		context = log.Output(zerolog.NewConsoleWriter()).With()
	} else {
		context = log.With()
	}

	zerolog.TimestampFieldName = "@timestamp"
	zerolog.LevelFieldName = FieldLevelName
	zerolog.ErrorFieldName = FieldErrorMessage
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000000"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	logger := context.
		Interface(FieldContext, l.context).
		Logger()

	return &logger
}

//Debug method for messages with level DEBUG
func (l *LoggerInstance) Debug() *zerolog.Event {
	return l.DefaultContext().Debug()
}

//Info method for messages with level INFO
func (l *LoggerInstance) Info() *zerolog.Event {
	return l.DefaultContext().Info()
}

//Error method for messages with level ERROR
func (l *LoggerInstance) Error() *zerolog.Event {
	return l.DefaultContext().Error()
}

//Warn method for messages with level WARNING
func (l *LoggerInstance) Warn() *zerolog.Event {
	return l.DefaultContext().Warn()
}

//StartMessage adds message with START postfix
func (l *LoggerInstance) StartMessage(msg string) {
	l.DefaultContext().Info().Msg(fmt.Sprintf("%s: %s", msg, "START"))
}

//FinishMessage adds message with FINISH postfix
func (l *LoggerInstance) FinishMessage(msg string) {
	l.DefaultContext().Info().Msg(fmt.Sprintf("%s: %s", msg, "FINISH"))
}
