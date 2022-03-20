package logx

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

var std *logX

type logX struct {
	name     string
	instance string

	level        Level
	reportCaller bool

	Formater Formatter
	Out      io.Writer
	mu       *sync.Mutex

	entryPool *sync.Pool
}

func NewLogx(name string) *logX {
	logger := &logX{
		name:         name,
		instance:     name + "-0",
		level:        InfoLevel,
		reportCaller: false,
		Out:          os.Stdout,
		mu:           new(sync.Mutex),
	}

	logger.Formater = &TextFormat{
		logger:   logger,
		colorful: true,
	}
	logger.entryPool = &sync.Pool{
		New: func() interface{} {
			return NewEntry(logger)
		},
	}
	return logger
}

func Init(name string, opts ...Option) {
	if std == nil {
		std = NewLogx(name)
	}
	for _, o := range opts {
		o(std)
	}
}

func SetName(name string) {
	std.name = name
}

func SetInstance(instance string) {
	std.instance = instance
}

func SetLevel(level Level) {
	std.SetLevel(level)
}

func SetReportCaller(reportCaller bool) {
	std.SetReportCaller(reportCaller)
}

func SetOutput(output io.Writer) {
	std.SetOutput(output)
}

func (l *logX) fireHooks() {

}

// Option Pattern

type Option func(l *logX)

func WithInstance(name string) Option {
	return func(l *logX) {
		l.instance = name
	}
}

func WithLevel(level string) Option {
	return func(l *logX) {
		l.SetLevel(ParseLevel(level))
	}
}

func WithReportCaller(reportCaller bool) Option {
	return func(l *logX) {
		l.SetReportCaller(reportCaller)
	}
}

func WithOutput(out io.Writer) Option {
	return func(l *logX) {
		l.SetOutput(out)
	}
}

// useful methods

func (l *logX) SetLevel(level Level) {
	l.level = level
}

func (l *logX) SetReportCaller(reportCaller bool) {
	l.reportCaller = reportCaller
}

func (l *logX) SetOutput(out io.Writer) {
	l.Out = out
}

// inner methods

func (l *logX) IsLevelEnabled(level Level) bool {
	return l.level >= level
}

func (l *logX) getEntry() *Entry {
	entry, ok := l.entryPool.Get().(*Entry)
	if ok {
		entry.Time = time.Now()
		entry.Fields = make(Fields, 4)
		return entry
	}
	return NewEntry(l)
}

func (l *logX) releaseEntry(entry *Entry) {
	entry.Fields = nil
	l.entryPool.Put(entry)
}

// Print family functions

func (l *logX) Log(depth int, level Level, args ...interface{}) {
	if l.IsLevelEnabled(level) {
		entry := l.getEntry()
		entry.Log(depth+1, level, fmt.Sprint(args...))
		l.releaseEntry(entry)
	}
}

func (l *logX) Debug(args ...interface{}) {
	l.Log(2, DebugLevel, args...)
}

func (l *logX) Info(args ...interface{}) {
	l.Log(2, InfoLevel, args...)
}

func (l *logX) Warn(args ...interface{}) {
	l.Log(2, WarnLevel, args...)
}

func (l *logX) Error(args ...interface{}) {
	l.Log(2, ErrorLevel, args...)
}

func (l *logX) Fatal(args ...interface{}) {
	l.Log(2, FatalLevel, args...)
	os.Exit(1)
}

func (l *logX) Panic(args ...interface{}) {
	l.Log(2, PanicLevel, args...)
	panic(fmt.Sprint(args...))
}

// Printf family functions

func (l *logX) Logf(depth int, level Level, format string, args ...interface{}) {
	if l.IsLevelEnabled(level) {
		entry := l.getEntry()
		entry.Log(depth+1, level, fmt.Sprintf(format, args...))
		l.releaseEntry(entry)
	}
}

func (l *logX) Debugf(format string, args ...interface{}) {
	l.Logf(2, DebugLevel, format, args...)
}

func (l *logX) Infof(format string, args ...interface{}) {
	l.Logf(2, InfoLevel, format, args...)
}

func (l *logX) Warnf(format string, args ...interface{}) {
	l.Logf(2, WarnLevel, format, args...)
}

func (l *logX) Errorf(format string, args ...interface{}) {
	l.Logf(2, ErrorLevel, format, args...)
}

func (l *logX) Fatalf(format string, args ...interface{}) {
	l.Logf(2, FatalLevel, format, args...)
	os.Exit(1)
}

func (l *logX) Panicf(format string, args ...interface{}) {
	l.Logf(2, PanicLevel, format, args...)
	panic(fmt.Sprintf(format, args...))
}

// Global Print family functions

func Log(level Level, args ...interface{}) {
	if std == nil {
		std = NewLogx("std")
	}
	std.Log(3, level, args...)
}

func Debug(args ...interface{}) { Log(DebugLevel, args...) }

func Info(args ...interface{}) { Log(InfoLevel, args...) }

func Warn(args ...interface{}) { Log(WarnLevel, args...) }

func Error(args ...interface{}) { Log(ErrorLevel, args...) }

func Panic(args ...interface{}) { Log(PanicLevel, args...) }

func Fatal(args ...interface{}) { Log(FatalLevel, args...) }

// Printf family functions

func Logf(level Level, format string, args ...interface{}) {
	if std == nil {
		std = NewLogx("std")
	}
	std.Logf(3, level, format, args...)
}

func Debugf(format string, args ...interface{}) { Logf(DebugLevel, format, args...) }

func Infof(format string, args ...interface{}) { Logf(InfoLevel, format, args...) }

func Warnf(format string, args ...interface{}) { Logf(WarnLevel, format, args...) }

func Errorf(format string, args ...interface{}) { Logf(ErrorLevel, format, args...) }

func Panicf(format string, args ...interface{}) { Logf(PanicLevel, format, args...) }

func Fatalf(format string, args ...interface{}) { Logf(FatalLevel, format, args...) }

// Utils functions

func ParseLevel(lvl string) Level {
	switch strings.ToLower(lvl) {
	case "panic":
		return PanicLevel
	case "fatal":
		return FatalLevel
	case "error":
		return ErrorLevel
	case "warn", "warning":
		return WarnLevel
	case "info":
		return InfoLevel
	case "debug":
		return DebugLevel
	}
	return InfoLevel
}
