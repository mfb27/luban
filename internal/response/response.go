package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mfb27/luban/internal/middleware"
)

// APIResponse 统一的API响应结构
type APIResponse struct {
	Code      int         `json:"code"`           // 状态码，8位自定义数字
	Message   string      `json:"message"`        // 响应消息
	Data      interface{} `json:"data,omitempty"` // 响应数据，可为空
	RequestID string      `json:"requestId"`      // 请求ID，用于追踪
}

// ResponseCode 定义8位自定义数字响应码
// 格式分类：
// 0xxxxxxx - 成功（0开头的8位数字）
// 1xxxxxxx - 客户端错误（1开头的8位数字）
// 2xxxxxxx - 服务端错误（2开头的8位数字）
// 3xxxxxxx - 认证相关错误（3开头的8位数字）
// 4xxxxxxx - 权限相关错误（4开头的8位数字）
// 5xxxxxxx - 业务相关错误（5开头的8位数字）

type ResponseCode int

const (
	// 成功类（0开头的8位数字）
	CodeSuccess ResponseCode = 0 // 成功

	// 客户端错误（1开头的8位数字）
	CodeError           ResponseCode = 10000000 // 请求错误
	CodeBadRequest      ResponseCode = 10000001 // 请求参数错误
	CodeInvalidParam    ResponseCode = 10000002 // 无效参数
	CodeRequiredParam   ResponseCode = 10000003 // 缺少必需参数
	CodeNotFound        ResponseCode = 10000004 // 资源未找到
	CodeConflict        ResponseCode = 10000005 // 冲突
	CodeTooManyRequests ResponseCode = 10000006 // 请求过多

	// 服务端错误（2开头的8位数字）
	CodeInternal           ResponseCode = 20000000 // 服务器内部错误
	CodeServiceUnavailable ResponseCode = 20000001 // 服务不可用
	CodeTimeout            ResponseCode = 20000002 // 请求超时
	CodeDatabaseError      ResponseCode = 20000003 // 数据库错误

	// 认证相关错误（3开头的8位数字）
	CodeAuthFailed   ResponseCode = 30000000 // 认证失败
	CodeTokenExpired ResponseCode = 30000001 // token过期
	CodeTokenInvalid ResponseCode = 30000002 // 无效的token
	CodeNoToken      ResponseCode = 30000003 // 缺少token

	// 权限相关错误（4开头的8位数字）
	CodeForbidden    ResponseCode = 40000000 // 禁止访问
	CodeNoPermission ResponseCode = 40000001 // 没有权限
	CodeRoleDenied   ResponseCode = 40000002 // 角色权限不足

	// 业务相关错误（5开头的8位数字）
	CodeBusinessError  ResponseCode = 50000000 // 业务错误
	CodeUserExists     ResponseCode = 50000001 // 用户已存在
	CodeUserNotFound   ResponseCode = 50000002 // 用户不存在
	CodePasswordError  ResponseCode = 50000003 // 密码错误
	CodeSessionExpired ResponseCode = 50000004 // 会话过期
	CodeUploadFailed   ResponseCode = 50000005 // 上传失败
	CodeDownloadFailed ResponseCode = 50000006 // 下载失败
)

// ResponseHelper 响应助手，简化响应处理
type ResponseHelper struct {
	c *gin.Context
}

// NewResponseHelper 创建新的响应助手
func NewResponseHelper(c *gin.Context) *ResponseHelper {
	return &ResponseHelper{c: c}
}

// Success 成功响应
func (r *ResponseHelper) Success(data interface{}) {
	response := APIResponse{
		Code:      int(CodeSuccess),
		Message:   "success",
		Data:      data,
		RequestID: middleware.GetRequestID(r.c),
	}
	r.c.JSON(http.StatusOK, response)
}

// Error 错误响应
func (r *ResponseHelper) Error(code ResponseCode, message string) {
	response := APIResponse{
		Code:      int(code),
		Message:   message,
		RequestID: middleware.GetRequestID(r.c),
	}

	r.c.JSON(http.StatusOK, response)
}

// SuccessWithMessage 带消息的成功响应
func (r *ResponseHelper) SuccessWithMessage(message string, data interface{}) {
	response := APIResponse{
		Code:      int(CodeSuccess),
		Message:   message,
		Data:      data,
		RequestID: middleware.GetRequestID(r.c),
	}
	r.c.JSON(http.StatusOK, response)
}

// ErrorWithDetails 错误响应，带详情
func (r *ResponseHelper) ErrorWithDetails(code ResponseCode, message string, details interface{}) {
	response := APIResponse{
		Code:      int(code),
		Message:   message,
		Data:      details,
		RequestID: middleware.GetRequestID(r.c),
	}

	r.c.JSON(http.StatusOK, response)
}

// PageResponse 分页响应结构
type PageResponse struct {
	Items      interface{} `json:"items"`      // 数据列表
	Total      int64       `json:"total"`      // 总记录数
	Page       int         `json:"page"`       // 当前页码
	PageSize   int         `json:"pageSize"`   // 每页大小
	TotalPages int         `json:"totalPages"` // 总页数
}

// SuccessWithPage 成功响应（分页数据）
func (r *ResponseHelper) SuccessWithPage(items interface{}, total int64, page, pageSize int) {
	response := APIResponse{
		Code:    int(CodeSuccess),
		Message: "success",
		Data: PageResponse{
			Items:      items,
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
		},
		RequestID: middleware.GetRequestID(r.c),
	}
	r.c.JSON(http.StatusOK, response)
}
