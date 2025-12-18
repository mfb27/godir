package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"godir/internal/common/ginx"
	"godir/internal/common/jwt"
	"godir/internal/common/logger"
	"godir/internal/common/svc"
	"godir/internal/handler"

	"go.uber.org/zap"
)

var configFile = flag.String("c", "config/local.yml", "配置文件路径")

func main() {
	flag.Parse()

	// 加载配置
	serviceContext, err := svc.Init(*configFile)
	if err != nil {
		fmt.Printf("加载配置文件失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	log := logger.InitWithConfig(logger.LogConfig{
		Output: serviceContext.Cfg.Log.Output,
		Format: serviceContext.Cfg.Log.Format,
	})

	log.With(zap.String("config", *configFile), zap.Int64("port", serviceContext.Cfg.Server.Port)).Info("应用启动")

	// 初始化JWT配置
	if err := jwt.Init(serviceContext.Cfg.JWT.SecretKey, serviceContext.Cfg.JWT.TokenExp); err != nil {
		log.Error("JWT初始化失败", zap.String("error", err.Error()))
		os.Exit(1)
	}

	engine := ginx.New(log)
	handler.RegisterRouter(engine)

	log.Info("服务器启动中", zap.String("address", fmt.Sprintf(":%d", serviceContext.Cfg.Server.Port)))

	go func() {
		if err := engine.Run(fmt.Sprintf(":%d", serviceContext.Cfg.Server.Port)); err != nil {
			log.Error("服务器启动失败", zap.String("error", err.Error()))
			os.Exit(2)
		}
	}()

	// graceful shutdown on SIGINT/SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("收到停止信号，正在关闭服务器...")
	log.Sync() // 确保所有日志都被写入
	log.Info("server stopped")
}
