package output

import (
	log "github.com/sirupsen/logrus"
)

type Logger struct {
	log.Logger
}

func NewLogger() *Logger {
	return &Logger{*log.New()}
}
