package test

import (
	"os"
	"time"

	"github.com/amahdian/golang-gin-boilerplate/pkg/logger"
	"github.com/amahdian/golang-gin-boilerplate/pkg/logger/logging"
)

func SetupTestingEnv() {
	time.Local = time.UTC
	logger.Configure(logging.DebugLevel, logging.TextFormat)
	_ = os.Setenv("PROFILE", "testing")
}
