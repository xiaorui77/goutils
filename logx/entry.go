package logx

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	packageName        = "github.com/xiaorui77/goutils/logx"
	minimumCallerDepth = 4
	maximumCallerDepth = 25
)

type Entry struct {
	logger *logX

	Time time.Time

	Fields Fields

	Level Level

	Message string

	Caller *runtime.Frame

	Buffer *bytes.Buffer
}

func NewEntry(l *logX) *Entry {
	return &Entry{
		logger: l,
		Fields: make(Fields, 4),
	}
}

var bufferPool = &sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func (e *Entry) Log(level Level, args ...interface{}) {
	if e.logger.IsLevelEnabled(level) {
		e.log(level, fmt.Sprint(args...))
	}
}

func (e *Entry) log(level Level, msg string) {
	e.Level = level
	e.Message = msg

	if e.Time.IsZero() {
		e.Time = time.Now()
	}

	if e.logger.reportCaller {
		e.Caller = getCaller()
	}

	// fire hooks
	e.logger.fireHooks()

	if e.logger.Out != nil {
		e.write()
	}
}

func (e *Entry) write() {
	buffer := bufferPool.Get().(*bytes.Buffer)
	defer func() {
		e.Buffer = nil
		buffer.Reset()
		bufferPool.Put(buffer)
	}()
	buffer.Reset()
	e.Buffer = buffer

	format, err := e.logger.Formater.Format(e)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to format logger, %v\n", err)
		return
	}

	e.logger.mu.Lock()
	defer e.logger.mu.Unlock()
	if _, err := e.logger.Out.Write(format); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to write to Output, %v\n", err)
	}
}

// Print functions

func (e *Entry) Debug(args ...interface{}) {
	e.Log(DebugLevel, args...)
}

func (e *Entry) Info(args ...interface{}) {
	e.Log(InfoLevel, args...)
}

func (e *Entry) Warn(args ...interface{}) {
	e.Log(WarnLevel, args...)
}

func (e *Entry) Error(args ...interface{}) {
	e.Log(ErrorLevel, args...)
}

func (e *Entry) Fatal(args ...interface{}) {
	e.Log(FatalLevel, args...)
}

func (e *Entry) Panic(args ...interface{}) {
	e.Log(PanicLevel, args...)
}

// Printf family functions

func (e *Entry) Debugf(format string, args ...interface{}) {
	e.Log(DebugLevel, fmt.Sprintf(format, args...))
}

func (e *Entry) Infof(format string, args ...interface{}) {
	e.Log(InfoLevel, fmt.Sprintf(format, args...))
}

func (e *Entry) Warnf(format string, args ...interface{}) {
	e.Log(WarnLevel, fmt.Sprintf(format, args...))
}

func (e *Entry) Errorf(format string, args ...interface{}) {
	e.Log(ErrorLevel, fmt.Sprintf(format, args...))
}

func (e *Entry) Fatalf(format string, args ...interface{}) {
	e.Log(FatalLevel, fmt.Sprintf(format, args...))
}

func (e *Entry) Panicf(format string, args ...interface{}) {
	e.Log(PanicLevel, fmt.Sprintf(format, args...))
}

// utils functions

func getCaller() *runtime.Frame {
	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if pkg != packageName {
			return &f //nolint:scopelint
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}
