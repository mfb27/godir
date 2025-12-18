package logger

import (
	"godir/internal/common/svc"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//	type Logger struct {
//		Logger *zap.Logger
//	}
var Logger *zap.SugaredLogger

// LogConfig 日志配置
type LogConfig struct {
	Output string
	Format string
}

// InitWithConfig 根据配置初始化全局 logger
func InitWithConfig(cfg LogConfig) *zap.SugaredLogger {
	// 解析日志级别
	var level zapcore.Level
	switch svc.Cfg().Env {
	case "local":
		level = zapcore.DebugLevel
	default:
		level = zapcore.InfoLevel
	}

	// 配置编码器
	var encoderConfig zapcore.EncoderConfig
	if cfg.Format == "console" {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
	}

	// 设置时间格式
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 创建编码器
	var encoder zapcore.Encoder
	if cfg.Format == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// 配置输出
	var writer zapcore.WriteSyncer
	switch cfg.Output {
	case "stdout":
		writer = zapcore.AddSync(os.Stdout)
	default:
		// 如果指定了文件路径，可以在这里处理
		writer = zapcore.AddSync(os.Stdout)
	}

	// 创建核心
	core := zapcore.NewCore(encoder, writer, level)

	// 创建logger
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)).Sugar()
	return Logger
}

// FromContext 从 context 取出 *zap.Logger，若没有则返回全局 Logger
func FromContext(c *gin.Context) *zap.SugaredLogger {
	val, ok := c.Get("logger")
	if ok {
		if log, ok := val.(*zap.SugaredLogger); ok && log != nil {
			return log
		}
	}
	return Logger
}

// 预留：如需从框架上下文中获取 logger，请在框架层调用 FromContext。

// func Info(msg string, fields ...zap.Field) {
// 	if logger != nil {
// 		logger.Info(msg, fields...)
// 	}
// }

// func Error(msg string, fields ...zap.Field) {
// 	if logger != nil {
// 		logger.Error(msg, fields...)
// 	}
// }

// func Debug(msg string, fields ...zap.Field) {
// 	if logger != nil {
// 		logger.Debug(msg, fields...)
// 	}
// }

// func Warn(msg string, fields ...zap.Field) {
// 	if logger != nil {
// 		logger.Warn(msg, fields...)
// 	}
// }

// func Fatal(msg string, fields ...zap.Field) {
// 	if logger != nil {
// 		logger.Fatal(msg, fields...)
// 	}
// }

// func Sync() {
// 	if logger != nil {
// 		_ = logger.Sync()
// 	}
// }
