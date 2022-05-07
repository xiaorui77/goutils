package logx

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
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
	Logger *LogX

	Time time.Time

	Fields Fields

	Level Level

	Message string

	Caller *runtime.Frame

	Buffer *bytes.Buffer
}

func NewEntry(l *LogX) *Entry {
	return &Entry{
		Logger: l,
		Time:   time.Now(),
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
	if !e.Logger.IsLevelEnabled(level) {
		return
	}

	e.Level = level
	e.Message = msg

	if e.Time.IsZero() {
		e.Time = time.Now()
	}

	if e.Logger.ReportCaller {
		e.Caller = GetCaller(calldepath + 1)
	}

	// fire hooks
	_ = e.Logger.fireHooks(e.Level, e)

	if e.Logger.Out != nil {
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

	format, err := e.Logger.Formatter.Format(e)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to format logger, %v\n", err)
		return
	}

	e.Logger.mu.Lock()
	defer e.Logger.mu.Unlock()
	if _, err := e.Logger.Out.Write(format); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to write to Output, %v\n", err)
	}
}

func (e *Entry) WithFields(fields Fields) *Entry {
	data := make(Fields, len(e.Fields)+len(fields))
	for k, v := range e.Fields {
		data[k] = v
	}
	for k, v := range fields {
		if t := reflect.TypeOf(v); t != nil &&
			t.Kind() != reflect.Func && t.Kind() != reflect.Ptr && t.Elem().Kind() != reflect.Func {
			data[k] = v
		}
	}
	return &Entry{
		Logger: e.Logger,
		Fields: data,
	}
}

func (e *Entry) WithField(k string, v interface{}) *Entry {
	return e.WithFields(Fields{k: v})
}

// Print functions

func (e *Entry) Debug(args ...interface{}) {
	e.Log(2, DebugLevel, fmt.Sprint(fmt.Sprint(args...)))
}

func (e *Entry) Info(args ...interface{}) {
	e.Log(2, InfoLevel, fmt.Sprint(args...))
}

func (e *Entry) Warn(args ...interface{}) {
	e.Log(1, WarnLevel, fmt.Sprint(args...))
}

func (e *Entry) Error(args ...interface{}) {
	e.Log(2, ErrorLevel, fmt.Sprint(args...))
}

func (e *Entry) Fatal(args ...interface{}) {
	e.Log(2, FatalLevel, fmt.Sprint(args...))
}

func (e *Entry) Panic(args ...interface{}) {
	e.Log(2, PanicLevel, fmt.Sprint(args...))
}

// Printf family functions

func (e *Entry) Debugf(format string, args ...interface{}) {
	e.Log(2, DebugLevel, fmt.Sprintf(format, fmt.Sprint(args...)))
}

func (e *Entry) Infof(format string, args ...interface{}) {
	e.Log(2, InfoLevel, fmt.Sprintf(format, fmt.Sprint(args...)))
}

func (e *Entry) Warnf(format string, args ...interface{}) {
	e.Log(2, WarnLevel, fmt.Sprintf(format, fmt.Sprint(args...)))
}

func (e *Entry) Errorf(format string, args ...interface{}) {
	e.Log(2, ErrorLevel, fmt.Sprintf(format, fmt.Sprint(args...)))
}

func (e *Entry) Fatalf(format string, args ...interface{}) {
	e.Log(2, FatalLevel, fmt.Sprintf(format, fmt.Sprint(args...)))
}

func (e *Entry) Panicf(format string, args ...interface{}) {
	e.Log(2, PanicLevel, fmt.Sprintf(format, fmt.Sprint(args...)))
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
