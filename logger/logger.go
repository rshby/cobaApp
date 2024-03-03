package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

func NewConsoleLog() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.DebugLevel)
	log.SetOutput(os.Stdout)
	return log
}
