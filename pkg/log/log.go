package log

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"os"
)

func InitLogger(jobID string) (log.Logger, log.Logger, log.Logger, log.Logger, log.Logger) {
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller, "jobID", jobID)
	logger = level.NewInjector(logger, level.InfoValue())

	cloneLogger := log.With(logger, "component", "clone")
	buildLogger := log.With(logger, "component", "build")
	pushLogger := log.With(logger, "component", "push")
	notifyLogger := log.With(logger, "component", "notify")

	return logger, cloneLogger, buildLogger, pushLogger, notifyLogger
}
