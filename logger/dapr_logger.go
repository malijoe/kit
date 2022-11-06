package logger

import (
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// daprLogger is the implementation for logrus.
type daprLogger struct {
	// name is the name of the logger that is published to log as a scope
	name string
	// logger is the instance of logrus logger
	logger *logrus.Entry
}

var DaprVersion = "unknown"

func newDaprLogger(name string) *daprLogger {
	newLogger := logrus.New()
	newLogger.SetOutput(os.Stdout)

	dl := &daprLogger{
		name: name,
		logger: newLogger.WithFields(logrus.Fields{
			logFieldScope: name,
			logFieldType:  LogTypeLog,
		}),
	}

	dl.EnableJSONOutput(defaultJSONOutput)
	return dl
}

func (l *daprLogger) EnableJSONOutput(enabled bool) {
	var formatter logrus.Formatter

	fieldMap := logrus.FieldMap{
		// If time field name is conflicted, logrus adds "fields." prefix.
		// So rename to unused field @time to avoid the conflict.
		logrus.FieldKeyTime:  logFieldTimeStamp,
		logrus.FieldKeyLevel: logFieldLevel,
		logrus.FieldKeyMsg:   logFieldMessage,
	}

	hostname, _ := os.Hostname()
	l.logger.Data = logrus.Fields{
		logFieldScope:    l.logger.Data[logFieldScope],
		logFieldType:     LogTypeLog,
		logFieldInstance: hostname,
		logFieldDaprVer:  DaprVersion,
	}

	if enabled {
		formatter = &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
			FieldMap:        fieldMap,
		}
	} else {
		formatter = &logrus.TextFormatter{
			TimestampFormat: time.RFC3339Nano,
			FieldMap:        fieldMap,
		}
	}

	l.logger.Logger.SetFormatter(formatter)
}

// SetAppID sets app_id field in the log. Default value is empty string.
func (l *daprLogger) SetAppID(id string) {
	l.logger = l.logger.WithField(logFieldAppID, id)
}

func toLogrusLevel(lvl LogLevel) logrus.Level {
	// ignore error becuase it will never happen
	l, _ := logrus.ParseLevel(string(lvl))
	return l
}

// SetOutputLevel sets log output level.
func (l *daprLogger) SetOutputLevel(lvl LogLevel) {
	l.logger.Logger.SetLevel(toLogrusLevel(lvl))
}

// IsOutputLevelEnabled returns true if the logger will output this LogLevel.
func (l *daprLogger) IsOutputLevelEnabled(level LogLevel) bool {
	return l.logger.Logger.IsLevelEnabled(toLogrusLevel(level))
}

// SetOutput sets the destination for the logs.
func (l *daprLogger) SetOutput(dst io.Writer) {
	l.logger.Logger.SetOutput(dst)
}

// WithLogType specify the log_type field in log. Default value is LogTypeLog.
func (l *daprLogger) WithLogType(logType string) Logger {
	return &daprLogger{
		name:   l.name,
		logger: l.logger.WithField(logFieldType, logType),
	}
}

func (l *daprLogger) WithFields(fields map[string]any) Logger {
	return &daprLogger{
		name:   l.name,
		logger: l.logger.WithFields(fields),
	}
}

func (l *daprLogger) Info(args ...any) {
	l.logger.Log(logrus.InfoLevel, args...)
}

func (l *daprLogger) Infof(format string, args ...any) {
	l.logger.Logf(logrus.InfoLevel, format, args...)
}

func (l *daprLogger) Debug(args ...any) {
	l.logger.Log(logrus.DebugLevel, args...)
}

func (l *daprLogger) Debugf(format string, args ...any) {
	l.logger.Logf(logrus.DebugLevel, format, args...)
}

func (l *daprLogger) Warn(args ...any) {
	l.logger.Log(logrus.WarnLevel, args...)
}

func (l *daprLogger) Warnf(format string, args ...any) {
	l.logger.Logf(logrus.WarnLevel, format, args...)
}

func (l *daprLogger) Error(args ...any) {
	l.logger.Log(logrus.ErrorLevel, args...)
}

func (l *daprLogger) Errorf(format string, args ...any) {
	l.logger.Logf(logrus.ErrorLevel, format, args...)
}

func (l *daprLogger) Fatal(args ...any) {
	l.logger.Fatal(args...)
}

func (l *daprLogger) Fatalf(format string, args ...any) {
	l.logger.Fatalf(format, args...)
}
