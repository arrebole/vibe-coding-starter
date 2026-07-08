# 项目级 Agent 约束

本文件是本项目所有 AI Agent 约束的权威来源。其他 Agent 入口文件只能引用或兼容本文件，不应复制一份可能漂移的规则。

本项目基于 Go + Gin + GORM + Redis + 出站 gRPC client + React + Tailwind CSS。后端位于仓库根目录，前端位于 `web/`。

## 工作规则

- 默认使用中文沟通和撰写文档。
- 修改后端代码前先阅读 `docs/architecture.md`。
- 保持 NestJS 风格的模块化分层：`entity/request/dto/repository/service/handler`。
- 新增业务模块时放在 `internal/modules/<module>`。
- 不要把数据库实体直接作为 HTTP 响应返回。
- 不要在 handler 中写数据库查询。
- 不要在 repository 中读取 Gin context 或调用外部 gRPC。
- 模块间只能通过 service 层定义的最小接口调用，接口由消费方定义。
- 跨模块依赖必须在 `cmd/server/main.go` 中注入，不要在模块内部自行初始化其他模块依赖。
- 不要跨模块导入对方的 `repository`、`entity`、`request`、`handler`。
- 本服务只调用其他服务的 gRPC 接口，不提供 gRPC server。
- 不要新增 gRPC 监听端口、server 注册或 grpc-gateway。
- proto 生成命令统一维护在 `Makefile` 的 `proto` 目标中。
- 不要新增 `docker-compose.yml`，除非用户明确要求。
- 项目技能统一存放在 `.claude/skills`，其他 Agent 目录不要维护技能副本。

## 验证

- 后端修改后优先运行 `go test ./...`。
- 前端修改后优先运行 `cd web && npm run build`。
