package log

import (
	"github.com/sirupsen/logrus"
)

const defaultLogLevel = logrus.WarnLevel

//New logger
func New(logLevel string) Logger {
	// control the log level
	if logLevel != "" {
		level, err := logrus.ParseLevel(logLevel)

		if err != nil {
			logrus.Errorf("the configured log level string is incorrect levelStr=%s error=%s", logLevel, err)
			logrus.SetLevel(defaultLogLevel)
			logrus.Warnf("using %s level instead", defaultLogLevel)
		} else {
			logrus.Info("Setting log level to " + logLevel)
			logrus.SetLevel(level)
		}
	} else {
		// set to default log level
		logrus.SetLevel(defaultLogLevel)
	}
	return logrus.StandardLogger()
}

// author github.com/rmukhamet/
