// internal/logger/logger.go
package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

func Init(logLevel string) {
	// Устанавливаем формат JSON для логов (удобно для структурированных логов)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	// Уровень логирования (info, debug, error и т.д.)
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)
}
