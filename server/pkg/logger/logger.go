package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

const (
	_defaultLevel = zapcore.DebugLevel
)

type Interface interface {
	Debug(string, ...interface{})
	Info(string, ...interface{})
	Warn(string, ...interface{})
	Error(string, ...interface{})
	Fatal(string, ...interface{})
}

type loggerT struct {
	sugaredLogger *zap.SugaredLogger
	atom          zap.AtomicLevel
}

var _ Interface = (*loggerT)(nil)

var logger loggerT

func SetLevel(level string) {
	var l zapcore.Level
	switch strings.ToLower(level) {
	case "debug":
		l = zapcore.DebugLevel
	case "info":
		l = zapcore.InfoLevel
	case "warn":
		l = zapcore.WarnLevel
	case "error":
		l = zapcore.ErrorLevel
	case "fatal":
		l = zapcore.FatalLevel
	}

	logger.atom.SetLevel(l)

}

func init() {
	zConfig := zap.NewDevelopmentEncoderConfig()
	zConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zConfig.EncodeCaller = zapcore.ShortCallerEncoder
	zConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05")

	atom := zap.NewAtomicLevelAt(_defaultLevel)

	z := zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zConfig),
			zapcore.AddSync(os.Stdout),
			atom,
		),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	logger = loggerT{
		sugaredLogger: z.Sugar(),
		atom:          atom,
	}
}

func (l loggerT) Info(format string, args ...interface{}) {
	l.sugaredLogger.Infof(format, args...)
}

func (l loggerT) Warn(format string, args ...interface{}) {
	l.sugaredLogger.Warnf(format, args...)
}

func (l loggerT) Debug(format string, args ...interface{}) {
	l.sugaredLogger.Debugf(format, args...)
}

func (l loggerT) Error(format string, args ...interface{}) {
	l.sugaredLogger.Errorf(format, args...)
}

func (l loggerT) Fatal(format string, args ...interface{}) {
	l.sugaredLogger.Fatalf(format, args...)
}

func Get() Interface {
	return logger
}
