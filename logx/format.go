package logx

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/yougtao/goutils/timeutils"
	"strings"
)

const (
	black  = 30
	red    = 31
	green  = 32
	yellow = 33
	blue   = 34
	purple = 35
	cyan   = 36
	gray   = 37
)

type stdFormat struct {
	colorful bool
	caller   bool
}

func (f *stdFormat) Format(entry *logrus.Entry) ([]byte, error) {
	var buffer *bytes.Buffer
	if entry.Buffer != nil {
		buffer = entry.Buffer
	} else {
		buffer = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format(timeutils.DefaultTimeFormat)
	buffer.WriteString(timestamp)

	level := coloring(levelString(entry.Level), levelColor(entry.Level), f.colorful)
	buffer.WriteString(" ")
	buffer.WriteString(level)
	buffer.WriteString(" - ")
	buffer.WriteString(entry.Message)

	if f.caller {
		caller := buildCaller(entry)
		if caller != "" {
			buffer.WriteString(" - ")
			buffer.WriteString(coloring(caller, green, f.colorful))
		}
	}

	// 2022-02-20 03:27:20 INFO main.go:28 - log info output
	buffer.WriteString("\n")
	return buffer.Bytes(), nil
}

func levelColor(level logrus.Level) int {
	switch level {
	case logrus.InfoLevel:
		return green
	case logrus.WarnLevel:
		return yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return red
	case logrus.DebugLevel, logrus.TraceLevel:
		return gray
	}
	return green
}

func levelString(level logrus.Level) string {
	switch level {
	case logrus.TraceLevel:
		return "TRACE"
	case logrus.DebugLevel:
		return "DEBGU"
	case logrus.InfoLevel:
		return " INFO"
	case logrus.WarnLevel:
		return " WARN"
	case logrus.ErrorLevel:
		return "ERROR"
	case logrus.FatalLevel:
		return "FATAL"
	case logrus.PanicLevel:
		return "PANIC"
	}
	return "UNKNOWN"
}

func buildCaller(entry *logrus.Entry) string {
	file := "xx"
	line := -1
	if fi, ok := entry.Data["file"]; ok {
		if f, ok := fi.(string); ok {
			file = f
		}
	}
	if file == "xx" {
		return ""
	}
	if li, ok := entry.Data["line"]; ok {
		if l, ok := li.(int); ok {
			line = l
		}
	}

	if index := strings.LastIndex(file, "/"); index >= 0 {
		file = file[index+1:]
	}
	return fmt.Sprintf("%s:%d", file, line)
}
