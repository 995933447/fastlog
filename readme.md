# 最快的日志框架 FastLog
> **异步写入 / 百万 QPS / 低 CPU 开销 / 精准文件行号 / 内置trace链路追踪**  
> **基准压测优于任意同类日志框架**

一款为极致吞吐与可观测性而生的日志框架：在高并发环境下保持 **百万级 QPS** 的稳定写入，**CPU 占用极低**，同时保留 **完整的调用位置信息（文件/行号/函数）** 与结构化字段，适合对性能极其敏感的场景。

---

## ✨ 特性一览
- **高性能**：单机百万 QPS 级别写入能力（异步、批量、零内存拷贝路径）。
- **低开销**：用户态缓冲 + 无锁 + 无阻塞热路径，显著降低 CPU 使用率与 GC 压力。生产环境实践中单进程一小时已写入30g以上日志，单机已写入1T以上日志，无明显增加服务器cpu负载。
- **真异步**：生产者-消费者模型，单线程/协程批量刷盘，自动背压与丢弃策略。
- **强可观测**：精准 `file:line:function`，纳秒级时间戳，自动注入trace链路追踪。
- **结构化**: print系列函数结构化输出。
- **可落地**：滚动切分（按大小/时间）、压缩归档、磁盘健康检测与告警。
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
	
	err = fastlog.InitDefaultLogger(func(msg *logger.Msg) {
		// TODO 发送告警
	})
	if err != nil {
		t.Fatal(err)
	}
	
	for i := 0; i < 1000000; i++ {
		fastlog.Bill("biz_bill", "hello infof! my name is:fastlog %v", i)
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
- 无锁 + 无伪共享：降低竞争。

- 批量 I/O：减少 syscall & fsync。

- 零拷贝路径：避免频繁分配与 GC。

- 调用信息缓存：近零额外开销。

- 可预测内存：尾延迟更稳。

## 🛡️ 可靠性
- 背压与降级：阻塞、丢弃、退让。

- 完整性：崩溃前 Sync() 最多丢一批。

- 自监控：内置水位、批量大小、flush 时间分位数指标。

- Tracing：自动注入 trace。

# Fastlog 使用示例与说明

本文档演示如何通过 `fastlog` 初始化日志配置、启动默认日志器，并进行高性能日志写入。该示例支持**文件切分、日志级别控制、异步写入、告警回调**等功能。

---

## 📦 初始化步骤

### 1. 导入依赖
```go
import (
    "testing"
    "github.com/995933447/fastlog"
    "github.com/995933447/logger"
)
```

### 2. 加载配置文件 & 初始化配置
#### fastlog.InitDefaultCfgLoader 会从配置文件加载默认配置(非必需，可传空字符串)，同时可以传入结构体设置默认参数。

```go
err := fastlog.InitDefaultCfgLoader("./test/log.toml", &logger.LogConf{
    File: logger.FileLogConf{
    MaxFileSizeBytes:            10000000000,    // 单个日志文件最大字节数（10G）
    LogInfoBeforeFileSizeBytes:  -1,            // 日志文件达到多大前输出Info 级别日志（-1 表示不限制）
    LogDebugBeforeFileSizeBytes: -1,            // 日志文件达到多大前输出Info 级别日志（-1 表示不限制）
    DebugMsgMaxLen:              1024,          // Debug 日志最大消息长度（超过会截断）,0表示不限制
    InfoMsgMaxLen:               1024,          // Info 日志最大消息长度（超过会截断）,0表示不限制
    MaxRemainFileNum:            2,             // 最多保留的日志文件个数
    Level:                       "DBG",         // 日志级别（DBG/INFO/IMP/WARN/ERR/PANIC/FATAL）
    DefaultLogDir:               "/var/work/logs/fastlog/log",  // 普通日志输出路径
    BillLogDir:                  "/var/work/logs/fastlog/bill", // 特殊重要类日志输出路径 
    StatLogDir:                  "/var/work/logs/fastlog/stat", // 统计类日志输出路径
    },
    AlertLevel: "WARN", // 告警级别
})
if err != nil {
    t.Fatal(err)
}
```
## ⚙️ 参数说明

| 参数 | 类型 | 说明                                          |
|------|------|---------------------------------------------|
| `MaxFileSizeBytes` | `int` | 单个日志文件最大字节数，超过后自动切分                         |
| `LogInfoBeforeFileSizeBytes` | `int` | 日志文件达到多大前输出debug 级别日志（-1 表示不限制） |
| `LogDebugBeforeFileSizeBytes` | `int` | 日志文件达到多大前输出Info 级别日志（-1 表示不限制）              |
| `DebugMsgMaxLen` | `int` | Debug 日志最大消息长度（超过会截断）,0表示不限制                |
| `InfoMsgMaxLen` | `int` | Info 日志最大消息长度（超过会截断）,0表示不限制                 |
| `MaxRemainFileNum` | `int` | 最多保留的日志文件个数                                 |
| `Level` | `string` | 最低输出日志级别（DBG/INFO/IMP/WARN/ERR/PANIC/FATAL） |
| `DefaultLogDir` | `string` | 普通日志保存目录                                    |
| `BillLogDir` | `string` | 特殊重要类日志输出路径                                 |
| `StatLogDir` | `string` | 统计日志保存目录                                    |


### 3. 初始化默认日志器并设置告警回调
#### 可以注册一个回调函数，当满足告警条件时执行（例如发送到监控平台）。
```go
err = fastlog.InitDefaultLogger(func(msg *logger.Msg) {
// TODO: 在这里实现告警逻辑，例如发送到告警群。注意这个方法是同步调用的，最好非阻塞逻辑实现
})
if err != nil {
t.Fatal(err)
}
```

### 日志写入示例
#### 批量写入 100 万条 Info 级别日志，适合性能测试或压测。
```go
for i := 0; i < 1_000_000; i++ {
    fastlog.Infof("hello infof! my name is:fastlog %v", i)
}
```

## 📚 常用的 API

以下为 `fastlog` 提供的日志记录 API 列表，涵盖不同日志级别的直接写入、格式化写入、打印模式，以及辅助功能方法。

---

### 1. 基础日志方法
接收任意类型参数（`interface{}`），直接输出指定级别日志。

| 方法 | 说明 |
|------|------|
| `Debug(content interface{})` | 输出 Debug 级别日志 |
| `Info(content interface{})` | 输出 Info 级别日志 |
| `Warn(content interface{})` | 输出 Warn 级别日志 |
| `Important(content interface{})` | 输出 Important 级别日志 |
| `Error(content interface{})` | 输出 Error 级别日志 |
| `Panic(content interface{})` | 输出 Panic 级别日志，并触发 panic |
| `Fatal(content interface{})` | 输出 Fatal 级别日志，并终止程序 |

---

### 2. 格式化日志方法
接收格式化字符串和参数（`fmt.Sprintf` 风格），适合动态拼接日志内容。

| 方法 | 说明 |
|------|------|
| `Debugf(format string, args ...interface{})` | 格式化输出 Debug 日志 |
| `Infof(format string, args ...interface{})` | 格式化输出 Info 日志 |
| `Importantf(format string, args ...interface{})` | 格式化输出 Important 日志 |
| `Warnf(format string, args ...interface{})` | 格式化输出 Warn 日志 |
| `Errorf(format string, args ...interface{})` | 格式化输出 Error 日志 |
| `Panicf(format string, args ...interface{})` | 格式化输出 Panic 日志，并触发 panic |
| `Fatalf(format string, args ...interface{})` | 格式化输出 Fatal 日志，并终止程序 |

---

### 3. Print 系列方法
接收多个参数，内部会进行拼接后输出。

| 方法 | 说明 |
|------|------|
| `PrintDebug(args ...interface{})` | 打印并输出 Debug 日志 |
| `PrintInfo(args ...interface{})` | 打印并输出 Info 日志 |
| `PrintImportant(args ...interface{})` | 打印并输出 Important 日志 |
| `PrintWarn(args ...interface{})` | 打印并输出 Warn 日志 |
| `PrintError(args ...interface{})` | 打印并输出 Error 日志 |
| `PrintPanic(args ...interface{})` | 打印并输出 Panic 日志，并触发 panic |
| `PrintFatal(args ...interface{})` | 打印并输出 Fatal 日志，并终止程序 |

---

### 4. 高级方法

| 方法 | 说明 |
|------|------|
| `WriteBySkipCall(level logger.Level, skipCall int, format string, args ...interface{})` | 在指定调用栈深度处输出日志（可用于封装日志 API 时保持正确的文件/行号） |
| `EnableStdoutPrinter()` | 启用标准输出打印日志 |
| `DisableStdoutPrinter()` | 禁用标准输出打印日志 |

---

### 💡 使用示例
```go
fastlog.Infof("server started at port %d", 8080)

fastlog.PrintError("failed to connect:", err)

fastlog.WriteBySkipCall(logger.LevelInfo, 2, "custom log with correct caller info")
```

### 5. 账单日志 API

`fastlog` 内置了**账单日志专用通道**，支持为不同账单类型独立记录日志文件，并支持回调事件。账单日志可用来保存需要独立区分的重要日志，如三方回调日志，定时作业执行等。

| 方法 | 说明 |
|------|------|
| `OnBill(fn func(billName string))` | 注册账单日志触发回调，当有账单日志写入时调用 |
| `Bill(billName string, format string, args ...interface{})` | 输出账单日志（Important 级别，格式化方式） |
| `PrintBill(billName string, args ...interface{})` | 输出账单日志（Important 级别，直接拼接参数） |
| `BillBySkipCall(skipCall int, billName string, format string, args ...interface{})` | 在指定调用栈深度处输出账单日志（格式化方式） |
| `PrintBillBySkipCall(skipCall int, billName string, args ...interface{})` | 在指定调用栈深度处输出账单日志（直接拼接参数） |

---

### 6. 账单日志实现原理

账单日志通过 `BillLoggerFactory` 按账单名称维护独立的 `logger.Logger` 实例：

- **多账单隔离**：不同 `billName` 对应不同的日志文件。
- **自动创建**：首次调用时会按配置文件中的 `BillLogDir` 自动创建对应日志文件。
- **回调机制**：`OnBill` 注册的回调会在每次账单日志写入后触发，可用于业务侧通知或统计。

```go
// 示例
fastlog.OnBill(func(billName string) {
    fmt.Println("%s日志写入\n", billName)
})

fastlog.Bill("order_success", "订单 %d 支付成功", 12345)
```

# 📊 MsgStat 内置消息统计日志工具

`MsgStat` 是一个基于 `fastlog` 的高性能消息统计工具，支持消息处理的实时统计与周期输出，适用于服务端消息处理、RPC 调用监控等场景。

---

## 🚀 功能特性

- **自动周期输出**：每隔 `Interval`（默认 1 分钟）输出统计结果到文件。
- **多维度统计**：记录消息总数、成功、失败、超时次数及耗时信息。
- **灵活扩展**：支持自定义统计输出回调（`additionMsgStatReportFunc`）,如上报prometheus。
- **异步写入**：基于通道缓存，避免阻塞业务流程。
- **支持文件日志**：按服务名生成独立统计日志文件。

---

## 📦 初始化

```go
import "github.com/995933447/fastlog"

func main() {
    // 初始化默认统计器
    fastlog.InitDefaultMsgStat("myService")
}
```
#### 初始化时会创建一个文件日志器，日志文件名为：{StatLogDir}/msgStat.{srvName}.log

## 📝 核心数据结构

### `ReportStatData`
| 字段 | 类型 | 说明 |
|------|------|------|
| `key` | `string` | 消息标识（如 msgid / rpcname） |
| `reportType` | `int` | 0: 累加统计；1: Set 统计 |
| `result` | `int` | 0: 成功；1/-1: 失败；2/-2: 超时 |
| `processTime` | `time.Duration` | 消息处理耗时 |

---

### `MsgStatData`
| 字段 | 类型 | 说明 |
|------|------|------|
| `Key` | `string` | 消息标识 |
| `Type` | `int32` | 0: 累加统计（输出后清零）；1: 重置型统计（输出后不清零） |
| `TotalMsgNum` | `int64` | 消息处理总数 |
| `SuccessMsgNum` | `int32` | 成功个数 |
| `FailMsgNum` | `int32` | 失败个数 |
| `TimeoutMsgNum` | `int32` | 超时个数 |
| `SumProcessTime` | `time.Duration` | 总处理耗时 |
| `MaxProcessTime` | `time.Duration` | 最大处理耗时 |
| `SumSuccProcessTime` | `time.Duration` | 成功请求总处理耗时 |
| `MaxSuccProcessTime` | `time.Duration` | 成功请求最大处理耗时 |

---

## 📚 API 说明

### 初始化与配置
| 方法 | 说明                       |
|------|--------------------------|
| `InitDefaultMsgStat(srvName string)` | 初始化默认统计器                 |
| `NewMsgStat(svrName string, additionMsgReport additionMsgReportFunc) *MsgStat` | 创建自定义统计器并指定额外统计输出回调      |
| `SetAdditionMsgReport(reportFunc additionMsgReportFunc)` | 设置额外统计回调函数,如上报prometheus |

---

### 数据上报
| 方法 | 说明 |
|------|------|
| `ReportStat(key string, result int, processTime time.Duration)` | 上报单条统计（累加型） |
| `ReportTotalStat(key string, result int)` | 上报总数统计（Set 型） |
| `(*MsgStat) ReportStat(key string, t int, result int, processTime time.Duration)` | 自定义上报类型（0: 累加；1: Set） |

### 使用例子
```
fastlog.InitDefaultMsgStat("myService")

// 处理成功
fastlog.ReportStat("LoginRPC", 0, 15*time.Millisecond)

// 处理失败
fastlog.ReportStat("LoginRPC", 1, 30*time.Millisecond)

// 记录总数
fastlog.ReportTotalStat("ActiveConnections", 150)
```
### 日志文件内容示例：
```
[2025-08-13 14:44:02.0161] [fastlog] [NoTrace][48] IMP github.com/995933447/fastlog.(*MsgStat).RunStat:stat.go:123 =========MsgStat begin=========
[2025-08-13 14:44:02.0161] [fastlog] [NoTrace][48] IMP github.com/995933447/fastlog.(*MsgStat).RunStat:stat.go:123 LoginRPC: Success = 1, Fail = 1, Timeout = 0, Total = 2, MaxTime = 30ms, AvgTime = 22.5ms, TotalTime = 45ms, MaxSuccTime = 15ms, AvgSuccTime = 7.5ms
[2025-08-13 14:44:02.0161] [fastlog] [NoTrace][48] IMP github.com/995933447/fastlog.(*MsgStat).RunStat:stat.go:123 ActiveConnections: Success = 150, Fail = 0, Timeout = 0, Total = 150, MaxTime = 0s, AvgTime = 0s, TotalTime = 0s, MaxSuccTime = 0s, AvgSuccTime = 0s
[2025-08-13 14:44:02.0161] [fastlog] [NoTrace][48] IMP github.com/995933447/fastlog.(*MsgStat).RunStat:stat.go:123 =========MsgStat end=========
