package logging

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	// ISO8601 is timestamp format used by Logger.
	ISO8601        = "2006-01-02T15:04:05.000Z07:00"
	loggerFieldKey = "logger"
)

type Logger struct {
	entry *logrus.Entry
	depth int
}

// Infof logs a message at level Info on the standard logger.
func (l Logger) Infof(format string, args ...interface{}) {
	l.sourced(l.depth).Infof(format, args...)
}

// sourced adds a source field to the logger that contains
// the file name and line where the logging happened.
func (l Logger) sourced(depth int) *logrus.Entry {
	_, file, line, ok := runtime.Caller(depth + 2)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	return l.entry.WithField(loggerFieldKey, fmt.Sprintf("%s:%d", file, line))
}

func NewLogger() *Logger {
	return &Logger{}
}
