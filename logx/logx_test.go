package logx

import "testing"

func TestPrintf(t *testing.T) {
	Init(OptReportCaller(true))
	Printf("hello world: %s", "hello")
	Debugf("hello world: %s", "debug")
	Infof("hello world: %s", "info")
	Warnf("hello world: %s", "warn")
	Errorf("hello world: %s", "error")
}
