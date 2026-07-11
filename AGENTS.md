# Codex 项目约束

## 权威约束来源

本项目以根目录 `CLAUDE.md` 作为所有 AI Agent 的主约束文件。Codex 在开始工作前必须读取并遵守 `CLAUDE.md`；如果本文件与其冲突，以 `CLAUDE.md` 为准。

项目技能统一存放在 `.claude/skills`。需要项目技能时从该目录读取，不要在其他 Agent 目录维护副本。

## 工作规则

- 修改代码前先阅读 `CLAUDE.md` 和 `docs/architecture.md`。
- 保持小步、聚焦的改动，不做无关重构。
- 避免覆盖用户已有改动。
- 后端修改后运行 `go test ./...`。
- 前端修改后运行 `cd web && npm run build`。
