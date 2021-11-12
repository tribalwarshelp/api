package internal

import (
	"github.com/Kichiyaki/appmode"
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"os"
)

const (
	sentryAppName = "twhelp-api"
)

func InitSentry(version string) error {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		Environment:      appmode.Get(),
		Release:          sentryAppName + "@" + version,
		Debug:            false,
		TracesSampleRate: 0.3,
	})
	if err != nil {
		return errors.Wrap(err, "sentry.Init")
	}

	return nil
}
