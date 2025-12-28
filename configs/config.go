package configs

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
)

// LogConfig 日志配置
type LogConfig struct {
	Level        string `toml:"level"`
	Format       string `toml:"format"`
	Output       string `toml:"output"`
	LogDir       string `toml:"log_dir"`
	LogFile      string `toml:"log_file"`
	MaxSize      int    `toml:"max_size"`
	MaxBackups   int    `toml:"max_backups"`
	MaxAge       int    `toml:"max_age"`
	Compress     bool   `toml:"compress"`
	EnableCaller bool   `toml:"enable_caller"`
	CallerSkip   int    `toml:"caller_skip"`
}

// GatewayConfig 网关配置
type GatewayConfig struct {
	URL             string `toml:"url"`
	BeedanceAddress string `toml:"beedance_address"`
}

// Config 全局配置
type Config struct {
	Log     LogConfig     `toml:"log"`
	Gateway GatewayConfig `toml:"gateway"`
}

// GlobalConfig 全局配置实例
var GlobalConfig *Config

var (
	once    sync.Once
	initErr error
)

// InitConfig 初始化配置
func InitConfig(configPath string) error {
	once.Do(func() {
		GlobalConfig = &Config{}

		// 读取配置文件
		data, err := os.ReadFile(configPath)
		if err != nil {
			initErr = fmt.Errorf("读取配置文件失败: %w", err)
			return
		}

		// 解析 TOML 配置
		if err := toml.Unmarshal(data, GlobalConfig); err != nil {
			initErr = fmt.Errorf("解析配置文件失败: %w", err)
			return
		}

		// 验证配置
		if err := validateConfig(GlobalConfig); err != nil {
			initErr = fmt.Errorf("配置验证失败: %w", err)
			return
		}
	})

	return initErr
}

// validateConfig 验证配置的有效性
func validateConfig(cfg *Config) error {
	// 验证日志级别
	validLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true, "fatal": true,
	}
	if !validLevels[cfg.Log.Level] {
		return fmt.Errorf("无效的日志级别: %s", cfg.Log.Level)
	}

	// 验证日志格式
	if cfg.Log.Format != "json" && cfg.Log.Format != "console" {
		return fmt.Errorf("无效的日志格式: %s, 必须是 'json' 或 'console'", cfg.Log.Format)
	}

	// 验证日志输出方式
	if cfg.Log.Output != "stdout" && cfg.Log.Output != "file" && cfg.Log.Output != "both" {
		return fmt.Errorf("无效的日志输出方式: %s, 必须是 'stdout', 'file' 或 'both'", cfg.Log.Output)
	}

	// 验证网关 URL
	if cfg.Gateway.URL == "" {
		return fmt.Errorf("网关 URL 不能为空")
	}

	return nil
}

func (c *LogConfig) GetLogFilePath() string {
	return filepath.Join(c.LogDir, c.LogFile)
}

func (c *LogConfig) EnsureLogDir() error {
	return os.MkdirAll(c.LogDir, 0755)
}
