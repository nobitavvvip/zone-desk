# ZoneDesk

## 简介

ZoneDesk 将 Linux 服务器转变为可通过任何现代浏览器访问的 Web 桌面环境。它无需 X11/Wayland，不依赖图形界面，在裸机上即可运行——你得到一个操作系统的桌面体验。

## 功能特性

- **Web 桌面**：桌面图标、窗口系统（拖拽/最小化/最大化/关闭）、底部任务栏、亮暗主题切换
- **文件管理器**：目录浏览、4 种视图模式、排序、导航历史、快捷目录
- **文件操作**：新建文件夹、上传/下载、重命名、删除、复制/剪切/粘贴（Ctrl+C/X/V）
- **文件预览**：文本文件内联预览、图片流预览
- **键盘快捷键**：F2 重命名、F5 刷新、Delete 删除、Backspace 后退、Enter 打开、Ctrl+A 全选
- **右键菜单**：完整的上下文菜单（打开/预览/下载/复制/剪切/重命名/删除/属性/新建/上传/粘贴/刷新/复制路径）
- **安全操作**：可配置的根目录访问开关、文件大小限制

## 适用场景

- 为团队提供浏览器可访问的共享文件工作区
- 在无显示器的服务器上获得图形化操作体验
- 作为轻量自运维入口，替代传统 SSH + 命令行流程
- 嵌入内网环境作为基础设施文件管理中枢

## 技术栈

| 层 | 技术 |
|------|--------|
| 后端 | Go 1.26 + Gin v1.12 |
| 前端 | Vue 3 + TypeScript + Naive UI + Pinia + Vite 6 |
| 存储 | YAML 配置 + JSON 文件持久化 |
| 部署 | 单二进制 + 静态文件，无外部依赖 |

## 快速开始

下载最新发布包，解压即用：

```bash
tar xzf zonedesk-*.tar.gz
cd zonedesk-*
./start.sh
```

访问 [http://localhost:7070](http://localhost:7070)

## 目录结构

```
zone-panel/
├── backend/                        # Go 后端
│   ├── cmd/server/main.go          # 程序入口
│   ├── internal/
│   │   ├── app/app.go              # 初始化与依赖注入
│   │   ├── config/                 # YAML 配置加载
│   │   ├── controller/             # HTTP 路由与处理器
│   │   ├── model/                  # 数据模型
│   │   ├── repository/             # JSON 文件持久层
│   │   └── service/                # 业务逻辑
│   ├── pkg/util/                   # 路径清理、统一响应
│   ├── go.mod
│   └── go.sum
├── frontend/                       # Vue 前端
│   └── src/
│       ├── App.vue                 # 根组件（主题 + 桌面）
│       ├── main.ts                 # 前端入口
│       ├── api/                    # Axios API 封装
│       ├── assets/styles/          # 全局样式
│       ├── components/             # 通用组件
│       │   ├── desktop/            # 桌面图标、任务栏
│       │   └── filemanager/        # 文件列表、工具栏等
│       ├── store/                  # Pinia 状态管理
│       └── views/                  # 页面级组件
│           ├── desktop/            # 桌面视图
│           └── filemanager/        # 文件管理器窗口
├── scripts/
│   ├── build.sh                    # 构建（Go + Vue）+ 打包 tar.gz
│   ├── dev.sh                      # 本地开发
│   └── deploy/                     # 部署脚本与配置
│       ├── start.sh                # 启动脚本
│       ├── stop.sh                 # 停止脚本
│       └── config.yaml             # 默认配置文件
├── Makefile                        # 快捷命令
└── target/                         # 构建产物
    ├── zonedesk                    # Go 二进制
    ├── dist/                       # 前端构建产物
    └── zonedesk-*.tar.gz           # 发布包
```

## 构建

```bash
# 一键构建 + 打包（Go + Vue + tar.gz）
./scripts/build.sh [版本号]
```

构建产物位于 `target/` 目录。

### 开发模式

前后端热重载：

```bash
./scripts/dev.sh
```

- 后端：`http://localhost:7070`
- 前端：`http://localhost:5173`（自动代理 `/api` 到 7070）

## 配置

`scripts/deploy/config.yaml`：

```yaml
server:
  host: "0.0.0.0"
  port: 7070

fileManager:
  defaultPath: "/"
  maxPreviewSize: 1048576      # 文本预览最大 1MB
  maxUploadSize: 1073741824    # 上传最大 1GB

security:
  enableLogin: false
  allowDangerousOperations: true

ui:
  theme: "dark"                # dark / light
```

所有存储路径均为相对路径，可整体迁移目录。

## 部署到服务器

```bash
# 构建并打包
./scripts/build.sh

# 部署：解压发布包即用
tar xzf target/zonedesk-*.tar.gz
cd zonedesk-*
./start.sh

# 升级：重新构建，解压覆盖
./scripts/build.sh
tar xzf target/zonedesk-*.tar.gz -C /srv/zonedesk --strip-components=1
/srv/zonedesk/stop.sh
/srv/zonedesk/start.sh
```

## API

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/health` | GET | 健康检查 |
| `/api/config` | GET | 获取前端配置 |
| `/api/files/list` | GET | 目录列表（支持排序） |
| `/api/files/stat` | GET | 文件/目录信息 |
| `/api/files/read` | GET | 读取文本文件 |
| `/api/files/download` | GET | 下载文件 |
| `/api/files/upload` | POST | 上传文件 |
| `/api/files/mkdir` | POST | 新建目录 |
| `/api/files/rename` | POST | 重命名 |
| `/api/files/delete` | POST | 删除 |
| `/api/files/copy` | POST | 复制（支持递归） |
| `/api/files/move` | POST | 移动 |
| `/api/shortcuts` | GET | 快捷目录列表 |
| `/api/shortcuts` | POST | 添加快捷目录 |
| `/api/shortcuts/:id` | DELETE | 删除快捷目录 |

## 版本规划

| 版本 | 功能 |
|------|------|
| v0.1 | Web 桌面 + 文件管理器 ✅ |
| v0.2 | 登录 + Markdown 日报 |
| v0.3 | Docker 容器管理 |
| v0.4 | 系统状态 + 端口管理 |

## 参与贡献

欢迎提交 Issue 和 Pull Request。详见 [CONTRIBUTING.md](CONTRIBUTING.md)。

## 许可

Apache 2.0。详见 [LICENSE](LICENSE)。
