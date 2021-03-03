package log

import (
	"github.com/go-kit/kit/log"
	"os"
)

func InitLogger(jobID string) (log.Logger, log.Logger, log.Logger, log.Logger) {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller, "jobID", jobID)

	cloneLogger := log.With(logger, "component", "clone")
	buildLogger := log.With(logger, "component", "build")
	pushLogger := log.With(logger, "components", "push")

	return logger, cloneLogger, buildLogger, pushLogger
}
