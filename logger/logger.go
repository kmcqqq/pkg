package logger

import (
	"fmt"
	"gitlab.bobbylive.cn/kongmengcheng/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"time"
)

var Log *zap.Logger

// 辅助函数，用于构建日志字段
func String(key string, value string) zap.Field {
	return zap.String(key, value)
}

func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}

func Int64(key string, value int64) zap.Field {
	return zap.Int64(key, value)
}

func Uint(key string, value uint) zap.Field {
	return zap.Uint(key, value)
}

func Float64(key string, value float64) zap.Field {
	return zap.Float64(key, value)
}

func Bool(key string, value bool) zap.Field {
	return zap.Bool(key, value)
}

// Err 返回错误类型的日志字段（重命名以避免与 Error 方法冲突）
func Err(err error) zap.Field {
	return zap.Error(err)
}

func Any(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

func Duration(key string, value time.Duration) zap.Field {
	return zap.Duration(key, value)
}

func InitLogger(cfg *config.LogConfig) error {
	var core zapcore.Core

	// 配置编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.EncodeDuration = zapcore.StringDurationEncoder

	// 根据配置选择编码器
	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 配置输出
	var writeSyncer zapcore.WriteSyncer
	if cfg.Output == "file" {
		// 确保日志目录存在
		if err := os.MkdirAll(filepath.Dir(cfg.File.Path), 0744); err != nil {
			return fmt.Errorf("create log directory failed: %w", err)
		}

		writeSyncer = zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.File.Path,
			MaxSize:    cfg.File.MaxSize, // megabytes
			MaxBackups: cfg.File.MaxBackups,
			MaxAge:     cfg.File.MaxAge, // days
			Compress:   false,
		})
	} else {
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	// 配置日志级别
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 创建核心
	core = zapcore.NewCore(encoder, writeSyncer, level)

	// 创建 logger
	Log = zap.New(
		core,
		zap.AddCaller(),                       // 添加调用者信息
		zap.AddCallerSkip(1),                  // 跳过一层调用栈
		zap.AddStacktrace(zapcore.ErrorLevel), // 错误时添加堆栈跟踪
	)

	// 替换全局 logger
	zap.ReplaceGlobals(Log)

	return nil
}

// Sync 刷新所有缓冲的日志
func Sync() error {
	return Log.Sync()
}

// Debug 使用方便的方式记录调试日志
func Debug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}

// Info 使用方便的方式记录信息日志
func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

// Warn 使用方便的方式记录警告日志
func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}

// Error 使用方便的方式记录错误日志
func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

// Fatal 使用方便的方式记录致命错误日志
func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
}
