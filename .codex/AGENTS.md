# Codex 项目约束

## 权威约束来源

本项目以 `.claude/CLAUDE.md` 作为所有 AI Agent 的主约束文件。项目技能统一存放在 `.claude/skills`。Codex 必须先读取并遵守 `.claude/CLAUDE.md`，本文件只提供 Codex 兼容入口和执行补充。

如果本文件与 `.claude/CLAUDE.md` 冲突，以 `.claude/CLAUDE.md` 为准。

## Codex 补充规则

- 修改代码前先阅读 `.claude/CLAUDE.md` 和 `docs/architecture.md`。
- 需要项目技能时读取 `.claude/skills`，不要在 `.codex/skills` 维护副本。
- 保持小步、聚焦的改动，不做无关重构。
- 避免覆盖用户已有改动。
- 后端修改后运行 `go test ./...`。
- 前端修改后运行 `cd web && npm run build`。
