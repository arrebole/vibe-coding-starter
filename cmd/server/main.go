package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/arrebole/vibe-coding-starter/internal/cache"
	"github.com/arrebole/vibe-coding-starter/internal/config"
	"github.com/arrebole/vibe-coding-starter/internal/database"
	"github.com/arrebole/vibe-coding-starter/internal/grpcclient"
	"github.com/arrebole/vibe-coding-starter/internal/logger"
	todohandler "github.com/arrebole/vibe-coding-starter/internal/modules/todo/handler"
	todorepository "github.com/arrebole/vibe-coding-starter/internal/modules/todo/repository"
	todoservice "github.com/arrebole/vibe-coding-starter/internal/modules/todo/service"
	"github.com/arrebole/vibe-coding-starter/internal/router"
)

func main() {
	cfg := config.MustLoad()
	log := logger.New(cfg.LogLevel)

	if cfg.AppEnv == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	db := database.MustOpen(cfg.DatabaseDSN)

	redisClient := cache.NewRedisClient(cfg.Redis)
	defer closeRedis(log, redisClient)

	demoClient := grpcclient.NewDemoClient(cfg.ExternalDemoGRPC, log)
	defer demoClient.Close()

	todoRepository := todorepository.NewGormRepository(db)
	todoCache := cache.NewTodoCache(redisClient, cfg.Redis.TodoTTL)
	todoService := todoservice.New(todoRepository, todoCache, demoClient, log)
	todoHandler := todohandler.New(todoService)

	engine := router.New(router.Dependencies{
		DB:          db,
		Redis:       redisClient,
		TodoHandler: todoHandler,
	})

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Info("HTTP 服务启动", slog.String("addr", cfg.HTTPAddr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("HTTP 服务异常退出", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	waitForShutdown(log, server)
}

func closeRedis(log *slog.Logger, client *redis.Client) {
	if err := client.Close(); err != nil {
		log.Warn("关闭 Redis 连接失败", slog.Any("error", err))
	}
}

func waitForShutdown(log *slog.Logger, server *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Info("收到关闭信号，开始优雅停止 HTTP 服务")
	if err := server.Shutdown(ctx); err != nil {
		log.Error("HTTP 服务优雅停止失败", slog.Any("error", err))
	}
}
