package logger

import (
	"io"
	"strings"
	"sync"
)

const (
	// LogTypeLog is the normal log type.
	LogTypeLog = "log"
	// LogTypeRequest is the request log type.
	LogTypeRequest = "request"

	// Field names that defines Daprdan log schema.
	logFieldTimeStamp = "time"
	logFieldLevel     = "level"
	logFieldType      = "type"
	logFieldScope     = "scope"
	logFieldMessage   = "msg"
	logFieldInstance  = "instance"
	logFieldDaprVer   = "ver"
	logFieldAppID     = "app_id"
)

type LogLevel string

const (
	// DebugLevel for verbose messages.
	DebugLevel LogLevel = "debug"
	// InfoLevel the default log level.
	InfoLevel LogLevel = "info"
	// WarnLevel for logging messages about possible issues.
	WarnLevel LogLevel = "warn"
	// ErrorLevel for logging errors.
	ErrorLevel LogLevel = "error"
	// FatalLevel for logging fatal messages. The system shuts down after logging the message.
	FatalLevel LogLevel = "fatal"

	// UndefinedLevel is for unrecognized log levels.
	UndefinedLevel LogLevel = "undefined"
)

// globalLoggers is the collection of Daprdan Loggers that are shared globally.
// TODO: User will disable/enable logger on demand.
var (
	globalLoggers     = map[string]Logger{}
	globalLoggersLock = sync.RWMutex{}
)

type Logger interface {
	// EnableJSONOutput enables JSON formatted output log
	EnableJSONOutput(enabled bool)

	// SetAppID sets daprd_id field in the logs. Default value is empty string
	SetAppID(id string)

	// SetOutputLevel sets the log output level.
	SetOutputLevel(outputLevel LogLevel)
	// SetOutput sets the destination for the logs.
	SetOutput(dst io.Writer)

	// IsOutputLevelEnabled returns true if the logger will output this LogLevel.
	IsOutputLevelEnabled(level LogLevel) bool

	// WithLogType specifies the log_type field in log. Default values is LogTypeLog.
	WithLogType(logType string) Logger

	// WithFields returns a logger with the added structured fields.
	WithFields(fields map[string]any) Logger

	// Info logs a message at the Info level.
	Info(args ...any)
	// Infof logs a message at the Info level.
	Infof(format string, args ...any)
	// Debug logs a message at the Debug level.
	Debug(args ...any)
	// Debugf logas a message at the Debug level.
	Debugf(format string, args ...any)
	// Warn logs a message at the Warn level.
	Warn(args ...any)
	// Warnf logs a message at the Warn level.
	Warnf(format string, args ...any)
	// Error logs a message at the Error level.
	Error(args ...any)
	// Errorf logs a message at the Error level.
	Errorf(format string, args ...any)
	// Fatal logs a message at the Fatal level then the process will exit with status set to 1.
	Fatal(args ...any)
	// Fatalf logs a message at the Fatal level then the process will exit with status set to 1.
	Fatalf(format string, args ...any)
}

// toLogLevel converts to LogLevel.
func toLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	}

	// unsupported log level
	return UndefinedLevel
}

func NewLogger(name string) Logger {
	globalLoggersLock.Lock()
	defer globalLoggersLock.Unlock()

	logger, ok := globalLoggers[name]
	if !ok {
		logger = newDaprLogger(name)
		globalLoggers[name] = logger
	}

	return logger
}

func getLoggers() map[string]Logger {
	globalLoggersLock.RLock()
	defer globalLoggersLock.RUnlock()

	l := map[string]Logger{}
	for k, v := range globalLoggers {
		l[k] = v
	}

	return l
}
