// Package logger provides log capabilities
package logger

import (
	"auri/config"
	"io"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/gobuffalo/buffalo"
)

// BuffaloLogger provides main buffalo logger
var BuffaloLogger buffalo.Logger

// IPALogger provides the logger for IPA operations log
var IPALogger *logrus.Logger

// AuriLogger provides the logger for all auri operations
var AuriLogger *logrus.Logger

// SetLoggers initializes loggers
func SetLoggers(buffaloLogger buffalo.Logger) error {
	BuffaloLogger = buffaloLogger

	return CreateFileLoggers()
}

// CreateFileLoggers sets up our own loggers to the files
func CreateFileLoggers() error {
	if config.IsProdBuild() { //using real files for production builds
		ipaFile, err := os.OpenFile(config.GetInstance().LogIPA, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
		if err != nil {
			return err
		}

		IPALogger = newFileLogger(ipaFile)

		auriFile, err := os.OpenFile(config.GetInstance().LogAuri, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
		if err != nil {
			return err
		}

		AuriLogger = newFileLogger(auriFile)
	} else { //development mode, using stdout
		IPALogger = newFileLogger(os.Stdout)
		AuriLogger = newFileLogger(os.Stdout)
	}

	return nil
}

// newFileLogger creates new logrus logger to a given output
func newFileLogger(output io.Writer) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	logger.SetOutput(output)

	return logger
}
