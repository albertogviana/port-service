package log_test

import (
	"os"
	"testing"

	"github.com/albertogviana/port-service/internal/log"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type LoggerUnitTestSuite struct {
	suite.Suite
}

func TestLoggerUnitTestSuite(t *testing.T) {
	suite.Run(t, new(LoggerUnitTestSuite))
}

func (t *LoggerUnitTestSuite) TestNewLogger() {
	t.Run("develop logger", func() {
		l := log.NewLogger(
			"develop",
			os.Stdout,
			true,
		)

		t.Equal("debug", l.GetLevel().String())
		t.Equal(
			&logrus.TextFormatter{
				FullTimestamp:   true,
				TimestampFormat: "2006-01-02 15:04:05",
			},
			l.Formatter.(*logrus.TextFormatter),
		)
	})
	t.Run("production logger", func() {
		l := log.NewLogger(
			"production",
			os.Stdout,
			false,
		)

		t.Equal("info", l.GetLevel().String())
		t.Equal(
			&logrus.JSONFormatter{},
			l.Formatter.(*logrus.JSONFormatter),
		)
	})
	t.Run("production logger with debug enable", func() {
		l := log.NewLogger(
			"production",
			os.Stdout,
			true,
		)

		t.Equal("debug", l.GetLevel().String())
		t.Equal(
			&logrus.JSONFormatter{},
			l.Formatter.(*logrus.JSONFormatter),
		)
	})
}
