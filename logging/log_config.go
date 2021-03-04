package logging

import (
	"github.com/deng00/go-base/config"
	"go.uber.org/zap/zapcore"
)

type LogConfig struct {
	isDebugMode            bool
	enableHandlerStd       bool
	enableHandlerFile      bool
	enableHandlerLogCenter bool
	alertLevel             zapcore.Level
	logLevel               zapcore.Level
	logFileDir             string
	alertChannel           AlertChannelI
}

func (l *LogConfig) EnableDebug() {
	l.isDebugMode = true
}

func (l *LogConfig) EnableHandlerStd() {
	l.enableHandlerStd = true
}

func (l *LogConfig) EnableHandlerFile() {
	l.enableHandlerFile = true
}

func (l *LogConfig) EnableHandlerLogCenter() {
	l.enableHandlerLogCenter = true
}

func (l *LogConfig) SetAlertLevel(level zapcore.Level) {
	l.alertLevel = level
}

func (l *LogConfig) SetAlertChannel(channel AlertChannelI) {
	l.alertChannel = channel
}

func (l *LogConfig) SetLogLevel(level zapcore.Level) {
	l.logLevel = level
}

func (l *LogConfig) SetLogFileDir(dir string) {
	l.logFileDir = dir
}

func GetLogConfig(config *config.Config) *LogConfig {
	logConfig := &LogConfig{}
	if config.GetBool("log.handlerStd") {
		logConfig.EnableHandlerStd()
	}
	if config.GetBool("log.handlerFile") {
		logConfig.EnableHandlerFile()
	}
	if config.GetString("log.fileDir") != "" {
		logConfig.SetLogFileDir(config.GetString("log.fileDir"))
	}
	// 告警配置
	alert := config.GetStringMapString("log.alert")
	if alert["channel"] != "" && alert["channel_param"] != "" {
		alertLevel := ErrorLevel
		err := alertLevel.Set(alert["level"])
		if err != nil {
			println("alert level config error, should be one of warn/error/dpanic/panic/fatal")
		}
		logConfig.SetAlertLevel(alertLevel)
		var alertChanel AlertChannelI
		switch alert["channel"] {
		case "dingding":
			alertChanel = NewDingDingAlertChanel(alert["channel_param"])
		default:
			println("ignore unknown alert channel, " + alert["channel"])
		}
		logConfig.SetAlertChannel(alertChanel)
	}
	// 日志等级
	level := DebugLevel
	err := level.Set(config.GetString("log.level"))
	if err != nil {
		println("log level config error, should be one of debug/info/warn/error/dpanic/panic/fatal")
	}
	if level == DebugLevel {
		logConfig.EnableDebug()
	}
	logConfig.SetLogLevel(level)
	return logConfig
}
