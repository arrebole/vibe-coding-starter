# 架构说明

本项目参考 NestJS 的模块化分层，但保留 Go 和 Gin 的常见命名习惯。每个业务模块放在 `internal/modules/<module>` 下，并在模块内维护完整分层。

## 分层职责

### entity

`entity` 是数据库实体层，只描述 GORM model、表名、字段、索引和持久化细节。数据库实体不能直接作为 HTTP 响应返回给前端。

### request

`request` 是 HTTP 请求参数层，只定义 Gin bind 使用的 body、query、path 参数结构。这里可以写基础校验 tag，但不能写业务规则。

### dto

`dto` 是数据传输层，定义 Service 输入输出、API 响应结构，以及 entity 到 response 的转换函数。DTO 用来隔离数据库结构和对外接口结构。

### repository

`repository` 是数据访问层，只负责通过 GORM 访问 PostgreSQL。它不读取 Gin context，不返回 HTTP 状态码，不调用外部 gRPC 服务。

### service

`service` 是业务逻辑层，负责业务校验、编排 repository、Redis 缓存和外部 gRPC client。Service 不依赖 Gin，也不直接拼 HTTP 响应。

### handler

`handler` 是 HTTP 适配层，负责绑定 request、调用 service、转换错误、返回统一 JSON 响应。Handler 不写数据库查询和复杂业务规则。

## 调用方向

```text
handler -> service -> repository -> entity
                 └-> cache
                 └-> grpcclient
```

禁止反向依赖，也禁止跨层直接调用。例如 handler 不能直接调用 repository，repository 不能调用 grpcclient。

## 模块间调用

`internal/modules` 下的模块之间只能通过 service 层定义的最小接口调用。跨模块调用必须由消费方 service 定义自己需要的接口，再在 `cmd/server/main.go` 中注入提供方 service 的实现。

允许的方向：

```text
moduleA/service -> moduleB/service interface
```

禁止的方向：

- 跨模块直接访问对方的 `repository`、`entity`、`request`、`handler`。
- handler 直接调用其他模块。
- repository 调用其他模块。
- 两个 service 互相依赖形成循环。
- 在模块内部自行初始化其他模块的 repository 或 service。

示例：订单模块需要读取用户信息时，由订单 service 定义最小接口。

```go
type UserReader interface {
	GetByID(ctx context.Context, id uint) (dto.UserResponse, error)
}

type Service struct {
	repository repository.Repository
	userReader UserReader
}
```

然后在应用入口组装依赖：

```go
userService := userservice.New(userRepository)
orderService := orderservice.New(orderRepository, userService)
```

这样可以让提供方模块保留数据库实体和数据访问细节，消费方模块只依赖自己真正需要的业务能力。

## gRPC 约定

本服务只消费其他服务的 gRPC 接口，不提供 gRPC server。

- 外部 proto 放在 `proto/external/<service>/v1`。
- 生成命令统一使用 `make proto`。
- 出站连接、超时、日志 interceptor 放在 `internal/grpcclient`。
- 业务模块只能在 service 层调用 gRPC client。

## 数据库 DDL 维护

数据库表结构的权威来源是 `sql/schema.sql`。该文件只存放表的创建语句、索引和注释，不存放测试数据、初始化业务数据或环境配置。

- 新增表、字段或索引时，必须同步更新 `sql/schema.sql`。
- 每个字段必须使用 `COMMENT ON COLUMN` 写明含义。
- 每张表应使用 `COMMENT ON TABLE` 写明用途。
- 后端启动不使用 GORM `AutoMigrate` 自动建表。
- `entity` 必须与 `sql/schema.sql` 保持一致，但不能替代 DDL。

## 新增业务模块

新增模块时复制 Todo 的结构：

```text
internal/modules/<module>/
├── entity
├── request
├── dto
├── repository
├── service
└── handler
```

然后在 `internal/router` 中注册模块路由，在 `cmd/server/main.go` 中完成依赖注入。
