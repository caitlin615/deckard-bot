package log

import (
	"github.com/handwritingio/deckard-bot/config"

	"github.com/evalphobia/logrus_sentry"
	"github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
)

func init() {
	if config.RuntimeEnv != "development" {
		addSentryHook()
	}
}

var sentryClient *raven.Client

// GetSentryClient returns the Sentry client for the application, creating it if
// necessary.
func GetSentryClient() *raven.Client {
	if sentryClient == nil {
		var err error
		sentryClient, err = raven.NewClient(
			config.SentryDSN,
			map[string]string{"environment": config.RuntimeEnv},
		)
		if err != nil {
			sentryClient = nil
			Panic(err)
		}
	}
	return sentryClient
}

// addSentryHook adds Sentry reporting to the standard logger.
func addSentryHook() {
	hook, err := logrus_sentry.NewSentryHook(config.SentryDSN, []logrus.Level{
		logrus.FatalLevel,
		logrus.ErrorLevel,
	})
	if err == nil {
		logrus.AddHook(hook)
	}
}
