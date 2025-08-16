package logger

import (
	"context"
	"log"

	"github.com/amahdian/golang-gin-boilerplate/global/env"
	"github.com/amahdian/golang-gin-boilerplate/pkg/logger/logging"
)

// global logger instance
var logger logging.Logger

func ConfigureFromEnvs(envs *env.Envs) logging.Logger {
	var err error
	logger, err = logging.NewLoggerFromEnv(envs)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	return logger
}

func Configure(level logging.LogLevel, logFormat logging.LogFormat) logging.Logger {
	var err error
	logger, err = logging.NewLogger(level, logFormat)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	return logger
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func WithCtx(ctx context.Context) logging.Logger {
	return logger.WithCtx(ctx)
}

func WithFields(fields logging.Fields) logging.Logger {
	return logger.WithFields(fields)
}

func Close() error {
	return logger.Close()
}
