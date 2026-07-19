package grpcserver

import (
	"context"
	"log/slog"
	"time"

	todogrpchandler "github.com/arrebole/vibe-coding-starter/internal/modules/todo/grpchandler"
	todoservice "github.com/arrebole/vibe-coding-starter/internal/modules/todo/service"
	todov1 "github.com/arrebole/vibe-coding-starter/proto/public/todo/v1"
	"google.golang.org/grpc"
)

type Dependencies struct {
	TodoService *todoservice.Service
}

func New(deps Dependencies, log *slog.Logger) *grpc.Server {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(unaryLogInterceptor(log)),
	)

	todov1.RegisterTodoServiceServer(server, todogrpchandler.New(deps.TodoService))
	return server
}

func unaryLogInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		startedAt := time.Now()
		response, err := handler(ctx, req)
		log.Debug("入站 gRPC 调用完成",
			slog.String("method", info.FullMethod),
			slog.Duration("cost", time.Since(startedAt)),
			slog.Any("error", err),
		)
		return response, err
	}
}
