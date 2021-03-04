package main

import (
	"github.com/deng00/go-base/logging"
	"testing"
	"time"
)

func TestSimpleStd(t *testing.T) {
	logConfig := &logging.LogConfig{}
	logConfig.EnableDebug()
	logConfig.EnableHandlerStd()
	logger := logging.GetLogger("test", "test", logConfig)
	logger.Info("show log type ", logging.String("type", "std"))
}

func TestLogFile(t *testing.T) {
	logConfig := &logging.LogConfig{}
	logConfig.EnableDebug()
	logConfig.EnableHandlerFile()
	logger := logging.GetLogger("test", "test", logConfig).Sugar()
	logger.Infof("show log type %s", "file")
}

func TestLogCenter(t *testing.T) {
	logConfig := &logging.LogConfig{}
	logConfig.EnableDebug()
	logConfig.EnableHandlerLogCenter()
	logger := logging.GetLogger("test", "test", logConfig).Sugar()
	logger.Infof("show log type %s", "log center")
}

func TestAlert(t *testing.T) {
	logConfig := &logging.LogConfig{}
	logConfig.EnableHandlerStd()
	logConfig.EnableHandlerFile()
	logConfig.SetAlertChannel(logging.NewDingDingAlertChanel("XXX"))
	logConfig.SetAlertLevel(logging.ErrorLevel)
	logger := logging.GetLogger("test", "test", logConfig).Sugar()
	logger.Errorf("show log type %s", "alert")
	time.Sleep(time.Second)
}
