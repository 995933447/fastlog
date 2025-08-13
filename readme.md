# 最快的日志框架 FastLog
> **异步写入 / 百万 QPS / 低 CPU 开销 / 精准文件行号 / 内置trace链路追踪**  
> **基准压测优于任意同类日志框架**

一款为极致吞吐与可观测性而生的日志框架：在高并发环境下保持**百万级 QPS** 的稳定写入，**CPU 占用极低**，同时保留**完整的调用位置信息（文件/行号/函数）**与结构化字段，适合对性能极其敏感的场景。

---

## ✨ 特性一览
- **高性能**：单机百万 QPS 级别写入能力（异步、批量、零内存拷贝路径）。
- **低开销**：用户态缓冲 + 无锁 + 无阻塞热路径，显著降低 CPU 使用率与 GC 压力。
- **真异步**：生产者-消费者模型，单线程/单协程批量刷盘，自动背压与丢弃策略。
- **强可观测**：精准 `file:line:function`，纳秒级时间戳，自动注入trace链路追踪。
- **结构化**： print系列函数支持结构化输出。
- **可落地**：滚动切分（按大小/时间）、压缩归档、定期删除、磁盘健康检测与告警。
- **易接入**：简单 API，默认即最佳实践；内置适配常见 APM/Tracing 生态。

---

## 🚀 快速上手（Go）

````
package fastlog

import (
	t "log"
	"testing"

	"github.com/995933447/fastlog/logger"
	"github.com/995933447/fastlog"
)

func BenchmarkLog(b *testing.B) {
	err := fastlog.InitDefaultCfgLoader("./test/log.toml", &logger.LogConf{
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
	err = fastlog.InitDefaultLogger(nil)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		fastlog.Debugf("debug fast log, i:%d", i)
	}
}

func TestInitDefaultLogger(t *testing.T) {
	err := fastlog.InitDefaultCfgLoader("./test/log.toml", &logger.LogConf{
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
	err = fastlog.InitDefaultLogger(nil)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 1000000; i++ {
		fastlog.Infof("hello infof! my name is:fastlog %v", i)
	}

	fastlog.Debug("hello debug! my name is fastlog")
	fastlog.Debugf("hello debugf! my name is:%s", "fastlog")
	fastlog.PrintDebug("print hello PrintDebug", "my name is fastlog")
	fastlog.Info("hello info! my name is fastlog")
	fastlog.Infof("hello infof! my name is:%s", "fastlog")
	fastlog.PrintInfo("hello PrintInfo", "my name is fastlog")
	fastlog.Important("hello important! my name is fastlog")
	fastlog.Importantf("hello importantf! my name is:%s", "fastlog")
	fastlog.PrintImportant("hello PrintImportant", "my name is fastlog")
	fastlog.Warn("hello warn! my name is fastlog")
	fastlog.Warnf("hello warnf! my name is:%s", "fastlog")
	fastlog.PrintWarn("hello PrintWarn", "my name is fastlog")
	fastlog.Error("hello error! my name is fastlog\r\n")
	fastlog.Errorf("hello errorf! my name is:%s", "fastlog\r\n")
	fastlog.PrintError("hello PrintError", "my name is fastlog")
	fastlog.Panic("hello panic! my name is fastlog")
	fastlog.Panicf("hello panicf! my name is:%s", "fastlog")
	fastlog.PrintPanic("hello PrintPanic", "my name is fastlog")
	fastlog.Fatal("hello fatal! my name is fastlog")
	fastlog.Fatalf("hello fatalf! my name is:%s", "fastlog")
	fastlog.PrintFatal("hello PrintFatal", "my name is fastlog")

	OnExit()
}
````
## 输出内容
````
// 时间戳                   模块名称    trace                协程id level  调用函数文件位置                                                  内容
[2025-08-13 10:23:23.1952] [fastlog] [mGwyIPZFWxgWg05vzBYA][35] INFO github.com/995933447/fastlog.TestInitDefaultLogger:log_test.go:58 hello infof! my name is:fastlog 0
````
## 🚀 基准压测
````
goos: darwin
goarch: arm64
pkg: github.com/995933447/fastlog
cpu: Apple M1 Pro
BenchmarkLog
BenchmarkLog-10    	 1732753	       672.4 ns/op
PASS
````

## 🧠 为什么更快
- 锁分离 + 无伪共享：降低竞争。

- 批量 I/O：减少 syscall & fsync。

- 零拷贝路径：避免频繁分配与 GC。

- 调用信息缓存：近零额外开销。

- 可预测内存：尾延迟更稳。

## 🛡️ 可靠性
- 背压与降级：阻塞、丢弃、退让。

- 完整性：崩溃前 Sync() 最多丢一批。

- 自监控：内置水位、批量大小、flush 时间分位数指标。

- Tracing：自动注入 trace。
