package logx

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

const defaultTimeFormat = "2006-01-02 15:04:05"

const (
	red    = 31
	green  = 32
	yellow = 33
	blue   = 36
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
	colorfulLevel := fmt.Sprintf("\u001B[%dm%s\u001B[0m", getLevelColor(entry.Level), strings.ToUpper(entry.Level.String()))

	log := fmt.Sprintf("%s %5s %s %s\n", timestamp, colorfulLevel, buildCaller(entry.Caller), entry.Message)
	buffer.WriteString(log)
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
	default:
		return green
	}
}

func buildCaller(caller *runtime.Frame) string {
	if caller == nil {
		return ""
	}
	fileName := caller.File
	if index := strings.LastIndex(fileName, "/"); index >= 0 {
		fileName = fileName[index+1:]
	}
	return fmt.Sprintf("%s:%d", fileName, caller.Line)
}
