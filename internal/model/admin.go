package model

import (
	"time"
)

// Admin 管理员模型
type Admin struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Status    string    `json:"status" gorm:"default:'active'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AdminUser 管理员视图用户数据
type AdminUser struct {
	ID             string     `json:"id" gorm:"primaryKey"`
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	LastLoginAt    *time.Time `json:"last_login_at"`
	MessageCount   int64      `json:"message_count"`
	SessionCount   int64      `json:"session_count"`
	DailyChatLimit int        `json:"daily_chat_limit"`
}

// AdminModel 管理员视图模型数据
type AdminModel struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name"`
	ModelID      string    `json:"model_id"`
	ProviderID   string    `json:"provider_id"`
	Status       string    `json:"status"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	MessageCount int64     `json:"message_count"`
}

// CreateAdminUserRequest 创建管理员用户请求
type CreateAdminUserRequest struct {
	Name           string `json:"name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=6"`
	Status         string `json:"status" binding:"oneof=active inactive"`
	DailyChatLimit *int   `json:"daily_chat_limit"` // 可选字段，nil表示使用默认值
}

// UpdateAdminUserRequest 更新管理员用户请求
type UpdateAdminUserRequest struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Status         string `json:"status" binding:"oneof=active inactive"`
	DailyChatLimit *int   `json:"daily_chat_limit"`
}

// CreateAdminModelRequest 创建管理员模型请求
type CreateAdminModelRequest struct {
	Name        string `json:"name" binding:"required"`
	ModelID     string `json:"model_id" binding:"required"`
	ProviderID  string `json:"provider_id" binding:"required"`
	Status      string `json:"status" binding:"oneof=active inactive"`
	Description string `json:"description"`
}

// UpdateAdminModelRequest 更新管理员模型请求
type UpdateAdminModelRequest struct {
	Name        string `json:"name"`
	ModelID     string `json:"model_id"`
	ProviderID  string `json:"provider_id"`
	Status      string `json:"status" binding:"oneof=active inactive"`
	Description string `json:"description"`
}

// BatchUserStatusRequest 批量更新用户状态请求
type BatchUserStatusRequest struct {
	UserIDs []string `json:"user_ids" binding:"required,min=1,max=50"`
	Status  string   `json:"status" binding:"required,oneof=active inactive"`
}

// BatchDeleteRequest 批量删除用户请求
type BatchDeleteRequest struct {
	UserIDs []string `json:"user_ids" binding:"required,min=1,max=50"`
}

// AdminAPIProvider 管理员视图的API提供商数据
type AdminAPIProvider struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	ProviderType     string    `json:"provider_type"`
	APIKey           string    `json:"api_key"`
	BaseURL          string    `json:"base_url"`           // API 基础URL
	AnthropicBaseURL string    `json:"anthropic_base_url"` // Anthropic 格式 API 基础URL
	Status           string    `json:"status"`
	Description      string    `json:"description"`
	Priority         int       `json:"priority"`
	CreatedAt        time.Time `json:"created_at"`
}

// CreateAPIProviderRequest 创建API提供商请求
type CreateAPIProviderRequest struct {
	Name             string `json:"name" binding:"required"`
	ProviderType     string `json:"provider_type" binding:"required,oneof=openai anthropic zhipu"`
	APIKey           string `json:"api_key" binding:"required"`
	BaseURL          string `json:"base_url"`
	AnthropicBaseURL string `json:"anthropic_base_url"`
	Status           string `json:"status" binding:"oneof=active inactive"`
	Description      string `json:"description"`
	Priority         int    `json:"priority"`
}

// UpdateAPIProviderRequest 更新API提供商请求
type UpdateAPIProviderRequest struct {
	Name             string `json:"name"`
	ProviderType     string `json:"provider_type" binding:"omitempty,oneof=openai anthropic zhipu"`
	APIKey           string `json:"api_key"`
	BaseURL          string `json:"base_url"`
	AnthropicBaseURL string `json:"anthropic_base_url"`
	Status           string `json:"status" binding:"omitempty,oneof=active inactive"`
	Description      string `json:"description"`
	Priority         int    `json:"priority"`
}

// BatchAPIProviderStatusRequest 批量更新API提供商状态请求
type BatchAPIProviderStatusRequest struct {
	ProviderIDs []string `json:"provider_ids" binding:"required,min=1,max=50"`
	Status      string   `json:"status" binding:"required,oneof=active inactive"`
}

// BatchDeleteAPIProviderRequest 批量删除API提供商请求
type BatchDeleteAPIProviderRequest struct {
	ProviderIDs []string `json:"provider_ids" binding:"required,min=1,max=50"`
}
