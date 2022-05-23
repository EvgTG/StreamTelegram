package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

const defaultLogLevel = logrus.WarnLevel

func New(logLevel string, toFile bool) Logger {
	logger := logrus.New()

	if logLevel != "" {
		level, err := logrus.ParseLevel(logLevel)
		if err != nil {
			logrus.Errorf("the configured log level string is incorrect levelStr=%s error=%s", logLevel, err)
			logger.SetLevel(defaultLogLevel)
			logrus.Warnf("using %s level instead", defaultLogLevel)
		} else {
			logrus.Info("Setting log level to " + logLevel)
			logger.SetLevel(level)
		}
	} else {
		logger.SetLevel(defaultLogLevel)
	}

	if toFile {
		file, err := os.OpenFile("files/logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logger.SetOutput(file)
		} else {
			panic(err)
		}
	}

	return logger
}
