# Vibe Template Skill

当需要在本项目中新增业务能力、API、前端页面或外部 gRPC client 时使用本技能。

## 新增后端模块

1. 在 `internal/modules/<module>` 下创建：
   - `entity`
   - `request`
   - `dto`
   - `repository`
   - `service`
   - `handler`
2. `entity` 只写 GORM model。
3. `request` 只写 HTTP 入参结构和基础 bind 校验。
4. `dto` 写 Service 输入输出、API 响应结构和转换函数。
5. `repository` 只访问数据库。
6. `service` 编排业务逻辑、缓存和外部 gRPC client。
7. `handler` 只做 Gin 参数绑定、错误转换和响应。
8. 在 `internal/router` 注册路由，在 `cmd/server/main.go` 注入依赖。

## 模块间调用

1. 模块之间只能通过 service 层定义的最小接口调用。
2. 跨模块接口由消费方 service 定义，只包含当前场景需要的方法。
3. 提供方 service 实现该接口，不要额外暴露 repository 或 entity。
4. 在 `cmd/server/main.go` 中完成依赖注入。
5. 不要在模块内部自行初始化其他模块依赖。
6. 不要导入其他模块的 `repository`、`entity`、`request`、`handler`。

## 新增外部 gRPC client

1. proto 放在 `proto/external/<service>/v1`。
2. 在 `Makefile` 的 `proto` 目标中维护生成命令。
3. 在 `internal/grpcclient` 中封装连接、超时、日志和关闭逻辑。
4. 只允许业务 service 调用 gRPC client。
5. 不要新增 gRPC server。

## 新增前端页面

1. 页面和 API 调用放在 `web/src`。
2. 使用 Tailwind CSS，保持工具型界面风格。
3. API 地址从 `VITE_API_BASE_URL` 读取。
4. 修改后运行 `cd web && npm run build`。
