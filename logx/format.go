package logx

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

const defaultTimeFormat = "2006-01-02 15:04:05"

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
}

func (f *stdFormat) Format(entry *logrus.Entry) ([]byte, error) {
	var buffer *bytes.Buffer
	if entry.Buffer != nil {
		buffer = entry.Buffer
	} else {
		buffer = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format(defaultTimeFormat)
	buffer.WriteString(timestamp)

	colorfulLevel := fmt.Sprintf("\033[%dm%5s\033[0m", getLevelColor(entry.Level), getLevelString(entry.Level))
	buffer.WriteString(" ")
	buffer.WriteString(colorfulLevel)

	caller := buildCaller(entry)
	if caller != "" {
		buffer.WriteString(" ")
		buffer.WriteString(caller)
	}

	buffer.WriteString(" - ")
	buffer.WriteString(entry.Message)
	buffer.WriteString("\n")

	// 2022-02-20 03:27:20 INFO main.go:28 - log info output
	return buffer.Bytes(), nil
}

func getLevelColor(level logrus.Level) int {
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

func getLevelString(level logrus.Level) string {
	switch level {
	case logrus.TraceLevel:
		return "TRACE"
	case logrus.DebugLevel:
		return "DEBUG"
	case logrus.InfoLevel:
		return "INFO"
	case logrus.WarnLevel:
		return "WARN"
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
