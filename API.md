# Luban API 文档

## 统一响应结构

所有 API 响应都使用统一的 `APIResponse` 结构：

```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "requestId": "uuid"
}
```

### 字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| `code` | int | 状态码，8位数字 |
| `message` | string | 响应消息 |
| `data` | object/array/null | 响应数据，成功时包含业务数据 |
| `requestId` | string | 请求ID，用于追踪 |

### 状态码说明

| 代码范围 | 说明 |
|---------|------|
| 0 | 成功 |
| 10000000-19999999 | 客户端错误 |
| 20000000-29999999 | 服务端错误 |
| 30000000-39999999 | 认证相关错误 |
| 40000000-49999999 | 权限相关错误 |
| 50000000-59999999 | 业务相关错误 |

### 常用状态码

| Code | 说明 |
|------|------|
| 0 | 成功 |
| 10000001 | 请求参数错误 |
| 10000002 | 无效参数 |
| 10000003 | 缺少必需参数 |
| 10000004 | 资源未找到 |
| 20000000 | 服务器内部错误 |
| 20000003 | 数据库错误 |
| 30000000 | 认证失败 |
| 30000003 | 缺少token |
| 40000000 | 禁止访问 |
| 40000001 | 没有权限 |
| 50000001 | 用户已存在 |
| 50000002 | 用户不存在 |
| 50000003 | 密码错误 |
| 50000005 | 上传失败 |

## 公共接口

### 1. 健康检查

**GET** `/api/health`

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "ok": true
  },
  "requestId": "uuid"
}
```

### 2. 获取模型列表

**GET** `/api/models`

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": "uuid",
      "name": "GLM-4 Flash",
      "model_id": "glm-4-flash",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "requestId": "uuid"
}
```

## 认证接口

### 3. 用户注册

**POST** `/api/auth/register`

**请求体：**
```json
{
  "name": "用户名",
  "email": "user@example.com",
  "password": "password123"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "User registered successfully",
  "data": {
    "token": "jwt_token",
    "user_id": "uuid",
    "name": "用户名",
    "email": "user@example.com",
    "avatar_url": null,
    "expires_at": "2024-01-02T00:00:00Z"
  },
  "requestId": "uuid"
}
```

### 4. 用户登录

**POST** `/api/auth/login`

**请求体：**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "Login successful",
  "data": {
    "token": "jwt_token",
    "user_id": "uuid",
    "name": "用户名",
    "email": "user@example.com",
    "avatar_url": null,
    "expires_at": "2024-01-02T00:00:00Z"
  },
  "requestId": "uuid"
}
```

## 用户接口（需要认证）

### 5. 获取当前用户信息

**GET** `/api/user/me`

**请求头：**
```
Authorization: Bearer {token}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "uuid",
    "name": "用户名",
    "email": "user@example.com",
    "avatar_url": "https://...",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "requestId": "uuid"
}
```

## 会话接口（需要认证）

### 6. 获取会话列表

**GET** `/api/sessions`

**请求头：**
```
Authorization: Bearer {token}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": "uuid",
      "title": "新对话",
      "model_id": "glm-4-flash",
      "user_id": "uuid",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "requestId": "uuid"
}
```

### 7. 创建会话

**POST** `/api/sessions`

**请求头：**
```
Authorization: Bearer {token}
```

**请求体：**
```json
{
  "title": "新对话",
  "model_id": "glm-4-flash"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "uuid",
    "title": "新对话",
    "model_id": "glm-4-flash",
    "user_id": "uuid",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "requestId": "uuid"
}
```

### 8. 获取会话消息列表

**GET** `/api/sessions/:id/messages`

**请求头：**
```
Authorization: Bearer {token}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": "uuid",
      "session_id": "uuid",
      "user_id": "uuid",
      "role": "user",
      "content": "你好",
      "created_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid",
      "session_id": "uuid",
      "user_id": "uuid",
      "role": "assistant",
      "content": "你好！有什么我可以帮助你的吗？",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "requestId": "uuid"
}
```

## 聊天接口（支持SSE流式响应）

### 9. 发送消息

**POST** `/api/chat`

**请求头：**
```
Authorization: Bearer {token} (可选，未登录时也可使用)
Content-Type: application/json
```

**请求体：**
```json
{
  "session_id": "uuid",
  "content": "你好",
  "model_id": "glm-4-flash",
  "attachment_urls": []
}
```

**响应格式：Server-Sent Events (SSE)**

该接口使用 SSE 流式响应，不返回标准的 APIResponse 格式。

**事件类型：**

| 事件 | 数据 | 说明 |
|------|------|------|
| session | `{"session_id": "uuid", "user_msg_id": "uuid"}` | 会话创建/确认 |
| content | `{"data": "文本内容"}` | 流式内容 |
| done | `{"assistant_id": "uuid", "session_id": "uuid", "is_new_conversation": false}` | 响应完成 |
| error | `{"error": "错误信息"}` | 错误发生 |

**SSE 响应示例：**
```
data: {"type":"session","session_id":"uuid","user_msg_id":"uuid"}

data: {"type":"content","data":"你好"}

data: {"type":"done","assistant_id":"uuid","session_id":"uuid","is_new_conversation":false}
```

## 文件上传接口（需要认证）

### 10. 上传文件

**POST** `/api/upload`

**请求头：**
```
Authorization: Bearer {token}
Content-Type: multipart/form-data
```

**请求体：**
```
file: (binary file)
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "uuid",
    "url": "https://minio.example.com/bucket/uuid.ext"
  },
  "requestId": "uuid"
}
```

## 管理员接口

### 11. 管理员登录

**POST** `/api/admin/login`

**请求体：**
```json
{
  "email": "admin@example.com",
  "password": "admin_password"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "admin_token",
    "admin": {
      "id": "admin",
      "email": "admin@example.com"
    }
  },
  "requestId": "uuid"
}
```

### 12. 获取管理员信息

**GET** `/api/admin/me`

**请求头：**
```
Authorization: Bearer {admin_token}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "admin",
    "email": "admin@example.com",
    "name": "管理员"
  },
  "requestId": "uuid"
}
```

### 13. 获取用户列表

**GET** `/api/admin/users`

**请求头：**
```
Authorization: Bearer {admin_token}
```

**查询参数：**
| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| search | string | - | 搜索用户名或邮箱 |
| status | string | - | 状态筛选 (active/inactive) |
| page | int | 1 | 页码 |
| page_size | int | 20 | 每页数量 |

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "users": [
      {
        "id": "uuid",
        "name": "用户名",
        "email": "user@example.com",
        "status": "active",
        "created_at": "2024-01-01T00:00:00Z",
        "last_login_at": "2024-01-01T00:00:00Z",
        "message_count": 10,
        "session_count": 2
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 20
  },
  "requestId": "uuid"
}
```

### 14. 创建用户

**POST** `/api/admin/users`

**请求头：**
```
Authorization: Bearer {admin_token}
```

**请求体：**
```json
{
  "name": "用户名",
  "email": "user@example.com",
  "password": "password123",
  "status": "active"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "uuid",
    "name": "用户名",
    "email": "user@example.com",
    "status": "active",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "requestId": "uuid"
}
```

### 15. 更新用户

**PUT** `/api/admin/users/:id`

**请求头：**
```
Authorization: Bearer {admin_token}
```

**请求体：**
```json
{
  "name": "新用户名",
  "status": "active",
  "password": "new_password"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "uuid",
    "name": "新用户名",
    "email": "user@example.com",
    "status": "active",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "requestId": "uuid"
}
```

### 16. 切换用户状态

**PATCH** `/api/admin/users/:id/status`

**请求头：**
```
Authorization: Bearer {admin_token}
```

**请求体：**
```json
{
  "status": "inactive"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "status": "inactive"
  },
  "requestId": "uuid"
}
```

### 17. 删除用户

**DELETE** `/api/admin/users/:id`

**请求头：**
```
Authorization: Bearer {admin_token}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "user deleted successfully"
  },
  "requestId": "uuid"
}
```

### 18. 获取模型列表（管理员）

**GET** `/api/admin/models`

**请求头：**
```
Authorization: Bearer {admin_token}
```

**查询参数：**
| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| page | int | 1 | 页码 |
| page_size | int | 20 | 每页数量 |

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "models": [
      {
        "id": "uuid",
        "name": "GLM-4 Flash",
        "model_id": "glm-4-flash",
        "status": "active",
        "description": "快速响应模型",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 20
  },
  "requestId": "uuid"
}
```

### 19. 创建模型

**POST** `/api/admin/models`

**请求头：**
```
Authorization: Bearer {admin_token}
```

**请求体：**
```json
{
  "name": "GLM-4 Flash",
  "model_id": "glm-4-flash",
  "status": "active",
  "description": "快速响应模型"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "uuid",
    "name": "GLM-4 Flash",
    "model_id": "glm-4-flash",
    "status": "active",
    "description": "快速响应模型",
    "created_at": "2024-01-01T00:00:00Z"
  },
  "requestId": "uuid"
}
```

### 20. 更新模型

**PUT** `/api/admin/models/:id`

**请求头：**
```
Authorization: Bearer {admin_token}
```

**请求体：**
```json
{
  "name": "新模型名",
  "model_id": "new-model-id",
  "status": "active",
  "description": "模型描述"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "uuid",
    "name": "新模型名",
    "model_id": "new-model-id",
    "status": "active",
    "description": "模型描述",
    "created_at": "2024-01-01T00:00:00Z"
  },
  "requestId": "uuid"
}
```

### 21. 切换模型状态

**PATCH** `/api/admin/models/:id/status`

**请求头：**
```
Authorization: Bearer {admin_token}
```

**请求体：**
```json
{
  "status": "inactive"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "status": "inactive"
  },
  "requestId": "uuid"
}
```

### 22. 删除模型

**DELETE** `/api/admin/models/:id`

**请求头：**
```
Authorization: Bearer {admin_token}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "model deleted successfully"
  },
  "requestId": "uuid"
}
```

## 前端响应处理规范

### 成功响应处理

```javascript
// api 函数已自动提取 data 字段
const data = await api('/api/user/me');
// data 直接是 APIResponse.data 的内容，无需再访问 data.data
```

### 错误响应处理

```javascript
try {
  const data = await api('/api/user/me');
} catch (error) {
  // error.message 包含 APIResponse.message
  console.error(error.message);
}
```

### 特殊情况：/api/chat SSE 接口

该接口使用 Server-Sent Events 流式响应，不返回标准的 APIResponse 格式，需要特殊处理。
