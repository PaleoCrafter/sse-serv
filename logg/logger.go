// MIT license, dtg [at] lengo [dot] org · 10/2018

package logg

import (
	"fmt"
	"os"
	"time"
)

// Logger ...
type Logger interface {
	Open() error
	Info(format string, msg ...interface{})
	Error(format string, msg ...interface{})
	Warn(format string, msg ...interface{})
}

type logger struct {
	filename string
	file     *os.File
}

// NewLogger ...
func NewLogger(filename string) Logger {
	return &logger{filename: filename}
}

func (l *logger) Open() error {
	var file *os.File
	var err error

	flags := os.O_CREATE | os.O_APPEND | os.O_WRONLY
	file, err = os.OpenFile(l.filename, flags, 0660)

	if err == nil {
		l.file = file
	}
	return err
}

func (l *logger) Info(format string, msg ...interface{}) {
	m := fmt.Sprintf(format, msg...)
	fmt.Fprintf(l.file, pattern(" "), m)
}

func (l *logger) Error(format string, msg ...interface{}) {
	m := fmt.Sprintf(format, msg...)
	fmt.Fprintf(l.file, pattern("E"), m)
}

func (l *logger) Warn(format string, msg ...interface{}) {
	m := fmt.Sprintf(format, msg...)
	fmt.Fprintf(l.file, pattern("W"), m)
}

func pattern(kind string) string {
	t := time.Now()
	s := " | " + kind + " | %s\n"
	return t.Format(time.RFC3339) + s
}
