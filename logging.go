package logging

import (
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
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
type Interface interface {
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

var loggers = map[string]Interface{}

// LogHandler provides an output formatter that implements io.Writer and a logging level
type LogHandler struct {
	Writer *log.Logger
	Level  int
}

func (h *LogHandler) emit(message string, a ...interface{}) {
	h.Writer.Printf(message, a...)
}

// NewHandler instance
func NewHandler(level interface{}, w io.Writer) (*LogHandler, error) {
	switch level.(type) {
	case int:
		if level.(int) < DEBUG || level.(int) > FATAL {
			log.Fatal("Level is not supported, use INFO/DEBUG/WARN/ERROR/FATAL.")
		}
		return &LogHandler{Writer: log.New(w, "", log.Ldate|log.Ltime|log.Lmicroseconds), Level: level.(int)}, nil
	case string:
		if _, ok := logLevels[level.(string)]; !ok {
			log.Fatal("Level is not supported, use INFO/DEBUG/WARN/ERROR/FATAL.")
		}
		return &LogHandler{Writer: log.New(w, "", log.Ldate|log.Ltime|log.Lmicroseconds), Level: logLevels[level.(string)]}, nil
	default:
		log.Fatal("Log level should either be of type int or string")
	}
	return &LogHandler{}, errors.New("Bad use of new LogHandler")
}

func (l *rootLogger) AddHandler(h *LogHandler) {
	l.handlers = append(l.handlers, h)
}

func formatLevel(level, message string) string {
	return fmt.Sprintf("%s: %s", level, message)
}

func (l *rootLogger) Debug(message string, a ...interface{}) {
	for _, handler := range l.handlers {
		if handler.Level <= DEBUG {
			handler.emit(formatLevel("DEBUG", message), a...)
		}
	}
}

func (l *rootLogger) Info(message string, a ...interface{}) {
	for _, handler := range l.handlers {
		if handler.Level <= INFO {
			fmt.Printf("%d", len(a))
			handler.emit(formatLevel("INFO", message), a...)
		}
	}
}

func (l *rootLogger) Warn(message string, a ...interface{}) {
	for _, handler := range l.handlers {
		if handler.Level <= WARN {
			handler.emit(formatLevel("WARNING", message), a...)
		}
	}
}

func (l *rootLogger) Error(message string, a ...interface{}) {
	for _, handler := range l.handlers {
		if handler.Level <= ERROR {
			handler.emit(formatLevel("ERROR", message), a...)
		}
	}
}

func (l *rootLogger) Fatal(message string, a ...interface{}) {
	for _, handler := range l.handlers {
		if handler.Level <= FATAL {
			handler.emit(formatLevel("FATAL", message), a...)
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
