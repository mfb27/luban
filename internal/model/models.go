package model

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	ID        string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID    string    `gorm:"type:varchar(36);index;not null" json:"user_id"`
	Title     string    `gorm:"type:varchar(200);not null" json:"title"`
	ModelID   string    `gorm:"type:varchar(64);not null" json:"model_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Message struct {
	ID        string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	SessionID string    `gorm:"type:varchar(36);index;not null" json:"session_id"`
	UserID    string    `gorm:"type:varchar(36);index" json:"user_id"`
	Role      string    `gorm:"type:varchar(16);index;not null" json:"role"` // user/assistant
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID             string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name           string    `gorm:"type:varchar(64);not null" json:"name"`
	Email          string    `gorm:"type:varchar(128);uniqueIndex;not null" json:"email"`
	PasswordHash   string    `gorm:"type:varchar(255)" json:"-"`            // Nullable for GitHub users
	GithubID       *string   `gorm:"type:varchar(64);uniqueIndex" json:"-"` // GitHub user ID, nullable (pointer = NULL when not set)
	AvatarURL      string    `gorm:"type:varchar(512)" json:"avatar_url"`
	Status         string    `gorm:"type:varchar(16);default:'active'" json:"status"`
	DailyChatLimit int       `gorm:"type:int;default:-1" json:"daily_chat_limit"` // -1表示无限制，0表示不允许请求，>0表示每日限制次数
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Model struct {
	ID           string         `gorm:"type:varchar(64);primaryKey" json:"id"`
	Name         string         `gorm:"type:varchar(128);not null" json:"name"`
	ModelID      string         `gorm:"type:varchar(128);uniqueIndex" json:"model_id"`
	ProviderID   string         `gorm:"type:varchar(64);index" json:"provider_id"` // 关联到 APIProvider
	Status       string         `gorm:"type:varchar(16);default:'active'" json:"status"`
	Description  string         `gorm:"type:text" json:"description"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete support
}

// APIProvider API提供商配置模型
type APIProvider struct {
	ID               string         `gorm:"type:varchar(64);primaryKey" json:"id"`
	Name             string         `gorm:"type:varchar(128);not null" json:"name"`                 // 显示名称 (如: OpenAI, Anthropic, 智谱AI)
	ProviderType     string         `gorm:"type:varchar(32);not null;default:'openai'" json:"provider_type"` // 提供商类型: openai/anthropic/zhipu
	APIKey           string         `gorm:"type:varchar(256);not null" json:"api_key"`              // API密钥
	BaseURL          string         `gorm:"type:varchar(512)" json:"base_url"`                     // API 基础URL
	AnthropicBaseURL string         `gorm:"type:varchar(512)" json:"anthropic_base_url"`           // Anthropic 格式 API 基础URL (用于 anthropic 类型)
	Status           string         `gorm:"type:varchar(16);default:'active'" json:"status"`       // 状态：active/inactive
	Description      string         `gorm:"type:text" json:"description"`                          // 描述
	Priority         int            `gorm:"type:int;default:0" json:"priority"`                    // 优先级（用于选择提供商）
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"` // 软删除支持
}

// GetBaseURL 根据提供商类型返回对应的 BaseURL
func (p *APIProvider) GetBaseURL() string {
	if p.ProviderType == "anthropic" {
		return p.AnthropicBaseURL
	}
	return p.BaseURL
}

type Attachment struct {
	ID        string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Bucket    string    `gorm:"type:varchar(128);not null" json:"bucket"`
	ObjectKey string    `gorm:"type:varchar(512);not null" json:"object_key"`
	URL       string    `gorm:"type:varchar(1024);not null" json:"url"`
	Type      string    `gorm:"type:varchar(16);index;not null" json:"type"` // image/video
	CreatedAt time.Time `json:"created_at"`
}

// UserDailyChatCount 用户每日对话次数统计
type UserDailyChatCount struct {
	ID        string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID    string    `gorm:"type:varchar(36);index;not null" json:"user_id"`
	Date      string    `gorm:"type:varchar(10);index;not null" json:"date"` // YYYY-MM-DD
	Count     int       `gorm:"type:int;default:0" json:"count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
