package logging

import (
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/omrikiei/go-logging/formatter"
)

const (
	DEBUG = 0
	INFO  = 1
	WARN  = 2
	ERROR = 3
	FATAL = 4
)

var logLevels = map[string]int{
	"DEBUG": 0,
	"INFO":  1,
	"WARN":  2,
	"ERROR": 3,
	"FATAL": 4,
}

var initialized uint32
var mu sync.Mutex

// Interface is used for the implementation of loggers with the service
type LoggerInterface interface {
	AddHandler(handler *io.Writer)
	Debug(message string, a ...interface{})
	Info(message string, a ...interface{})
	Warn(message string, a ...interface{})
	Error(message string, a ...interface{})
	Fatal(message string)
	SetFormatter(template string, a ...interface{})
}

// rootLogger implements the LoggerInterface and is used to log messages
type rootLogger struct {
	handlers  []*LogHandler
	formatter string
}

var instance *rootLogger
var once sync.Once

// Get provides an instance getter for the logging object
func Get() *rootLogger {
	once.Do(func() {
		instance = &rootLogger{}
	})
	return instance
}

// var loggers = map[string]LoggerInterface{}

// LogHandler provides an output formatter that implements io.Writer and a logging level
type LogHandler struct {
	Writer    *io.Writer
	Level     int
	Formatter *formatter.LogFormatter
}

// Sets a new formatter to the log handler
func (h *LogHandler) SetFormatter(pattern string) *LogHandler {
	h.Formatter = formatter.NewFormatter(pattern)
	return h
}

func (h *LogHandler) emit(message *formatter.LogMessage) {
	logFormatter := *h.Formatter
	logFormatter.Format(h.Writer, *message)
}

// NewHandler instance
func NewHandler(level interface{}, w io.Writer) (*LogHandler, error) {
	switch level.(type) {
	case int:
		if level.(int) < DEBUG || level.(int) > FATAL {
			log.Fatal("Level is not supported, use INFO/DEBUG/WARN/ERROR/FATAL.")
		}
		return &LogHandler{Writer: &w, Level: level.(int), Formatter: formatter.DefaultFormatter}, nil
	case string:
		if _, ok := logLevels[level.(string)]; !ok {
			log.Fatal("Level is not supported, use INFO/DEBUG/WARN/ERROR/FATAL.")
		}
		return &LogHandler{Writer: &w, Level: logLevels[level.(string)], Formatter: formatter.DefaultFormatter}, nil
	default:
		return &LogHandler{Writer: &w, Level: DEBUG, Formatter: formatter.DefaultFormatter}, nil
	}
}

func (l *rootLogger) AddHandler(h *LogHandler) {
	l.handlers = append(l.handlers, h)
}

func formatLevel(levelno int, level, message string, args ...interface{}) *formatter.LogMessage {
	return &formatter.LogMessage{
		fmt.Sprintf(message+"\n", args...),
		level,
		levelno,
	}
}

func (l *rootLogger) Debug(message string, a ...interface{}) {
	for _, handler := range l.handlers {
		if handler.Level <= DEBUG {
			handler.emit(formatLevel(DEBUG, "DEBUG", message, a...))
		}
	}
}

func (l *rootLogger) Info(message string, a ...interface{}) {
	for _, handler := range l.handlers {
		if handler.Level <= INFO {
			handler.emit(formatLevel(INFO, "INFO", message, a...))
		}
	}
}

func (l *rootLogger) Warn(message string, a ...interface{}) {
	for _, handler := range l.handlers {
		if handler.Level <= WARN {
			handler.emit(formatLevel(WARN, "WARNING", message, a...))
		}
	}
}

func (l *rootLogger) Error(message string, a ...interface{}) {
	for _, handler := range l.handlers {
		if handler.Level <= ERROR {
			handler.emit(formatLevel(ERROR, "ERROR", message, a...))
		}
	}
}

func (l *rootLogger) Fatal(message string, a ...interface{}) {
	for _, handler := range l.handlers {
		if handler.Level <= FATAL {
			handler.emit(formatLevel(FATAL, "FATAL", message, a...))
		}
	}
}

// Indirect calls to log/handlers from the package itself will be directed at the root logger

// AddHandler called from the looging package will invoke in the rootLogger
func AddHandler(h *LogHandler) {
	Get().AddHandler(h)
}

// Debug is called from the package and implements in the rootLogger instance
func Debug(message string, a ...interface{}) {
	Get().Debug(message, a...)
}

// Info is called from the package and implements in the rootLogger instance
func Info(message string, a ...interface{}) {
	Get().Info(message, a...)
}

// Warn is called from the package and implements in the rootLogger instance
func Warn(message string, a ...interface{}) {
	Get().Warn(message, a...)
}

// Error is called from the package and implements in the rootLogger instance
func Error(message string, a ...interface{}) {
	Get().Error(message, a...)
}

// Fatal is called from the package and implements in the rootLogger instance
func Fatal(message string, a ...interface{}) {
	Get().Fatal(message, a...)
}
