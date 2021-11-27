package logx

import "github.com/sirupsen/logrus"

var std = New()

// Level type
type Level uint32

const (
	TraceLevel Level = iota
	DebugLevel
)

type logX struct {
	logrus *logrus.Logger
}

func New() *logX {
	log := logrus.New()
	log.SetReportCaller(true)
	log.SetFormatter(&stdFormat{})

	return &logX{
		logrus: log,
	}
}

// Print family functions

func Trace(args ...interface{}) { std.logrus.Trace(args) }

func Debug(args ...interface{}) { std.logrus.Debug(args) }

func Print(args ...interface{}) { std.logrus.Print(args) }

func Info(args ...interface{}) { std.logrus.Print(args) }

func Warn(args ...interface{}) { std.logrus.Warn(args) }

func Error(args ...interface{}) { std.logrus.Error(args) }

func Panic(args ...interface{}) { std.logrus.Panic(args) }

func Fatal(args ...interface{}) { std.logrus.Fatal(args) }

// Printf family functions

func Tracef(format string, args ...interface{}) { std.logrus.Tracef(format, args) }

func Debugf(format string, args ...interface{}) { std.logrus.Debugf(format, args) }

func Printf(format string, args ...interface{}) { std.logrus.Printf(format, args) }

func Infof(format string, args ...interface{}) { std.logrus.Infof(format, args) }

func Warnf(format string, args ...interface{}) { std.logrus.Warnf(format, args) }

func Errorf(format string, args ...interface{}) { std.logrus.Errorf(format, args) }

func Panicf(format string, args ...interface{}) { std.logrus.Panicf(format, args) }

func Fatalf(format string, args ...interface{}) { std.logrus.Fatalf(format, args) }
