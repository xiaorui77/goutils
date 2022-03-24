package logx

import (
	"fmt"
	"github.com/xiaorui77/goutils/logx/formatters"
	"github.com/xiaorui77/goutils/logx/hooks"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

var std = NewLogx("std")
var once sync.Once

type LogX struct {
	Name     string
	instance string

	level        Level
	ReportCaller bool

	Formatter Formatter
	Out       io.Writer
	mu        *sync.Mutex
	hooks     map[Level][]Hook

	entryPool *sync.Pool
}

func NewLogx(name string, opts ...Option) *LogX {
	logger := &LogX{
		Name:         name,
		instance:     name + "-0",
		level:        InfoLevel,
		ReportCaller: false,
		Out:          os.Stdout,
		mu:           new(sync.Mutex),
	}

	logger.Formatter = formatters.NewTextFormatter(logger, true)
	logger.entryPool = &sync.Pool{
		New: func() interface{} {
			return NewEntry(logger)
		},
	}

	// handle options
	for _, o := range opts {
		o(logger)
	}
	return logger
}

func Init(name string, opts ...Option) {
	once.Do(func() {
		if std == nil || std.Name == "std" {
			std = NewLogx(name, opts...)
		}
	})
}

func SetName(name string) {
	std.Name = name
}

func SetInstance(instance string) {
	std.instance = instance
}

func SetLevel(level Level) {
	std.SetLevel(level)
}

func SetLevelS(level string) {
	std.SetLevel(ParseLevel(level))
}

func SetReportCaller(reportCaller bool) {
	std.SetReportCaller(reportCaller)
}

func SetOutput(output io.Writer) {
	std.SetOutput(output)
}

func AddHook(hook Hook) {
	std.AddHook(hook)
}

// Option Pattern functions

type Option func(l *LogX)

func WithInstance(name string) Option {
	return func(l *LogX) {
		l.instance = name
	}
}

func WithLevel(level Level) Option {
	return func(l *LogX) {
		l.SetLevel(level)
	}
}

func WithReportCaller(reportCaller bool) Option {
	return func(l *LogX) {
		l.SetReportCaller(reportCaller)
	}
}

func WithOutput(out io.Writer) Option {
	return func(l *LogX) {
		l.SetOutput(out)
	}
}

func WithHook(hook Hook) Option {
	return func(l *LogX) {
		l.AddHook(hook)
	}
}

func WithEsHook(host string) Option {
	return func(l *LogX) {
		l.AddHook(hooks.NewEsHook(host))
	}
}

// useful methods

func (l *LogX) SetLevel(level Level) {
	l.level = level
}

func (l *LogX) SetReportCaller(reportCaller bool) {
	l.ReportCaller = reportCaller
}

func (l *LogX) SetOutput(out io.Writer) {
	l.Out = out
}

// inner methods

func (l *LogX) IsLevelEnabled(level Level) bool {
	return l.level >= level
}

func (l *LogX) getEntry() *Entry {
	entry, ok := l.entryPool.Get().(*Entry)
	if ok {
		entry.Time = time.Now()
		entry.Fields = make(Fields, 4)
		return entry
	}
	return NewEntry(l)
}

func (l *LogX) releaseEntry(entry *Entry) {
	entry.Fields = nil
	l.entryPool.Put(entry)
}

// Print family functions

func (l *LogX) Log(depth int, level Level, args ...interface{}) {
	if l.IsLevelEnabled(level) {
		entry := l.getEntry()
		entry.Log(depth+1, level, fmt.Sprint(args...))
		l.releaseEntry(entry)
	}
}

func (l *LogX) Debug(args ...interface{}) {
	l.Log(2, DebugLevel, args...)
}

func (l *LogX) Info(args ...interface{}) {
	l.Log(2, InfoLevel, args...)
}

func (l *LogX) Warn(args ...interface{}) {
	l.Log(2, WarnLevel, args...)
}

func (l *LogX) Error(args ...interface{}) {
	l.Log(2, ErrorLevel, args...)
}

func (l *LogX) Fatal(args ...interface{}) {
	l.Log(2, FatalLevel, args...)
	os.Exit(1)
}

func (l *LogX) Panic(args ...interface{}) {
	l.Log(2, PanicLevel, args...)
	panic(fmt.Sprint(args...))
}

// Printf family functions

func (l *LogX) Logf(depth int, level Level, format string, args ...interface{}) {
	if l.IsLevelEnabled(level) {
		entry := l.getEntry()
		entry.Log(depth+1, level, fmt.Sprintf(format, args...))
		l.releaseEntry(entry)
	}
}

func (l *LogX) Debugf(format string, args ...interface{}) {
	l.Logf(2, DebugLevel, format, args...)
}

func (l *LogX) Infof(format string, args ...interface{}) {
	l.Logf(2, InfoLevel, format, args...)
}

func (l *LogX) Warnf(format string, args ...interface{}) {
	l.Logf(2, WarnLevel, format, args...)
}

func (l *LogX) Errorf(format string, args ...interface{}) {
	l.Logf(2, ErrorLevel, format, args...)
}

func (l *LogX) Fatalf(format string, args ...interface{}) {
	l.Logf(2, FatalLevel, format, args...)
	os.Exit(1)
}

func (l *LogX) Panicf(format string, args ...interface{}) {
	l.Logf(2, PanicLevel, format, args...)
	panic(fmt.Sprintf(format, args...))
}

// Global Print family functions

// Log 可以打印指定级别的日志,
// 如果想打印出调用的方法是, 请不要直接使用这个方法, 可以封装一层, 因为它会跳过y
func Log(level Level, args ...interface{}) {
	std.Log(3, level, args...)
}

func Debug(args ...interface{}) { Log(DebugLevel, args...) }

func Info(args ...interface{}) { Log(InfoLevel, args...) }

func Warn(args ...interface{}) { Log(WarnLevel, args...) }

func Error(args ...interface{}) { Log(ErrorLevel, args...) }

func Fatal(args ...interface{}) {
	Log(FatalLevel, args...)
	os.Exit(1)
}

func Panic(args ...interface{}) {
	Log(PanicLevel, args...)
	panic(fmt.Sprint(args...))
}

// Printf family functions

func Logf(level Level, format string, args ...interface{}) {
	std.Logf(3, level, format, args...)
}

func Debugf(format string, args ...interface{}) { Logf(DebugLevel, format, args...) }

func Infof(format string, args ...interface{}) { Logf(InfoLevel, format, args...) }

func Warnf(format string, args ...interface{}) { Logf(WarnLevel, format, args...) }

func Errorf(format string, args ...interface{}) { Logf(ErrorLevel, format, args...) }

func Fatalf(format string, args ...interface{}) {
	Logf(FatalLevel, format, args...)
	os.Exit(1)
}

func Panicf(format string, args ...interface{}) {
	Logf(PanicLevel, format, args...)
	panic(fmt.Sprintf(format, args...))
}

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
