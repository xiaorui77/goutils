package logx

import (
	"io"
	"os"
	"strings"
	"sync"
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

	logger.Formatter = NewTextFormatter(logger, true)
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
