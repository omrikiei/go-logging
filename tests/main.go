package main

import (
	"os"

	"github.com/omrikiei/go-logging"
	"github.com/omrikiei/go-logging/formatter"
)

func main() {
	// Handler with default formatter
	handler, err := logging.NewHandler(logging.DEBUG, os.Stdout)
	if err != nil {
		panic("got an error")
	}
	logging.AddHandler(handler)
	logging.Debug("Testing a debug message")
	logging.Info("Testing an info message")
	logging.Warn("Testing a warning meesage")
	logging.Error("Testing an error message")
	// Setting a different formatter
	handler.SetFormatter("{{ created }}; {{ fileline }}; {{.LevelNum}}; {{.Message}}")
	logging.Debug("Testing a debug message")
	logging.Info("Testing an info message")
	logging.Warn("Testing a warning meesage")
	logging.Error("Testing an error message")
	// Return to the previous formatter
	handler.SetFormatter(formatter.DefaultFormatterPattern)
	str := "hello"
	logging.Debug("Testing a debug message with arguments %s:%d", str, 0)
	logging.Info("Testing an info message with arguments %s:%d", str, 1)
	logging.Warn("Testing a warning meesage with arguments %s:%d", str, 2)
	logging.Error("Testing an error message with arguments %s:%d", str, 3)
}
