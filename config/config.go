// Package config provides access to application-wide configuration.
// All configuration should be via environment variables. See http://12factor.net/config.
package config

import "os"

var (
	// SlackAPIURL is the Slack API URL
	SlackAPIURL = getEnvDefault("SLACK_API_URL", "https://slack.com/api")

	// RuntimeEnv e.g. "production", "staging", "development"
	RuntimeEnv = getEnvDefault("RUNTIME_ENV", "development")

	// SentryDSN is the Sentry data source name, the url where the client should
	// report errors.
	SentryDSN = os.Getenv("SENTRY_DSN")

	// AWSRegion is the primary aws region
	AWSRegion = getEnvDefault("AWS_REGION", "us-east-1")
)

func getEnvDefault(key string, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}
