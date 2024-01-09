package log

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var StandardLogger = logrus.StandardLogger()

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func GetLevel() logrus.Level {
	return logrus.GetLevel()
}

func Open(path string, level string) error {
	if level != "" {
		lvl, err := logrus.ParseLevel(level)
		if err != nil {
			return err
		}
		logrus.SetLevel(lvl)
	}

	switch strings.ToLower(path) {
	case "":
		break
	case "stdout":
		logrus.SetOutput(os.Stdout)
	case "stderr":
		logrus.SetOutput(os.Stderr)
	default:
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		logrus.SetOutput(file)
	}
	return nil
}

func Close() error {
	file, isFile := StandardLogger.Out.(*os.File)
	if !isFile {
		logrus.Debug("not logging to a file, ignoring close request")
		return nil
	}
	logrus.SetOutput(os.Stdout)
	return file.Close()
}

func Rotate() error {
	file, isFile := StandardLogger.Out.(*os.File)
	if !isFile {
		logrus.Debug("not logging to a file, ignoring rotate request")
		return nil
	}

	path := file.Name()

	logrus.SetOutput(os.Stdout)
	if err := file.Close(); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	logrus.SetOutput(file)

	return nil
}
