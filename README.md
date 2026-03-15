# Luban - 智能聊天应用

Luban 是一个基于 Go 语言开发的智能聊天应用，提供与 AI 模型对话的 Web 界面。该应用采用了清晰的分层架构，使用 MySQL 进行数据持久化，Redis 缓存，以及 MinIO 文件存储。

## 特性

- 🚀 基于 Go 的高性能后端
- 💬 支持 AI 模型对话
- 🗃️ MySQL 数据持久化
- ⚡ Redis 缓存支持
- 📁 MinIO 文件存储
- 🔐 GORM ORM 支持
- 🌐 原生 HTML/CSS/JS 前端
- ⚖️ 结构化日志记录

## 快速开始

### 环境要求

- Go 1.19+
- Docker (可选，用于依赖服务)
- MySQL 5.7+
- Redis 6.0+
- MinIO (可选)

### 安装和运行

1. 克隆项目
```bash
git clone https://github.com/yourusername/luban.git
cd luban
```

2. 构建应用
```bash
# 构建应用程序
go build -o luban ./cmd

# 直接运行服务器
go run ./cmd/server.go

# 使用自定义配置运行
go run ./cmd/server.go --config ./config.yaml

# 使用 CLI
./luban server --config ./config.yaml
```

### 使用 Docker Compose 启动依赖

```bash
# 启动所有依赖服务
docker-compose up -d

# 应用将在以下地址可用：
# - HTTP 服务: http://localhost:8080
# - MinIO 控制台: http://localhost:9001
```

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行测试并显示详细输出
go test -v ./...

# 运行特定测试文件
go test ./internal/handler
```

## 项目结构

```
luban/
├── cmd/                    # 应用入口点和 CLI 命令
│   ├── root.go
│   └── server.go
├── internal/              # 私有应用代码
│   ├── handler/           # HTTP 处理器和路由注册
│   ├── model/            # 数据模型和业务逻辑
│   ├── config/           # 配置管理
│   ├── cache/            # Redis 缓存实现
│   ├── storage/          # MinIO 文件存储
│   └── logger/           # 日志配置
├── frontend/             # 前端静态文件（HTML/CSS/JS）
│   ├── css/              # 样式文件
│   ├── js/               # JavaScript 文件
│   ├── images/           # 图片资源
│   ├── index.html        # 主页面
│   ├── build.sh          # 构建脚本
│   └── dev-server.sh     # 开发服务器脚本
├── web/                  # 旧版前端文件（已迁移到 frontend/）
├── config.yaml           # 配置文件
├── docker-compose.yml    # Docker Compose 配置
├── go.mod               # Go 模块文件
└── go.sum               # Go 模块依赖校验
```

## 核心架构

### 核心组件

- **cmd/**: 应用入口点和 CLI 命令 (root.go, server.go)
- **internal/**: 私有应用代码
  - **handler/**: HTTP 处理器和路由注册
  - **model/**: 数据模型和业务逻辑
  - **config/**: 配置管理
  - **cache/**: Redis 缓存实现
  - **storage/**: MinIO 文件存储
  - **logger/**: 日志配置

### 关键组件说明

#### AppDeps 依赖注入模式
依赖通过 `internal/handler/app.go` 中的 `AppDeps` 结构注入：
```go
type AppDeps struct {
    Cfg   *config.Config
    Log   *zap.Logger
    DB    *gorm.DB
    Redis *redis.Client
    MinIO *storage.MinIO
}
```

#### 路由注册
路由在 `registerRoutes()` 方法中注册：
- `/api/chat`: 主要聊天接口 (POST)
- `/api/models`: 获取可用模型列表
- `/api/sessions`: 会话管理
- `/api/upload`: 文件上传到 MinIO

#### 数据模型
应用使用 GORM 模型：
- `Session`: 聊天会话，包含模型 ID 和标题
- `Message`: 单个消息，包含角色（用户/助手）和内容
- `User`: 用户档案（未完全实现）
- `Model`: 可用 AI 模型
- `Attachment`: 上传文件的元数据

## 配置说明

配置通过 Viper 加载，支持：
- YAML 文件 (config.yaml)
- 环境变量（前缀为 `LUBAN_`）
- 缺失设置的默认值

主要配置部分：
- `server`: HTTP 服务器设置
- `mysql`: 数据库连接
- `redis`: 缓存配置
- `minio`: 文件存储设置
- `web`: 静态文件服务

## 数据库架构

应用在启动时会自动创建以下表：
- sessions (id, title, model_id, created_at, updated_at)
- messages (id, session_id, role, content, created_at)
- users (id, name, avatar_url, created_at, updated_at)
- models (id, name, created_at, updated_at)
- attachments (id, bucket, object_key, url, type, created_at)

## 开发注意事项

1. **聊天接口**: 目前返回模拟响应，需要集成实际的 AI 模型 API。
2. **模型管理**: 基础模型列表存在，但未实现真正的模型适配器。
3. **身份验证**: 未实现用户身份验证系统。
4. **前端**: 原生 HTML/CSS/JS 界面，支持开发和生产构建。

## 重要设计模式

1. **依赖注入**: 所有依赖都是注入的，不是在本地创建。
2. **优雅关闭**: 使用上下文和超时进行干净关闭。
3. **结构化日志**: 使用带有结构字段的 Zap 日志记录器。
4. **配置绑定**: 使用 Viper 进行配置，配合 mapstructure 绑定。

## 贡献指南

欢迎提交 Pull Request！请确保：

1. 代码符合项目的代码风格
2. 添加适当的测试
3. 更新必要的文档

## 许可证

[MIT License](LICENSE)