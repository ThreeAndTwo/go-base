package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

type Level = zapcore.Level

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = zapcore.DebugLevel
	// InfoLevel is the default logging priority.
	InfoLevel = zapcore.InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = zapcore.WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = zapcore.ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel = zapcore.DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel = zapcore.PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = zapcore.FatalLevel
)

type Field = zap.Field
type Logger = zap.Logger
type SugaredLogger = zap.SugaredLogger

// String constructs a field with the given key and value.
func String(key string, val string) Field {
	return zap.String(key, val)
}

func Bool(key string, val bool) Field {
	return zap.Bool(key, val)
}

func Float64(key string, val float64) Field {
	return zap.Float64(key, val)
}

func Int64(key string, val int64) Field {
	return zap.Int64(key, val)
}

func Uint64(key string, val uint64) Field {
	return zap.Uint64(key, val)
}

func Any(key string, value interface{}) Field {
	return zap.Any(key, value)
}

func GetLogger(serviceName, moduleName string, config *LogConfig) *Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)
	consoleEncodeConfig := zap.NewDevelopmentEncoderConfig()
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(config.logLevel)
	var options []zap.Option
	var cores []zapcore.Core
	if config.enableHandlerFile {
		logDir := config.logFileDir
		if logDir == "" {
			logDir = "./logs"
		}
		hook := lumberjack.Logger{
			Filename:   logDir + "/" + serviceName + "-" + moduleName + ".log", // 日志文件路径
			MaxSize:    1024,                                                   // 每个日志文件保存的最大尺寸 单位：M
			MaxBackups: 30,                                                     // 日志文件最多保存多少个备份
			MaxAge:     7,                                                      // 文件最多保存多少天
			Compress:   true,                                                   // 是否压缩
		}
		cores = append(cores, zapcore.NewCore(jsonEncoder, zapcore.AddSync(&hook), atomicLevel))
	}
	if config.enableHandlerStd {
		consoleEncodeConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 颜色等级
		consoleEncoder := zapcore.NewConsoleEncoder(consoleEncodeConfig)
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), atomicLevel))
	}
	if config.enableHandlerLogCenter {
		logCenterHandler := LogCenterSyncer{}
		cores = append(cores, zapcore.NewCore(jsonEncoder, zapcore.AddSync(&logCenterHandler), atomicLevel))
	}
	if config.alertChannel != nil {
		if config.alertLevel <= InfoLevel {
			config.alertLevel = ErrorLevel
		}
		highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= config.alertLevel
		})
		alertHandler := AlertSyncer{}
		alertHandler.setChannel(serviceName, config.alertChannel)
		cores = append(cores, zapcore.NewCore(jsonEncoder, zapcore.AddSync(&alertHandler), highPriority))
	}
	if len(cores) == 0 {
		println("go-base-logger: you should set at least one log handler")
	}
	var core = zapcore.NewTee(cores...)
	// 构造日志
	var logger *Logger
	if config.isDebugMode {
		// 开启开发模式，堆栈跟踪
		options = append(options, zap.AddCaller())
		// 开启文件及行号
		options = append(options, zap.Development())
	}
	filed := zap.Fields(zap.String("service", serviceName), zap.String("module", moduleName))
	options = append(options, filed)
	logger = zap.New(core, options...)
	return logger
}
