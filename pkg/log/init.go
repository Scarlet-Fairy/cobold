package log

import (
	"github.com/go-kit/kit/log"
	"os"
)

func InitLogger() log.Logger {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	return logger
}
