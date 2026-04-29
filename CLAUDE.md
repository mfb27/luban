# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

Luban 是一个基于 Go 语言开发的智能聊天应用，提供与智谱 AI (ZhipuAI) 模型对话的 Web 界面。应用采用 Gin 框架，使用 MySQL 进行数据持久化，Redis 缓存，MinIO 文件存储，并包含完整的后台管理系统。

## 用户端页面目录
D:\space\luban-frontend

## 后台管理页面目录
D:\space\luban-admin

## 常用命令

### 构建和运行
```bash
# 构建应用
go build -o luban ./cmd

# 直接运行服务器
go run ./cmd/server.go

# 使用自定义配置运行
go run ./cmd/server.go --config ./config.yaml

# 使用 CLI
./luban server --config ./config.yaml
```

### 开发环境
```bash
# 启动依赖服务
docker-compose up -d

# 应用地址:
# - HTTP: http://localhost:8080
# - MinIO Console: http://localhost:9001
```

### 测试
```bash
# 运行所有测试
go test ./...

# 详细输出
go test -v ./...

# 运行特定包的测试
go test ./internal/handler

# 运行单个测试
go test -v -run TestFunctionName ./...
```

## 核心架构

### 目录结构
- **cmd/**: 入口点和 CLI 命令 (root.go, server.go)
- **internal/**: 私有应用代码
  - **handler/**: HTTP 处理器和路由注册
  - **model/**: GORM 数据模型
  - **config/**: Viper 配置管理
  - **cache/**: Redis 缓存
  - **storage/**: MinIO 文件存储
  - **logger/**: Zap 日志配置
  - **middleware/**: Gin 中间件 (认证、日志、CORS)
  - **response/**: 统一 API 响应结构
  - **admin/**: 管理员认证服务
  - **auth/**: 用户认证
  - **zhipu/**: 智谱 AI 客户端

### App 结构 (internal/handler/app.go)
依赖通过 `AppDeps` 注入，应用启动时初始化：
```go
type App struct {
    Engine           *gin.Engine
    adminAuthService *admin.AdminAuthService
    cfg              *config.Config
    log              *zap.Logger
    db               *gorm.DB
    redis            *redis.Client
    minio            *storage.MinIO
    zhipu            *zhipu.Client
}
```

### 路由注册
API 路由在 `registerRoutes()` 中注册：
- `/api/health`: 健康检查
- `/api/auth/register`, `/api/auth/login`: 用户认证
- `/api/models`: 获取模型列表 (公开)
- `/api/chat`: 聊天接口 (支持 SSE 流式响应)
- `/api/user/me`: 用户信息 (需认证)
- `/api/sessions`: 会话管理 (需认证)
- `/api/upload`: 文件上传 (需认证)
- `/api/admin/*`: 后台管理 API (需管理员认证)

### 后台管理系统
管理员功能在 `internal/handler/admin.go` 中实现：
- 用户管理: 列表、创建、更新、删除、批量操作
- 模型管理: 列表、创建、更新、删除、批量操作
- JWT 认证: 管理员登录使用 JWT，24 小时过期
- 初始管理员: 首次启动时自动创建 (配置在 config.yaml)

### 数据模型 (GORM)
- `Session`: 聊天会话 (id, user_id, title, model_id)
- `Message`: 消息 (id, session_id, user_id, role, content)
- `User`: 用户 (id, name, email, password_hash, status)
- `Model`: AI 模型 (id, name, model_id, status, description)
- `Attachment`: 文件附件 (id, bucket, object_key, url, type)
- `Admin`: 管理员 (id, name, email, password, status)

### 统一响应结构 (internal/response/response.go)
所有 API 必须使用 `APIResponse` 结构：
```go
type APIResponse struct {
    Code      int         `json:"code"`      // 8位自定义状态码
    Message   string      `json:"message"`
    Data      interface{} `json:"data,omitempty"`
    RequestID string      `json:"requestId"` // 请求追踪 ID
}
```
响应码分类：
- `0`: 成功
- `1xxxxxxx`: 客户端错误
- `2xxxxxxx`: 服务端错误
- `3xxxxxxx`: 认证错误
- `4xxxxxxx`: 权限错误
- `5xxxxxxx`: 业务错误

使用 `response.NewResponseHelper(c)` 简化响应处理。

### 智谱 AI 集成 (internal/zhipu/client.go)
- 支持流式和非流式聊天
- JWT Token 自动生成 (按智谱 AI 规范)
- SSE 流式响应处理

### 配置管理 (config.yaml)
配置通过 Viper 加载，支持环境变量覆盖 (`LUBAN_` 前缀):
- `server`: HTTP 服务器配置
- `mysql`: 数据库 DSN
- `redis`: 缓存配置
- `minio`: 文件存储配置
- `zhipuai`: 智谱 AI API 密钥
- `admin`: 管理员 JWT 密钥和初始账号

## 重要模式

1. **依赖注入**: 所有依赖通过 `AppDeps` 注入，不在本地创建
2. **中间件链**: RequestID → Logger → CORS → Auth (按需)
3. **请求追踪**: 每个请求有唯一 RequestID，日志和响应都包含
4. **软删除**: Model 表使用 GORM 软删除
5. **会话归属**: 匿名会话可转换为认证用户会话
6. **批量操作**: 批量 API 支持事务，失败返回部分结果

## 当前实现状态

- 聊天接口: 完整实现，支持智谱 AI 流式响应
- 模型管理: 完整实现
- 用户管理: 完整实现
- 认证系统: 用户和管理员认证均已实现
- 后台管理: 完整实现，包含批量操作