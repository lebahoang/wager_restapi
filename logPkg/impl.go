package logPkg

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logInstance *logrus.Logger

func Init() {
	logInstance = logrus.New()
	// Log as JSON instead of the default ASCII formatter.
	logInstance.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logInstance.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logInstance.SetLevel(logrus.InfoLevel)
}

func GetLog() *logrus.Logger {
	return logInstance
}
