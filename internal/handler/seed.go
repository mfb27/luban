package handler

import (
	"errors"

	"github.com/google/uuid"
	"github.com/mfb27/luban/internal/auth"
	"github.com/mfb27/luban/internal/model"
	"gorm.io/gorm"
)

func (a *App) seedIfNeeded() {
	// default user
	var u model.User
	err := a.db.First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Hash admin password
		hashedPassword, _ := auth.HashPasswordBcrypt("admin123")

		_ = a.db.Create(&model.User{
			ID:           uuid.NewString(),
			Name:         "Admin User",
			Email:        "admin@luban.com",
			PasswordHash: hashedPassword,
			AvatarURL:    "https://img.alicdn.com/imgextra/i3/O1CN01QLt9r31b7x4MN6qUL_!!6000000003419-2-tps-116-116.png",
		}).Error
	}

	// models
	var count int64
	_ = a.db.Model(&model.Model{}).Count(&count).Error
	if count == 0 {
		_ = a.db.Create(&[]model.Model{
			{ID: "glm-4-flash", Name: "GLM-4 Flash (快速免费)"},
			{ID: "glm-4", Name: "GLM-4 (通用)"},
			{ID: "glm-4-plus", Name: "GLM-4 Plus (更强能力)"},
		}).Error
	}
}

