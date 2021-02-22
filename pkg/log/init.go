package log

import "github.com/sirupsen/logrus"

var Logger = logrus.New()

func InitLogger(level string) {
	Logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	parsedLevel, err := logrus.ParseLevel(level)
	if err != nil {
		Logger.Panicf("Level: %s is not a valid log level", level)
	}

	Logger.SetLevel(parsedLevel)
}
