package utils

import (
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var log = logrus.New()

// InitLogger initializes the logger, sets the formatter and the output file
func InitLogger() error {
	log.Formatter = &logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	}
	file, err := os.OpenFile(viper.GetString("logging.path"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return errors.Wrap(err, "failed to open log file")
	}
	log.Out = file
	return nil
}

// Info logs an info message
func Info(args ...interface{}) {
	log.Info(args...)
}

// Warn logs a warning message
func Warn(args ...interface{}) {
	log.Warn(args...)
}

// Error logs an error message
func Error(args ...interface{}) {
	log.Error(args...)
}
