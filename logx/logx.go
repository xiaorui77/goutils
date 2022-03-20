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

func newLogx(name string) *logX {
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
			return NewEntry(std)
		},
	}
	return logger
}

func Init(name string, opts ...Option) {
	if std == nil {
		std = newLogx(name)
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
		l.reportCaller = reportCaller
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

func (l *logX) Log(level Level, args ...interface{}) {
	if l.IsLevelEnabled(level) {
		entry := l.getEntry()
		entry.Log(level, args...)
		l.releaseEntry(entry)
	}
}

func (l *logX) Debug(args ...interface{}) {
	l.Log(DebugLevel, args...)
}

func (l *logX) Info(args ...interface{}) {
	l.Log(InfoLevel, args...)
}

func (l *logX) Warn(args ...interface{}) {
	l.Log(WarnLevel, args...)
}

func (l *logX) Error(args ...interface{}) {
	l.Log(ErrorLevel, args...)
}

func (l *logX) Fatal(args ...interface{}) {
	l.Log(FatalLevel, args...)
	os.Exit(1)
}

func (l *logX) Panic(args ...interface{}) {
	l.Log(PanicLevel, args...)
	panic(fmt.Sprint(args...))
}

// Printf family functions

func (l *logX) Logf(level Level, format string, args ...interface{}) {
	if l.IsLevelEnabled(level) {
		entry := l.getEntry()
		entry.Log(level, fmt.Sprintf(format, args...))
		l.releaseEntry(entry)
	}
}

func (l *logX) Debugf(format string, args ...interface{}) {
	l.Logf(DebugLevel, format, args...)
}

func (l *logX) Infof(format string, args ...interface{}) {
	l.Logf(InfoLevel, format, args...)
}

func (l *logX) Warnf(format string, args ...interface{}) {
	l.Logf(WarnLevel, format, args...)
}

func (l *logX) Errorf(format string, args ...interface{}) {
	l.Logf(ErrorLevel, format, args...)
}

func (l *logX) Fatalf(format string, args ...interface{}) {
	l.Logf(FatalLevel, format, args...)
	os.Exit(1)
}

func (l *logX) Panicf(format string, args ...interface{}) {
	l.Logf(PanicLevel, format, args...)
	panic(fmt.Sprintf(format, args...))
}

// Global Print family functions

func Log(level Level, args ...interface{}) {
	if std == nil {
		std = newLogx("std")
	}
	std.Log(level, args...)
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
		std = newLogx("std")
	}
	std.Logf(level, format, args...)
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
