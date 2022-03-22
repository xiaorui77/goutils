package demo

import (
	"github.com/xiaorui77/goutils/logx"
	"os"
)

// logx use demo.
func _() {
	logx.Init("demo-name", logx.WithInstance("demo-001"),
		logx.WithLevel(logx.DebugLevel), logx.WithReportCaller(true), logx.WithOutput(os.Stdout))

	logx.Debugf("debug log %s", "debug")
	logx.Infof("info log %s", "info")
}
