package logging

import (
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

//NewDefaultLogger returns an default instance of logger
func NewDefaultLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.SetReportCaller(true)
	return logger
}

func NewZapLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction(zap.AddCaller(), zap.Development())
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	return sugar
}
