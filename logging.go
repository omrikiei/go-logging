package logging

import (
	"fmt"
	"io"
	"log"
	"os"
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

// LoggerInterface is used for the implementation of loggers with the service
// I really don't see any reason not to used the default logger with it's great formatting, but
// you can implement this interface.
type LoggerInterface interface {
	AddHandler(handler *io.Writer)
	Debug(message string, a ...interface{})
	Info(message string, a ...interface{})
	Warn(message string, a ...interface{})
	Error(message string, a ...interface{})
	Fatal(message string)
	SetFormatter(template string, a ...interface{})
}

// rootLogger implements the ;oggerInterface and is used to log messages
type rootLogger struct {
	handlers  []*LogHandler
	formatter string
}

var instance *rootLogger
var once sync.Once

func get() *rootLogger {
	once.Do(func() {
		instance = &rootLogger{}
	})
	return instance
}

// var loggers = map[string]LoggerInterface{}

// LogHandler provides an output formatter that implements io.Writer, logging level and a formatter
// The best way to instantiate it is by using the NewHandler function:
// 			handler, err := logging.NewHandler(logging.DEBUG, os.Stdout)
//			if err != nil {
//				panic("got an error")
//			}
// Once a handler is created with can pass is to the LoggerInterface implementation with AddHandler
type LogHandler struct {
	Writer    *io.Writer
	Level     int
	Formatter *formatter.LogFormatter
}

// SetFormatter sets a new formatter to the log handler
// receives a pattern string which implements the same logic of "text/template"
// This template will be used when emitting new log records
func (h *LogHandler) SetFormatter(pattern string) *LogHandler {
	h.Formatter = formatter.NewFormatter(pattern)
	return h
}

func (h *LogHandler) emit(message *formatter.LogMessage) {
	logFormatter := *h.Formatter
	logFormatter.Format(h.Writer, *message)
}

// NewHandler will receive a loglevel(either int(0-4) or string) and an io.Writer implementer
// and return an address to a LogHandler instance
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

var defaultHandler, err = NewHandler(DEBUG, os.Stdout)

func formatLevel(levelno int, level, message string, args ...interface{}) *formatter.LogMessage {
	return &formatter.LogMessage{
		fmt.Sprintf(message+"\n", args...),
		level,
		levelno,
	}
}

func dispatchMessage(handlers []*LogHandler, levelNum int, level string, message string, args ...interface{}) {
	if len(handlers) == 0 {
		handlers = []*LogHandler{defaultHandler}
	}
	for _, handler := range handlers {
		if handler.Level <= levelNum {
			handler.emit(formatLevel(levelNum, level, message, args...))
		}
	}
}

func (l *rootLogger) Debug(message string, a ...interface{}) {
	dispatchMessage(l.handlers, DEBUG, "DEBUG", message, a...)
}

func (l *rootLogger) Info(message string, a ...interface{}) {
	dispatchMessage(l.handlers, INFO, "INFO", message, a...)
}

func (l *rootLogger) Warn(message string, a ...interface{}) {
	dispatchMessage(l.handlers, WARN, "WARNING", message, a...)
}

func (l *rootLogger) Error(message string, a ...interface{}) {
	dispatchMessage(l.handlers, ERROR, "ERROR", message, a...)
}

func (l *rootLogger) Fatal(message string, a ...interface{}) {
	dispatchMessage(l.handlers, FATAL, "FATAL", message, a...)
}

// Indirect calls to log/handlers from the package itself will be directed at the root logger

// AddHandler will receive a pointer to a log handler and
// add this handler to the logging module, we will log to multiple loggers
func AddHandler(h *LogHandler) {
	get().AddHandler(h)
}

// Debug will emit a debug message to the handlers that are configured to log debug messages
func Debug(message string, a ...interface{}) {
	get().Debug(message, a...)
}

// Info will emit a debug message to the handlers that are configured to log info messages
func Info(message string, a ...interface{}) {
	get().Info(message, a...)
}

// Warn will emit a debug message to the handlers that are configured to log warning messages
func Warn(message string, a ...interface{}) {
	get().Warn(message, a...)
}

// Error will emit a debug message to the handlers that are configured to log error messages
func Error(message string, a ...interface{}) {
	get().Error(message, a...)
}

// Fatal will emit a debug message to the handlers that are configured to log fatal messages
func Fatal(message string, a ...interface{}) {
	get().Fatal(message, a...)
}
