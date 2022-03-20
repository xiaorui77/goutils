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

// Log is entry point to the log package.
// @param calldepath: An additional call number of lines to skip
func (e *Entry) Log(calldepath int, level Level, msg string) {
	if !e.logger.IsLevelEnabled(level) {
		return
	}

	e.Level = level
	e.Message = msg

	if e.Time.IsZero() {
		e.Time = time.Now()
	}

	if e.logger.reportCaller {
		e.Caller = GetCaller(calldepath + 1)
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
	e.Log(1, DebugLevel, fmt.Sprint(fmt.Sprint(args...)))
}

func (e *Entry) Info(args ...interface{}) {
	e.Log(1, InfoLevel, fmt.Sprint(args...))
}

func (e *Entry) Warn(args ...interface{}) {
	e.Log(1, WarnLevel, fmt.Sprint(args...))
}

func (e *Entry) Error(args ...interface{}) {
	e.Log(1, ErrorLevel, fmt.Sprint(args...))
}

func (e *Entry) Fatal(args ...interface{}) {
	e.Log(1, FatalLevel, fmt.Sprint(args...))
}

func (e *Entry) Panic(args ...interface{}) {
	e.Log(1, PanicLevel, fmt.Sprint(args...))
}

// Printf family functions

func (e *Entry) Debugf(format string, args ...interface{}) {
	e.Log(1, DebugLevel, fmt.Sprintf(format, fmt.Sprint(args...)))
}

func (e *Entry) Infof(format string, args ...interface{}) {
	e.Log(1, InfoLevel, fmt.Sprintf(format, fmt.Sprint(args...)))
}

func (e *Entry) Warnf(format string, args ...interface{}) {
	e.Log(1, WarnLevel, fmt.Sprintf(format, fmt.Sprint(args...)))
}

func (e *Entry) Errorf(format string, args ...interface{}) {
	e.Log(1, ErrorLevel, fmt.Sprintf(format, fmt.Sprint(args...)))
}

func (e *Entry) Fatalf(format string, args ...interface{}) {
	e.Log(1, FatalLevel, fmt.Sprintf(format, fmt.Sprint(args...)))
}

func (e *Entry) Panicf(format string, args ...interface{}) {
	e.Log(1, PanicLevel, fmt.Sprintf(format, fmt.Sprint(args...)))
}

// utils functions

func GetCaller(skip int) *runtime.Frame {
	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(skip+1, pcs)
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
