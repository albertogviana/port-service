package log

import (
	"io"

	"github.com/sirupsen/logrus"
)

// NewLogger returns an instance of logger ready to use.
func NewLogger(environment string, output io.Writer, debugLevel bool) *logrus.Logger {
	l := logrus.New()

	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	if environment == "production" {
		l.SetFormatter(&logrus.JSONFormatter{})
	}

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	l.SetOutput(output)

	if debugLevel || environment == "develop" {
		// Only log the warning severity or above.
		l.SetLevel(logrus.DebugLevel)
	}

	return l
}
