package fastlog

import (
	t "log"
	"testing"

	"github.com/995933447/fastlog/logger"
)

func BenchmarkLog(b *testing.B) {
	err := InitDefaultCfgLoader("./test/log.toml", &logger.LogConf{
		File: logger.FileLogConf{
			MaxFileSizeBytes:            10000000,
			LogInfoBeforeFileSizeBytes:  -1,
			LogDebugBeforeFileSizeBytes: -1,
			DebugMsgMaxLen:              15,
			MaxRemainFileNum:            2,
			Level:                       "DBG",
			DefaultLogDir:               "/var/work/logs/fastlog/log",
			BillLogDir:                  "/var/work/logs/fastlog/bill",
			StatLogDir:                  "/var/work/logs/fastlog/stat",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	err = InitDefaultLogger(nil)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		Debugf("debug fast log, i:%d", i)
	}
}

func TestInitDefaultLogger(t *testing.T) {
	err := InitDefaultCfgLoader("./test/log.toml", &logger.LogConf{
		File: logger.FileLogConf{
			MaxFileSizeBytes:            10000000,
			LogInfoBeforeFileSizeBytes:  -1,
			LogDebugBeforeFileSizeBytes: -1,
			DebugMsgMaxLen:              15,
			MaxRemainFileNum:            2,
			Level:                       "DBG",
			DefaultLogDir:               "/var/work/logs/fastlog/log",
			BillLogDir:                  "/var/work/logs/fastlog/bill",
			StatLogDir:                  "/var/work/logs/fastlog/stat",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	err = InitDefaultLogger(nil)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 1000000; i++ {
		Infof("hello infof! my name is:fastlog %v\r\n", i)
	}

	Debug("hello debug! my name is fastlog")
	Debugf("hello debugf! my name is:%s", "fastlog")
	PrintDebug("print hello PrintDebug", "my name is fastlog")
	Info("hello info! my name is fastlog")
	Infof("hello infof! my name is:%s", "fastlog")
	PrintInfo("hello PrintInfo", "my name is fastlog")
	Important("hello important! my name is fastlog")
	Importantf("hello importantf! my name is:%s", "fastlog")
	PrintImportant("hello PrintImportant", "my name is fastlog")
	Warn("hello warn! my name is fastlog")
	Warnf("hello warnf! my name is:%s", "fastlog")
	PrintWarn("hello PrintWarn", "my name is fastlog")
	Error("hello error! my name is fastlog\r\n")
	Errorf("hello errorf! my name is:%s", "fastlog\r\n")
	PrintError("hello PrintError", "my name is fastlog")
	Panic("hello panic! my name is fastlog")
	Panicf("hello panicf! my name is:%s", "fastlog")
	PrintPanic("hello PrintPanic", "my name is fastlog")
	Fatal("hello fatal! my name is fastlog")
	Fatalf("hello fatalf! my name is:%s", "fastlog")
	PrintFatal("hello PrintFatal", "my name is fastlog")

	OnExit()
}
