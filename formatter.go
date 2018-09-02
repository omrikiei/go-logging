package logging

import (
	"fmt"
	"io"
	"log"
	"runtime"
	"text/template"
	"time"
)

// this magic number represents the stack location of the invoking log command (.Debug, .Error etc...)
// there is only a need to touch it if the stack call changes in the number of calls
var stackLocation = 16

func asctime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func created() int64 {
	return time.Now().Unix()
}

func filename() string {
	_, filename, _, ok := runtime.Caller(stackLocation)
	if !ok {
		return "<unknown>"
	}
	return filename
}

func lineno() int {
	_, _, lineno, ok := runtime.Caller(stackLocation)
	if !ok {
		return -1
	}
	return lineno
}

func filenameAndLineno() string {
	_, filename, lineno, ok := runtime.Caller(stackLocation)
	if !ok {
		return "<unknown>: -1"
	}
	return fmt.Sprintf("%s: %d", filename, lineno)
}

var functions = template.FuncMap{
	"asctime":  asctime,
	"created":  created,
	"filename": filename,
	"lineno":   lineno,
	"fileline": filenameAndLineno,
}

// LogMessage is a basic represantation of the message we write to the log stream
// Message is the message string
// Level is the debug level(string)
// LevelNum is the debug level int representation
type LogMessage struct {
	Message  string
	Level    string
	LevelNum int
}

// LogFormatter is a formatter that parses log messages
// it implements a template based on the "text/template" package
type LogFormatter struct {
	Template *template.Template
}

// New is the constructor for a LogFormatter, it
// accepts a pattern and binds the log formatting functions and pattern to the LogFormatters template
func NewFormatter(pattern string) (*LogFormatter, error) {
	template, err := template.New("logTemplate").Funcs(functions).Parse(pattern)
	if err != nil {
		return nil, err
	}
	return &LogFormatter{
		template,
	}, nil
}

// Format is the function that the formatter uses to format a string and write it to an io.Writer
func (logFormatter *LogFormatter) Format(writer *io.Writer, message *LogMessage) {
	err := logFormatter.Template.Execute(*writer, message)
	if err != nil {
		log.Printf("logging package failed to emit a message to %s", writer)
	}
}

// DefaultFormatterPattern is the pattern used when no formatter is set
const DefaultFormatterPattern = "{{ asctime }}; {{ fileline }}; {{.Level}}; {{.Message}}"

// DefaultFormatter is the default formatter that will be used by the logging module, it implements the DefaultFormatterPattern
var DefaultFormatter, _ = NewFormatter(DefaultFormatterPattern)
