package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/arrebole/vibe-coding-starter/internal/config"
	"github.com/arrebole/vibe-coding-starter/internal/database"
	"github.com/arrebole/vibe-coding-starter/internal/grpcclient"
	"github.com/arrebole/vibe-coding-starter/internal/grpcserver"
	"github.com/arrebole/vibe-coding-starter/internal/logger"
	todohandler "github.com/arrebole/vibe-coding-starter/internal/modules/todo/handler"
	todorepository "github.com/arrebole/vibe-coding-starter/internal/modules/todo/repository"
	todoservice "github.com/arrebole/vibe-coding-starter/internal/modules/todo/service"
	"github.com/arrebole/vibe-coding-starter/internal/router"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.MustLoad()
	log := logger.New(cfg.LogLevel)

	if cfg.AppEnv == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	db := database.MustOpen(cfg.DatabaseDSN)

	demoClient := grpcclient.NewDemoClient(cfg.ExternalDemoGRPC, log)
	defer demoClient.Close()

	todoRepository := todorepository.NewGormRepository(db)
	todoService := todoservice.New(todoRepository, demoClient, log)
	todoHandler := todohandler.New(todoService)

	engine := router.New(router.Dependencies{
		DB:          db,
		TodoHandler: todoHandler,
	})

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	grpcListener, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Error("gRPC 监听端口失败", slog.String("addr", cfg.GRPCAddr), slog.Any("error", err))
		os.Exit(1)
	}
	grpcServer := grpcserver.New(grpcserver.Dependencies{
		TodoService: todoService,
	}, log)

	go func() {
		log.Info("HTTP 服务启动", slog.String("addr", cfg.HTTPAddr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("HTTP 服务异常退出", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	go func() {
		log.Info("gRPC 服务启动", slog.String("addr", cfg.GRPCAddr))
		if err := grpcServer.Serve(grpcListener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Error("gRPC 服务异常退出", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	waitForShutdown(log, server, grpcServer)
}

func waitForShutdown(log *slog.Logger, httpServer *http.Server, grpcServer *grpc.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Info("收到关闭信号，开始优雅停止 HTTP 服务")
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error("HTTP 服务优雅停止失败", slog.Any("error", err))
	}

	log.Info("开始优雅停止 gRPC 服务")
	grpcStopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(grpcStopped)
	}()

	select {
	case <-grpcStopped:
	case <-ctx.Done():
		log.Error("gRPC 服务优雅停止超时", slog.Any("error", ctx.Err()))
		grpcServer.Stop()
	}
}
