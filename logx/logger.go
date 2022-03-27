package logx

import (
	"fmt"
	"io"
	"os"
	"time"
)

func (l *LogX) WithField(key string, value interface{}) *Entry {
	entry := l.getEntry()
	defer l.releaseEntry(entry)
	return entry.WithField(key, value)
}

func (l *LogX) WithFields(fields Fields) *Entry {
	entry := l.getEntry()
	defer l.releaseEntry(entry)
	return entry.WithFields(fields)
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
