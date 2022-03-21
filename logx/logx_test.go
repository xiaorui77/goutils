package logx

import "testing"

func TestPrintf(t *testing.T) {
	Init("test", WithInstance("test"), WithReportCaller(true), WithLevel(DebugLevel), WithLevelS("debug"))

	Debugf("hello world: %s", "debug")
	Infof("hello world: %s", "info")
	Warnf("hello world: %s", "warn")
	Errorf("hello world: %s", "error")
}
