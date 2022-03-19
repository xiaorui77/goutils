package logx

import (
	"io"
	"runtime"

	"github.com/sirupsen/logrus"
)

var std = New()

type logX struct {
	reportCaller bool

	logrus *logrus.Logger
}

func New() *logX {
	log := logrus.New()
	log.SetFormatter(&stdFormat{
		colorful: true,
		caller:   true,
	})

	return &logX{
		logrus: log,
	}
}

func Init(opts ...Option) {
	for _, o := range opts {
		o(std)
	}
}

func SetLevel(level logrus.Level) {
	std.logrus.SetLevel(level)
}

func SetOutput(output io.Writer) {
	std.logrus.SetOutput(output)
}

// Option Pattern

type Option func(l *logX)

func OptLevel(level string) Option {
	return func(l *logX) {
		lv, err := logrus.ParseLevel(level)
		if err != nil {
			// default Info
			lv = logrus.InfoLevel
		}
		l.logrus.SetLevel(logrus.Level(lv))
	}
}

func OptReportCaller(reportCaller bool) Option {
	return func(l *logX) {
		l.reportCaller = reportCaller
	}
}

func OptOutput() Option {
	return func(l *logX) {
		l.logrus.SetOutput(l.logrus.Out)
	}
}

// Print family functions

func Log(level logrus.Level, args ...interface{}) {
	var field logrus.Fields
	if std.reportCaller {
		_, file, line, _ := runtime.Caller(2)
		field = logrus.Fields{"file": file, "line": line}
	}
	std.logrus.WithFields(field).Log(level, args...)
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
	var field logrus.Fields
	if std.reportCaller {
		_, file, line, _ := runtime.Caller(2)
		field = logrus.Fields{"file": file, "line": line}
	}
	std.logrus.WithFields(field).Logf(level, format, args...)
}

func Tracef(format string, args ...interface{}) { Logf(logrus.TraceLevel, format, args...) }

func Debugf(format string, args ...interface{}) { Logf(logrus.DebugLevel, format, args...) }

func Printf(format string, args ...interface{}) { Logf(logrus.InfoLevel, format, args...) }

func Infof(format string, args ...interface{}) { Logf(logrus.InfoLevel, format, args...) }

func Warnf(format string, args ...interface{}) { Logf(logrus.WarnLevel, format, args...) }

func Errorf(format string, args ...interface{}) { Logf(logrus.ErrorLevel, format, args...) }

func Panicf(format string, args ...interface{}) { Logf(logrus.PanicLevel, format, args...) }

func Fatalf(format string, args ...interface{}) { Logf(logrus.FatalLevel, format, args...) }
