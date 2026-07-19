# Vibe Coding Start

这是一个 AI 友好的全栈初始项目模板。后端使用 Go、Gin、GORM、Redis 和 gRPC client/server；前端使用 React、Vite 和 Tailwind CSS。

## 技术栈

- 后端：Go、Gin、GORM、PostgreSQL、Redis、gRPC client/server
- 前端：React、Vite、Tailwind CSS
- 分层：参考 NestJS 的模块化组织方式
- 文档：中文优先

## 目录结构

```text
.
├── cmd/server                 # HTTP 服务入口
├── internal
│   ├── cache                  # Redis client 基础封装
│   ├── config                 # 环境变量配置
│   ├── database               # PostgreSQL 初始化
│   ├── grpcclient             # 出站 gRPC 客户端封装
│   ├── grpcserver             # 入站 gRPC Server 注册与拦截器
│   ├── logger                 # 结构化日志
│   ├── modules/todo           # Todo 业务模块
│   │   ├── entity             # 数据库实体
│   │   ├── request            # HTTP 请求参数结构
│   │   ├── dto                # Service 输入输出与 API 响应结构
│   │   ├── repository         # 数据访问层
│   │   ├── service            # 业务逻辑层
│   │   └── handler            # HTTP 入口层
│   ├── response               # 统一 HTTP 响应
│   └── router                 # 路由注册
├── proto/external             # 外部 gRPC 服务 proto
├── proto/public               # 本服务对外 gRPC proto
├── sql/schema.sql             # 数据库 DDL 权威来源
├── docs                       # 中文架构文档
└── web                        # React 前端项目
```

## 本地启动

先准备 PostgreSQL，然后复制环境变量。Redis 是可用基础设施，只有某个功能明确需要缓存时才接入。

```bash
cp .env.example .env
```

按本机情况修改 `.env` 中的 `DATABASE_DSN`。`HTTP_ADDR` 默认 `:8080`，`GRPC_ADDR` 默认 `:9090`。如果后续功能明确使用 Redis，再配置 `REDIS_ADDR`。前端也从这份根目录 `.env` 读取 `VITE_API_BASE_URL`。

初始化数据库表结构：

```bash
psql "$DATABASE_DSN" -f sql/schema.sql
```

本项目不在服务启动时执行 GORM `AutoMigrate`。数据库表、索引和字段注释统一维护在 `sql/schema.sql`。

启动后端：

```bash
make dev
```

启动前端：

```bash
make web-install
make web-dev
```

## 常用命令

```bash
make test       # 后端测试
make tidy       # 整理 Go 依赖
make proto      # 生成 gRPC client/server 代码
make web-build  # 前端构建
```

`proto/external` 和 `proto/public` 下的 Go 生成代码会提交到版本库。修改 proto 后运行 `make proto` 并一并提交生成文件。

## API

- `GET /healthz`
- `GET /api/v1/todos`
- `POST /api/v1/todos`
- `PATCH /api/v1/todos/:id`
- `DELETE /api/v1/todos/:id`

统一响应格式：

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

## gRPC API

- `public.todo.v1.TodoService/ListTodos`

Todo gRPC 服务监听 `GRPC_ADDR`，`ListTodos` 复用现有 Todo service 的列表查询能力。

## AI 协作约束

- 修改代码前先阅读 `docs/architecture.md`。
- 新业务优先创建 `internal/modules/<module>`。
- HTTP 入参放在 `request`，出参放在 `dto`，数据库实体放在 `entity`。
- `handler` 不写数据库查询，`repository` 不读取 Gin context，`entity` 不直接返回给前端。
- 新增或修改数据库字段时同步更新 `sql/schema.sql`，每个字段都要有注释。
- 不要默认给 CRUD 加 Redis 缓存，只有功能明确要求缓存时才接入。
- 本服务支持出站 gRPC client 和入站 gRPC server；对外 gRPC adapter 只调用 service 层。
