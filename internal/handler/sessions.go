package handler

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mfb27/luban/internal/model"
	"github.com/mfb27/luban/internal/response"
)

type createSessionReq struct {
	Title   string `json:"title"`
	ModelID string `json:"model_id"`
}

func (a *App) listSessions(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.NewResponseHelper(c).Error(response.CodeNoPermission, "unauthorized")
		return
	}

	var sessions []model.Session
	if err := a.db.Where("user_id = ?", userID).Order("updated_at desc").Find(&sessions).Error; err != nil {
		response.NewResponseHelper(c).Error(response.CodeDatabaseError, err.Error())
		return
	}
	response.NewResponseHelper(c).Success(sessions)
}

func (a *App) createSession(c *gin.Context) {
	var req createSessionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.NewResponseHelper(c).Error(response.CodeInvalidParam, "invalid json")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.NewResponseHelper(c).Error(response.CodeNoPermission, "unauthorized")
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		req.Title = "新对话"
	}
	if req.ModelID == "" {
		req.ModelID = "qwen-plus"
	}

	now := time.Now()
	s := model.Session{
		ID:        uuid.NewString(),
		UserID:    userID.(string),
		Title:     req.Title,
		ModelID:   req.ModelID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := a.db.Create(&s).Error; err != nil {
		response.NewResponseHelper(c).Error(response.CodeDatabaseError, err.Error())
		return
	}
	response.NewResponseHelper(c).Success(s)
}

func (a *App) listMessages(c *gin.Context) {
	sessionID := c.Param("id")

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.NewResponseHelper(c).Error(response.CodeNoPermission, "unauthorized")
		return
	}

	var msgs []model.Message
	// For authenticated users, only return their own messages
	if err := a.db.Where("session_id = ? AND user_id = ?", sessionID, userID).Order("created_at asc").Find(&msgs).Error; err != nil {
		response.NewResponseHelper(c).Error(response.CodeDatabaseError, err.Error())
		return
	}
	response.NewResponseHelper(c).Success(msgs)
}

