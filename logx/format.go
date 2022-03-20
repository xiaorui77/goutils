package logx

import (
	"bytes"
	"fmt"
	"github.com/xiaorui77/goutils/timeutils"
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

type TextFormat struct {
	logger   *logX
	colorful bool
}

func (f *TextFormat) Format(entry *Entry) ([]byte, error) {
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

	if f.logger.reportCaller && entry.Caller != nil {
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

func levelColor(level Level) int {
	switch level {
	case InfoLevel:
		return green
	case WarnLevel:
		return yellow
	case ErrorLevel, FatalLevel, PanicLevel:
		return red
	case DebugLevel:
		return gray
	}
	return green
}

func levelString(level Level) string {
	switch level {
	case DebugLevel:
		return "DEBGU"
	case InfoLevel:
		return " INFO"
	case WarnLevel:
		return " WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	case PanicLevel:
		return "PANIC"
	}
	return "UNKNOWN"
}

func buildCaller(entry *Entry) string {
	file := entry.Caller.File
	line := entry.Caller.Line

	if index := strings.LastIndex(file, "/"); index >= 0 {
		file = file[index+1:]
	}
	return fmt.Sprintf("%s:%d", file, line)
}
