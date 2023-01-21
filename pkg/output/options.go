package output

type Options struct {
	Logger    *Logger
	Verbosity Verbosity
}

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
