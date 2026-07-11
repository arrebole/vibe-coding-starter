package grpcclient

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/arrebole/vibe-coding-starter/internal/config"
	demov1 "github.com/arrebole/vibe-coding-starter/proto/external/demo/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ErrDemoClientNotConfigured = errors.New("外部 demo gRPC 服务未配置")

type DemoClient interface {
	Ping(ctx context.Context, name string) (string, error)
	Close() error
}

type demoClient struct {
	cfg    config.GRPCClientConfig
	log    *slog.Logger
	conn   *grpc.ClientConn
	client demov1.DemoServiceClient
}

func NewDemoClient(cfg config.GRPCClientConfig, log *slog.Logger) DemoClient {
	if cfg.Addr == "" {
		return &demoClient{cfg: cfg, log: log}
	}

	conn, err := grpc.NewClient(
		cfg.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(unaryLogInterceptor(log)),
	)
	if err != nil {
		log.Warn("初始化外部 demo gRPC client 失败", slog.String("addr", cfg.Addr), slog.Any("error", err))
		return &demoClient{cfg: cfg, log: log}
	}

	return &demoClient{
		cfg:    cfg,
		log:    log,
		conn:   conn,
		client: demov1.NewDemoServiceClient(conn),
	}
}

func (c *demoClient) Ping(ctx context.Context, name string) (string, error) {
	if c.cfg.Addr == "" || c.client == nil {
		return "", ErrDemoClientNotConfigured
	}

	ctx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
	defer cancel()

	c.log.Debug("准备调用外部 demo gRPC 服务", slog.String("addr", c.cfg.Addr), slog.String("name", name))
	response, err := c.client.Ping(ctx, &demov1.PingRequest{Name: name})
	if err != nil {
		return "", err
	}
	return response.GetMessage(), nil
}

func (c *demoClient) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func unaryLogInterceptor(log *slog.Logger) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req any,
		reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		startedAt := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		log.Debug("外部 gRPC 调用完成",
			slog.String("method", method),
			slog.Duration("cost", time.Since(startedAt)),
			slog.Any("error", err),
		)
		return err
	}
}
