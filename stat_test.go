package fastlog

import (
	"testing"
	"time"

	"github.com/995933447/fastlog/logger"
)

func TestStatReport(t *testing.T) {
	err := InitDefaultCfgLoader("./test/log.toml", &logger.LogConf{
		File: logger.FileLogConf{
			MaxFileSizeBytes:            10000000,
			LogInfoBeforeFileSizeBytes:  -1,
			LogDebugBeforeFileSizeBytes: -1,
			DebugMsgMaxLen:              1024,
			InfoMsgMaxLen:               1024,
			MaxRemainFileNum:            2,
			Level:                       "DBG",
			DefaultLogDir:               "/var/work/logs/fastlog/log",
			BillLogDir:                  "/var/work/logs/fastlog/bill",
			StatLogDir:                  "/var/work/logs/fastlog/stat",
		},
		AlertLevel: "WARN",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = InitDefaultLogger(nil)
	if err != nil {
		t.Fatal(err)
	}

	InitDefaultMsgStat("myService")

	// 处理成功
	ReportStat("LoginRPC", 0, 15*time.Millisecond)

	// 处理失败
	ReportStat("LoginRPC", 1, 30*time.Millisecond)

	// 记录总数
	ReportTotalStat("ActiveConnections", 150)

	OnExit()
}
