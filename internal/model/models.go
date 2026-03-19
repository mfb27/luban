package model

import "time"


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
	ID           string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(64);not null" json:"name"`
	Email        string    `gorm:"type:varchar(128);uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"` // Hide from JSON
	AvatarURL    string    `gorm:"type:varchar(512)" json:"avatar_url"`
	Status       string    `gorm:"type:varchar(16);default:'active'" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Model struct {
	ID          string    `gorm:"type:varchar(64);primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(128);not null" json:"name"`
	ModelID     string    `gorm:"type:varchar(128);uniqueIndex" json:"model_id"`
	Status      string    `gorm:"type:varchar(16);default:'active'" json:"status"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Attachment struct {
	ID        string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Bucket    string    `gorm:"type:varchar(128);not null" json:"bucket"`
	ObjectKey string    `gorm:"type:varchar(512);not null" json:"object_key"`
	URL       string    `gorm:"type:varchar(1024);not null" json:"url"`
	Type      string    `gorm:"type:varchar(16);index;not null" json:"type"` // image/video
	CreatedAt time.Time `json:"created_at"`
}


