package logger

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
)

type LogConf struct {
	File       FileLogConf `json:"file"`
	AlertLevel string      `json:"alert_level"`
}

type FileLogConf struct {
	FilePrefix                  string `json:"file_prefix"`
	MaxFileSizeBytes            int64  `json:"max_file_size_bytes"`              // 文件最大容量,字节为单位
	DefaultLogDir               string `json:"default_log_dir"`                  // 默认日志目录
	ExceptionLogDir             string `json:"exception_log_dir"`                // 异常日志目录
	BillLogDir                  string `json:"bill_log_dir"`                     // bill日志目录
	StatLogDir                  string `json:"stat_log_dir"`                     // stat日志目录
	Level                       string `json:"level"`                            // 日志最小级别
	DebugMsgMaxLen              int32  `json:"debug_msg_max_len"`                // debug日志消息最大长度,-1或者0代表不限制
	InfoMsgMaxLen               int32  `json:"info_msg_max_len"`                 // info日志消息最大长度,-1或者0代表不限制
	LogDebugBeforeFileSizeBytes int64  `json:"log_debug_before_file_size_bytes"` // 文件允许写入debug日志的大小阀值,-1代表不限制
	LogInfoBeforeFileSizeBytes  int64  `json:"log_info_before_file_size_bytes"`  // 文件允许写入info日志的大小阀值,-1代表不限制
	FileMaxRemainDays           int    `json:"file_max_remain_days"`             // 文件最大保留天数
	MaxRemainFileNum            int    `json:"max_remain_file_num"`              // 保留文件数量
	CompressFrequentHours       int    `json:"compress_frequent_hours"`          // 压缩频率小时数
	CompressAfterReachBytes     int64  `json:"compress_after_reach_bytes"`       // 压缩最小文件大小
}

func (f *FileLogConf) GetLevel() Level {
	return TransStrToLevel(f.Level)
}

const defaultReloadCfgFileIntervalSec = 10

func NewConfLoader(cfgFile string, reloadCfgFileIntervalSec uint32, defaultLogCfg *LogConf) (*ConfLoader, error) {
	var loader ConfLoader
	if reloadCfgFileIntervalSec <= 0 {
		reloadCfgFileIntervalSec = defaultReloadCfgFileIntervalSec
	}

	loader.reloadCfgFileIntervalSec = reloadCfgFileIntervalSec
	loader.cfgFile = cfgFile

	loader.opLogCfgMu.Lock()
	loader.cfg = defaultLogCfg
	loader.defaultLogCfg = defaultLogCfg
	loader.opLogCfgMu.Unlock()

	if err := loader.loadFile(); err != nil {
		return nil, err
	}

	loader.init()

	return &loader, nil
}

type ConfLoader struct {
	cfgFile                  string
	cfg                      *LogConf
	defaultLogCfg            *LogConf
	opLogCfgMu               sync.RWMutex
	reloadCfgFileIntervalSec uint32
}

func (c *ConfLoader) GetConf() *LogConf {
	c.opLogCfgMu.RLock()
	defer c.opLogCfgMu.RUnlock()
	return c.cfg
}

func (c *ConfLoader) init() {
	go func() {
		for {
			if err := c.loadFile(); err != nil {
				fmt.Println(err)
			}
			time.Sleep(time.Duration(c.reloadCfgFileIntervalSec) * time.Second)
		}
	}()
}

func (c *ConfLoader) loadFile() error {
	if c.cfgFile == "" {
		c.opLogCfgMu.Lock()
		defer c.opLogCfgMu.Unlock()
		c.cfg = c.defaultLogCfg
		return nil
	}

	if _, err := os.Stat(c.cfgFile); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		c.opLogCfgMu.Lock()
		defer c.opLogCfgMu.Unlock()
		c.cfg = c.defaultLogCfg
		return nil
	}

	var cfg LogConf
	if _, err := toml.DecodeFile(c.cfgFile, &cfg); err != nil {
		return err
	}
	c.opLogCfgMu.Lock()
	defer c.opLogCfgMu.Unlock()
	if cfg.File.Level == "" {
		cfg.File.Level = c.defaultLogCfg.File.Level
	}
	if cfg.File.DefaultLogDir == "" {
		cfg.File.DefaultLogDir = c.defaultLogCfg.File.DefaultLogDir
	}
	if cfg.File.BillLogDir == "" {
		cfg.File.BillLogDir = c.defaultLogCfg.File.BillLogDir
	}
	if cfg.File.StatLogDir == "" {
		cfg.File.StatLogDir = c.defaultLogCfg.File.StatLogDir
	}
	c.cfg = &cfg

	return nil
}

func (c *ConfLoader) SetDefaultLogConf(cfg *LogConf) {
	if cfg == nil {
		return
	}
	c.opLogCfgMu.Lock()
	defer c.opLogCfgMu.Unlock()
	if cfg.File.Level == "" {
		cfg.File.Level = c.defaultLogCfg.File.Level
	}
	if cfg.File.DefaultLogDir == "" {
		cfg.File.DefaultLogDir = c.defaultLogCfg.File.DefaultLogDir
	}
	if cfg.File.BillLogDir == "" {
		cfg.File.BillLogDir = c.defaultLogCfg.File.BillLogDir
	}
	if cfg.File.StatLogDir == "" {
		cfg.File.StatLogDir = c.defaultLogCfg.File.StatLogDir
	}
	if cfg.File.ExceptionLogDir == "" {
		cfg.File.ExceptionLogDir = c.defaultLogCfg.File.ExceptionLogDir
	}
	c.defaultLogCfg = cfg
}
