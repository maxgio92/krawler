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

// Verbosity type
type Verbosity uint32

const (
	PanicLevel Verbosity = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)
