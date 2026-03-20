package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mfb27/luban/internal/auth"
	"github.com/mfb27/luban/internal/middleware"
	"github.com/mfb27/luban/internal/model"
	"github.com/mfb27/luban/internal/response"
)

func (a *App) getUser(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.NewResponseHelper(c).Error(response.CodeNoPermission, "User not authenticated")
		return
	}

	var u model.User
	if err := a.db.First(&u, "id = ?", userID).Error; err != nil {
		response.NewResponseHelper(c).Error(response.CodeNotFound, "User not found")
		return
	}

	response.NewResponseHelper(c).Success(u)
}

func (a *App) register(c *gin.Context) {
	var input auth.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.NewResponseHelper(c).Error(response.CodeInvalidParam, err.Error())
		return
	}

	// Check if user already exists
	var existingUser model.User
	if err := a.db.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		response.NewResponseHelper(c).Error(response.CodeUserExists, "User with this email already exists")
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPasswordBcrypt(input.Password)
	if err != nil {
		response.NewResponseHelper(c).Error(response.CodeInternal, "Failed to hash password")
		return
	}

	// Create user
	user := model.User{
		ID:           uuid.New().String(),
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := a.db.Create(&user).Error; err != nil {
		response.NewResponseHelper(c).Error(response.CodeInternal, "Failed to create user")
		return
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		response.NewResponseHelper(c).Error(response.CodeInternal, "Failed to generate token")
		return
	}

	loginResponse := auth.LoginResponse{
		Token:     token,
		UserID:    user.ID,
		Name:      user.Name,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	response.NewResponseHelper(c).SuccessWithMessage("User registered successfully", loginResponse)
}

func (a *App) login(c *gin.Context) {
	var input auth.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.NewResponseHelper(c).Error(response.CodeInvalidParam, err.Error())
		return
	}

	// Find user
	var user model.User
	if err := a.db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		response.NewResponseHelper(c).Error(response.CodeAuthFailed, "Invalid email or password")
		return
	}

	// Check password
	if !auth.CheckPasswordBcrypt(input.Password, user.PasswordHash) {
		response.NewResponseHelper(c).Error(response.CodeAuthFailed, "Invalid email or password")
		return
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		response.NewResponseHelper(c).Error(response.CodeInternal, "Failed to generate token")
		return
	}

	loginResponse := auth.LoginResponse{
		Token:     token,
		UserID:    user.ID,
		Name:      user.Name,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	response.NewResponseHelper(c).SuccessWithMessage("Login successful", loginResponse)
}

