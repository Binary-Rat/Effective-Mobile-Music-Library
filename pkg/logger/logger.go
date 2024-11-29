package logger

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func ConfigureLogger() *log.Logger {
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	if strings.ToLower(os.Getenv("DEBUG")) == "true" {
		logger.SetLevel(log.DebugLevel)
	} else {
		logger.SetLevel(log.InfoLevel)
	}
	return logger
}
