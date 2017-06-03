package instahelper

import (
	"os"

	"strings"

	"github.com/sirupsen/logrus"
)

// Log is the main logger for instahelper
var Log = logrus.New()

// SetLoggingLevel Sets the logging level of the app Debug, Error, Fatal, Info, Panic, Warn
func SetLoggingLevel(level string) {
	level = strings.ToLower(level)
	switch level {
	case "debug":
		Log.Level = logrus.DebugLevel

	case "error":
		Log.Level = logrus.ErrorLevel

	case "fatal":
		Log.Level = logrus.FatalLevel

	case "info":
		Log.Level = logrus.InfoLevel

	case "panic":
		Log.Level = logrus.PanicLevel

	case "warn":
		Log.Level = logrus.WarnLevel

	default:
		Log.Info("Invalid logging level given. Logging level set to default of Warn.")
	}
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Set log severity
	Log.Level = logrus.InfoLevel
}
