package main

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/getsentry/raven-go"
)

var allLevels = []logrus.Level{
	logrus.DebugLevel,
	logrus.InfoLevel,
	logrus.WarnLevel,
	logrus.FatalLevel,
	logrus.ErrorLevel,
	logrus.PanicLevel,
}

// sentryHook implements the logrus.Hook interface so that errors and events
// are logged to sentry if the SENTRY_DSN environment variable is set
type sentryHook struct {
	client *raven.Client
}

func newSentryHook(dsn string, tags map[string]string) (logrus.Hook, error) {
	client, err := raven.NewClient(dsn, tags)
	if err != nil {
		return nil, err
	}

	return &sentryHook{
		client: client,
	}, nil
}

func (s *sentryHook) Levels() []logrus.Level {
	return allLevels
}

func (s *sentryHook) Fire(e *logrus.Entry) error {
	packet := raven.NewPacket(e.Message)

	packet.Level = s.ravenLevel(e.Level)

	switch packet.Level {
	case raven.ERROR, raven.WARNING, raven.FATAL:
		packet.Interfaces = append(packet.Interfaces, raven.NewStacktrace(4, 0, nil))
	}

	_, cerr := s.client.Capture(packet, s.ravenTags(e.Data))

	return <-cerr
}

func (s *sentryHook) ravenLevel(l logrus.Level) raven.Severity {
	switch l {
	case logrus.DebugLevel:
		return raven.DEBUG
	case logrus.InfoLevel:
		return raven.INFO
	case logrus.WarnLevel:
		return raven.WARNING
	case logrus.ErrorLevel:
		return raven.ERROR
	case logrus.FatalLevel, logrus.PanicLevel:
		return raven.FATAL
	}

	return raven.ERROR
}

func (s *sentryHook) ravenTags(data logrus.Fields) map[string]string {
	out := make(map[string]string, len(data))

	for k, v := range data {
		out[k] = fmt.Sprint(v)
	}

	return out
}
