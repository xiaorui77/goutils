package logx

import (
	"fmt"
	"io"
	"os"
)

func SetName(name string) {
	std.Name = name
}

func SetInstance(instance string) {
	std.Instance = instance
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

func WithFields(fields Fields) *Entry {
	return std.WithFields(fields)
}

func WithField(key string, value interface{}) *Entry {
	return std.WithField(key, value)
}

func F(key string, value interface{}) *Entry {
	return std.WithField(key, value)
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
