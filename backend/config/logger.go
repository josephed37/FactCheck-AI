package config

import (
	"os"
	"github.com/sirupsen/logrus"
)

// Log is a global logger instance that can be used throughout the application.
var Log *logrus.Logger

// InitLogger initializes the global logger with a JSON format.
func InitLogger() {
	Log = logrus.New()

	// Set the formatter to output logs in JSON format.
	Log.SetFormatter(&logrus.JSONFormatter{})

	// Set the output to standard out (the console).
	// In a production environment, this could be a file or a log management service.
	Log.SetOutput(os.Stdout)

	// Set the default logging level.
	// We can filter out less important logs by setting this to WarnLevel or ErrorLevel.
	Log.SetLevel(logrus.InfoLevel)
}