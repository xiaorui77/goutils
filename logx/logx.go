package logx

import (
	"github.com/sirupsen/logrus"
	"io"
	"runtime"
)

var std = New()

// Level type
type Level uint32

type logX struct {
	logrus *logrus.Logger
}

func New() *logX {
	log := logrus.New()
	log.SetFormatter(&stdFormat{})

	return &logX{
		logrus: log,
	}
}

func SetLevel(level logrus.Level) {
	std.logrus.SetLevel(level)
}

func SetOutput(output io.Writer) {
	std.logrus.SetOutput(output)
}

// Print family functions

func Log(level logrus.Level, args ...interface{}) {
	_, file, line, _ := runtime.Caller(2)
	std.logrus.WithFields(logrus.Fields{"file": file, "line": line}).Log(level, args...)
}

func Trace(args ...interface{}) { Log(logrus.TraceLevel, args...) }

func Debug(args ...interface{}) { Log(logrus.DebugLevel, args...) }

func Print(args ...interface{}) { Log(logrus.InfoLevel, args...) }

func Info(args ...interface{}) { Log(logrus.InfoLevel, args...) }

func Warn(args ...interface{}) { Log(logrus.WarnLevel, args...) }

func Error(args ...interface{}) { Log(logrus.ErrorLevel, args...) }

func Panic(args ...interface{}) { Log(logrus.PanicLevel, args...) }

func Fatal(args ...interface{}) { Log(logrus.FatalLevel, args...) }

// Printf family functions

func Logf(level logrus.Level, format string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(2)
	std.logrus.WithFields(logrus.Fields{"file": file, "line": line}).Logf(level, format, args...)
}

func Tracef(format string, args ...interface{}) { Logf(logrus.TraceLevel, format, args...) }

func Debugf(format string, args ...interface{}) { Logf(logrus.DebugLevel, format, args...) }

func Printf(format string, args ...interface{}) { Logf(logrus.InfoLevel, format, args...) }

func Infof(format string, args ...interface{}) { Logf(logrus.InfoLevel, format, args...) }

func Warnf(format string, args ...interface{}) { Logf(logrus.WarnLevel, format, args...) }

func Errorf(format string, args ...interface{}) { Logf(logrus.ErrorLevel, format, args...) }

func Panicf(format string, args ...interface{}) { Logf(logrus.PanicLevel, format, args...) }

func Fatalf(format string, args ...interface{}) { Logf(logrus.FatalLevel, format, args...) }
