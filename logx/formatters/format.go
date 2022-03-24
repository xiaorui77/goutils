package formatters

import (
	"bytes"
	"fmt"
	"github.com/xiaorui77/goutils/coloring"
	"github.com/xiaorui77/goutils/logx"
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

type TextFormatter struct {
	logger   *logx.LogX
	colorful bool
}

func NewTextFormatter(logger *logx.LogX, colorful bool) *TextFormatter {
	return &TextFormatter{
		logger:   logger,
		colorful: colorful,
	}
}

func (f *TextFormatter) Format(entry *logx.Entry) ([]byte, error) {
	var buffer *bytes.Buffer
	if entry.Buffer != nil {
		buffer = entry.Buffer
	} else {
		buffer = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format(timeutils.DefaultTimeFormat)
	buffer.WriteString(timestamp)

	level := coloring.Coloring(levelString(entry.Level), levelColor(entry.Level), f.colorful)
	buffer.WriteString(" ")
	buffer.WriteString(level)
	buffer.WriteString(" - ")
	buffer.WriteString(entry.Message)

	if f.logger.ReportCaller && entry.Caller != nil {
		caller := buildCaller(entry)
		if caller != "" {
			buffer.WriteString(" - ")
			buffer.WriteString(coloring.Coloring(caller, green, f.colorful))
		}
	}

	// 2022-02-20 03:27:20 INFO main.go:28 - log info output
	buffer.WriteString("\n")
	return buffer.Bytes(), nil
}

func levelColor(level logx.Level) int {
	switch level {
	case logx.InfoLevel:
		return green
	case logx.WarnLevel:
		return yellow
	case logx.ErrorLevel, logx.FatalLevel, logx.PanicLevel:
		return red
	case logx.DebugLevel:
		return gray
	}
	return green
}

func levelString(level logx.Level) string {
	switch level {
	case logx.DebugLevel:
		return "DEBGU"
	case logx.InfoLevel:
		return " INFO"
	case logx.WarnLevel:
		return " WARN"
	case logx.ErrorLevel:
		return "ERROR"
	case logx.FatalLevel:
		return "FATAL"
	case logx.PanicLevel:
		return "PANIC"
	}
	return "UNKNOWN"
}

func buildCaller(entry *logx.Entry) string {
	file := entry.Caller.File
	line := entry.Caller.Line

	if index := strings.LastIndex(file, "/"); index >= 0 {
		file = file[index+1:]
	}
	return fmt.Sprintf("%s:%d", file, line)
}
