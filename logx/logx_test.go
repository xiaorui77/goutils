package logx

import (
	"os"
	"testing"
)

func TestPrintf(t *testing.T) {
	Init("test", WithInstance("test"), WithReportCaller(true), WithLevel(DebugLevel))

	Debugf("hello world: %s", "debug")
	Infof("hello world: %s", "info")
	Warnf("hello world: %s", "warn")
	Errorf("hello world: %s", "error")

	Debug("hello world: debug")
	Info("hello world: info")
	Warn("hello world: warn")
	Error("hello world: error")
}

func TestSet(t *testing.T) {
	Init("test", WithInstance("test-set"))

	SetName("test")
	SetInstance("test-set")
	SetLevel(DebugLevel)
	SetLevelS("info")
	SetReportCaller(true)
	SetOutput(os.Stdout)
}
