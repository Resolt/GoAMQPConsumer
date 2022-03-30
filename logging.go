package main

import (
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.StandardLogger()
}

func logInfo(args ...interface{}) {
	logger.Info(args)
}

func logWarn(args ...interface{}) {
	logger.Warn(args)
}

func logError(args ...interface{}) {
	logger.Error(args)
}

func logFatal(args ...interface{}) {
	logger.Fatal(args)
}
