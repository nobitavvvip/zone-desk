# 贡献指南

感谢你考虑为 ZoneDesk 贡献代码！

## 报告问题

- 使用 GitHub Issues 报告 bug 或功能请求
- 描述问题时请提供：环境信息、复现步骤、期望行为和实际行为

## 提交 Pull Request

1. Fork 本仓库
2. 创建功能分支：`git checkout -b feat/your-feature`
3. 提交改动
4. 确保构建通过：`make build`
5. 发起 Pull Request

## 开发环境

```bash
# 安装依赖
make frontend

# 开发模式（前后端热重载）
make dev
```

## 代码规范

- Go 代码使用 `gofmt` 格式化
- 前端代码遵循项目内已有的风格
- 保持注释简洁，中英文均可
