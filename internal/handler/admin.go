package handler

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mfb27/luban/internal/admin"
	"github.com/mfb27/luban/internal/model"
	"github.com/mfb27/luban/internal/response"
	"gorm.io/gorm"
)

// CreateAdminMiddleware 创建管理员认证中间件
func CreateAdminMiddleware(authService *admin.AdminAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.NewResponseHelper(c).Error(response.CodeNoToken, "authorization header is required")
			c.Abort()
			return
		}

		// 去掉 "Bearer " 前缀
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		claims, err := authService.ValidateToken(token)
		if err != nil {
			response.NewResponseHelper(c).Error(response.CodeTokenInvalid, "invalid token")
			c.Abort()
			return
		}

		// 将管理员信息存入上下文
		c.Set("admin_id", claims.UserID)
		c.Set("admin_email", claims.Email)
		c.Next()
	}
}

// Admin App
type AdminApp struct {
	authService *admin.AdminAuthService
	db          *gorm.DB
}

// NewAdminApp 创建管理员应用
func NewAdminApp(authService *admin.AdminAuthService, db *gorm.DB) *AdminApp {
	return &AdminApp{
		authService: authService,
		db:          db,
	}
}

// adminLogin 管理员登录
func (a *AdminApp) adminLogin(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NewResponseHelper(c).Error(response.CodeInvalidParam, "invalid request")
		return
	}

	token, err := a.authService.Login(req.Email, req.Password)
	if err != nil {
		switch err {
		case admin.ErrAdminNotFound:
			response.NewResponseHelper(c).Error(response.CodeNotFound, "admin not found")
		case admin.ErrInvalidPassword:
			response.NewResponseHelper(c).Error(response.CodeInvalidParam, "invalid password")
		case admin.ErrAdminDisabled:
			response.NewResponseHelper(c).Error(response.CodeForbidden, "admin account is disabled")
		default:
			response.NewResponseHelper(c).Error(response.CodeInternal, "login failed")
		}
		return
	}

	response.NewResponseHelper(c).Success(gin.H{
		"token": token,
		"admin": gin.H{
			"id":    "admin", // 简化处理
			"email": req.Email,
		},
	})
}

// adminGetMe 获取当前管理员信息
func (a *AdminApp) adminGetMe(c *gin.Context) {
	adminID := c.GetString("admin_id")
	adminEmail := c.GetString("admin_email")

	response.NewResponseHelper(c).Success(gin.H{
		"id":    adminID,
		"email": adminEmail,
		"name":  "管理员", // 简化处理
	})
}

// adminGetUsers 获取用户列表
func (a *AdminApp) adminGetUsers(c *gin.Context) {
	var users []model.AdminUser

	// 构建查询
	query := a.db.Model(&model.User{}).
		Select(`
			u.id, u.name, u.email, u.created_at, u.status,
			(SELECT MAX(created_at) FROM messages WHERE user_id = u.id) as last_login_at,
			(SELECT COUNT(*) FROM messages WHERE user_id = u.id) as message_count,
			(SELECT COUNT(DISTINCT session_id) FROM messages WHERE user_id = u.id) as session_count
		`).Table("users as u")

	// 搜索条件
	if search := c.Query("search"); search != "" {
		query = query.Where("name LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 状态过滤
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	query = query.Offset(offset).Limit(pageSize)

	if err := query.Scan(&users).Error; err != nil {
		response.NewResponseHelper(c).Error(response.CodeDatabaseError, "failed to load users")
		return
	}

	// 获取总数
	var total int64
	a.db.Model(&model.User{}).Count(&total)

	response.NewResponseHelper(c).Success(gin.H{
		"users": users,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}

// adminCreateUser 创建用户
func (a *AdminApp) adminCreateUser(c *gin.Context) {
	var req model.CreateAdminUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NewResponseHelper(c).Error(response.CodeInvalidParam, "invalid request")
		return
	}

	// 检查邮箱是否已存在
	var existingUser model.User
	if err := a.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		response.NewResponseHelper(c).Error(response.CodeInvalidParam, "email already exists")
		return
	}

	// 创建用户
	user := model.User{
		ID:           generateUserID(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: req.Password, // 简化处理，实际应该加密
	}

	if err := a.db.Create(&user).Error; err != nil {
		response.NewResponseHelper(c).Error(response.CodeDatabaseError, "failed to create user")
		return
	}

	response.NewResponseHelper(c).Success(user)
}

// adminUpdateUser 更新用户
func (a *AdminApp) adminUpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var req model.UpdateAdminUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NewResponseHelper(c).Error(response.CodeInvalidParam, "invalid request")
		return
	}

	// 检查用户是否存在
	var user model.User
	if err := a.db.First(&user, "id = ?", userID).Error; err != nil {
		response.NewResponseHelper(c).Error(response.CodeNotFound, "user not found")
		return
	}

	// 更新用户 - 只更新非零值字段
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.Password != "" {
		updates["password_hash"] = req.Password // 简化处理，实际应该加密
	}

	if len(updates) == 0 {
		response.NewResponseHelper(c).Error(response.CodeInvalidParam, "no fields to update")
		return
	}

	if err := a.db.Model(&user).Updates(updates).Error; err != nil {
		response.NewResponseHelper(c).Error(response.CodeDatabaseError, "failed to update user")
		return
	}

	// 重新获取更新后的用户数据
	a.db.First(&user, "id = ?", userID)
	response.NewResponseHelper(c).Success(user)
}

// adminToggleUserStatus 切换用户状态
func (a *AdminApp) adminToggleUserStatus(c *gin.Context) {
	userID := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NewResponseHelper(c).Error(response.CodeInvalidParam, "invalid request")
		return
	}

	// 检查用户是否存在
	var user model.User
	if err := a.db.First(&user, "id = ?", userID).Error; err != nil {
		response.NewResponseHelper(c).Error(response.CodeNotFound, "user not found")
		return
	}

	if err := a.db.Model(&user).Update("status", req.Status).Error; err != nil {
		response.NewResponseHelper(c).Error(response.CodeDatabaseError, "failed to update user status")
		return
	}

	response.NewResponseHelper(c).Success(gin.H{"status": req.Status})
}

// adminDeleteUser 删除用户
func (a *AdminApp) adminDeleteUser(c *gin.Context) {
	userID := c.Param("id")

	// 检查用户是否存在
	var user model.User
	if err := a.db.First(&user, "id = ?", userID).Error; err != nil {
		response.NewResponseHelper(c).Error(response.CodeNotFound, "user not found")
		return
	}

	// 删除用户（实际应该标记为删除，而不是物理删除）
	if err := a.db.Delete(&user).Error; err != nil {
		response.NewResponseHelper(c).Error(response.CodeDatabaseError, "failed to delete user")
		return
	}

	response.NewResponseHelper(c).Success(gin.H{"message": "user deleted successfully"})
}

// adminBatchUpdateUserStatus 批量更新用户状态
func (a *AdminApp) adminBatchUpdateUserStatus(c *gin.Context) {
	var req model.BatchUserStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NewResponseHelper(c).Error(response.CodeInvalidParam, err.Error())
		return
	}

	// 去重用户ID
	uniqueUserIDs := make(map[string]bool)
	for _, id := range req.UserIDs {
		uniqueUserIDs[id] = true
	}
	req.UserIDs = make([]string, 0, len(uniqueUserIDs))
	for id := range uniqueUserIDs {
		req.UserIDs = append(req.UserIDs, id)
	}

	// 验证所有用户存在
	var count int64
	if err := a.db.Model(&model.User{}).Where("id IN ?", req.UserIDs).Count(&count).Error; err != nil {
		response.NewResponseHelper(c).Error(response.CodeDatabaseError, "failed to verify users")
		return
	}
	if int(count) != len(req.UserIDs) {
		response.NewResponseHelper(c).Error(response.CodeNotFound, "one or more users not found")
		return
	}

	// 批量更新
	result := a.db.Model(&model.User{}).
		Where("id IN ?", req.UserIDs).
		Update("status", req.Status)

	if result.Error != nil {
		response.NewResponseHelper(c).Error(response.CodeDatabaseError, "failed to update user status")
		return
	}

	response.NewResponseHelper(c).Success(gin.H{
		"updated_count": result.RowsAffected,
	})
}

// adminGetModels 获取模型列表
func (a *AdminApp) adminGetModels(c *gin.Context) {
	// NEW VERSION - This should be executed now
	var models []model.AdminModel

	// 先获取基础模型数据
	if err := a.db.Model(&model.Model{}).
		Select("id, name, model_id, status, description, created_at").
		Find(&models).Error; err != nil {
		response.NewResponseHelper(c).Error(response.CodeDatabaseError, "failed to load models")
		return
	}

	// 确保即使没有数据也返回空数组而不是null
	if models == nil {
		models = []model.AdminModel{}
	}

	// 获取总数
	var total int64
	a.db.Model(&model.Model{}).Count(&total)

	response.NewResponseHelper(c).Success(gin.H{
		"models": models,
		"total": total,
		"page":  1,
		"page_size": 20,
	})
}

// getModelsData 获取模型数据的辅助函数
func (a *AdminApp) getModelsData() ([]model.AdminModel, error) {
	var models []model.AdminModel

	// 尝试直接查询
	err := a.db.Model(&model.Model{}).
		Select("id, name, model_id, status, description, created_at").
		Find(&models).Error

	if err != nil {
		return nil, err
	}

	// 检查是否为nil
	if models == nil {
		models = []model.AdminModel{}
	}

	return models, nil
}

// adminCreateModel 创建模型
func (a *AdminApp) adminCreateModel(c *gin.Context) {
	var req model.CreateAdminModelRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NewResponseHelper(c).Error(response.CodeInvalidParam, "invalid request")
		return
	}

	// 检查模型ID是否已存在
	var existingModel model.Model
	if err := a.db.Where("model_id = ?", req.ModelID).First(&existingModel).Error; err == nil {
		response.NewResponseHelper(c).Error(response.CodeInvalidParam, "model_id already exists")
		return
	}

	// 创建模型
	modelData := model.Model{
		ID:        generateID(),
		Name:      req.Name,
		ModelID:   req.ModelID,
		Status:    req.Status,
	}

	if req.Description != "" {
		modelData.Description = req.Description
	}

	if err := a.db.Create(&modelData).Error; err != nil {
		response.NewResponseHelper(c).Error(response.CodeDatabaseError, "failed to create model")
		return
	}

	response.NewResponseHelper(c).Success(modelData)
}

// adminUpdateModel 更新模型
func (a *AdminApp) adminUpdateModel(c *gin.Context) {
	modelID := c.Param("id")

	var req model.UpdateAdminModelRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NewResponseHelper(c).Error(response.CodeInvalidParam, "invalid request")
		return
	}

	// 检查模型是否存在
	var modelData model.Model
	if err := a.db.First(&modelData, "id = ?", modelID).Error; err != nil {
		response.NewResponseHelper(c).Error(response.CodeNotFound, "model not found")
		return
	}

	// 更新模型 - 只更新非零值字段
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.ModelID != "" {
		updates["model_id"] = req.ModelID
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	// description 允许为空字符串，需要特殊处理
	updates["description"] = req.Description

	if len(updates) == 0 {
		response.NewResponseHelper(c).Error(response.CodeInvalidParam, "no fields to update")
		return
	}

	if err := a.db.Model(&modelData).Updates(updates).Error; err != nil {
		response.NewResponseHelper(c).Error(response.CodeDatabaseError, "failed to update model")
		return
	}

	// 重新获取更新后的模型数据
	a.db.First(&modelData, "id = ?", modelID)
	response.NewResponseHelper(c).Success(modelData)
}

// adminToggleModelStatus 切换模型状态
func (a *AdminApp) adminToggleModelStatus(c *gin.Context) {
	modelID := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NewResponseHelper(c).Error(response.CodeInvalidParam, "invalid request")
		return
	}

	// 检查模型是否存在
	var modelData model.Model
	if err := a.db.First(&modelData, "id = ?", modelID).Error; err != nil {
		response.NewResponseHelper(c).Error(response.CodeNotFound, "model not found")
		return
	}

	if err := a.db.Model(&modelData).Update("status", req.Status).Error; err != nil {
		response.NewResponseHelper(c).Error(response.CodeDatabaseError, "failed to update model status")
		return
	}

	response.NewResponseHelper(c).Success(gin.H{"status": req.Status})
}

// adminDeleteModel 删除模型
func (a *AdminApp) adminDeleteModel(c *gin.Context) {
	modelID := c.Param("id")

	// 检查模型是否存在
	var modelData model.Model
	if err := a.db.First(&modelData, "id = ?", modelID).Error; err != nil {
		response.NewResponseHelper(c).Error(response.CodeNotFound, "model not found")
		return
	}

	// 检查是否有关联的消息
	var messageCount int64
	a.db.Model(&model.Message{}).Where("model_id = ?", modelID).Count(&messageCount)
	if messageCount > 0 {
		response.NewResponseHelper(c).Error(response.CodeForbidden, "cannot delete model with associated messages")
		return
	}

	// 删除模型
	if err := a.db.Delete(&modelData).Error; err != nil {
		response.NewResponseHelper(c).Error(response.CodeDatabaseError, "failed to delete model")
		return
	}

	response.NewResponseHelper(c).Success(gin.H{"message": "model deleted successfully"})
}

// registerAdminRoutes 注册管理员路由
func (a *App) registerAdminRoutes() {
	// 创建管理员应用
	adminApp := NewAdminApp(a.adminAuthService, a.db)

	// 公开路由
	admin := a.Engine.Group("/api/admin")
	{
		admin.POST("/login", adminApp.adminLogin)
	}

	// 需要认证的路由
	adminGroup := admin.Use(CreateAdminMiddleware(a.adminAuthService))
	{
		adminGroup.GET("/me", adminApp.adminGetMe)

		// 用户管理
		adminGroup.GET("/users", adminApp.adminGetUsers)
		adminGroup.POST("/users", adminApp.adminCreateUser)
		adminGroup.PUT("/users/:id", adminApp.adminUpdateUser)
		adminGroup.PATCH("/users/:id/status", adminApp.adminToggleUserStatus)
		adminGroup.DELETE("/users/:id", adminApp.adminDeleteUser)

		// 模型管理
		adminGroup.GET("/models", adminApp.adminGetModels)
		adminGroup.POST("/models", adminApp.adminCreateModel)
		adminGroup.PUT("/models/:id", adminApp.adminUpdateModel)
		adminGroup.PATCH("/models/:id/status", adminApp.adminToggleModelStatus)
		adminGroup.DELETE("/models/:id", adminApp.adminDeleteModel)
	}

	// Serve admin index page (public)
	a.Engine.GET("/admin/login.html", func(c *gin.Context) {
		c.File("./admin/login.html")
	})

	// Admin index page requires authentication
	a.Engine.GET("/admin/index.html", func(c *gin.Context) {
		// Check if user is authenticated
		token := c.GetHeader("Authorization")
		if token == "" {
			// Not authenticated, redirect to login
			c.Redirect(http.StatusFound, "/admin/login.html")
			return
		}

		// Validate token
		_, err := a.adminAuthService.ValidateToken(token)
		if err != nil {
			// Invalid token, redirect to login
			c.Redirect(http.StatusFound, "/admin/login.html")
			return
		}

		// Serve the admin page
		c.File("./admin/index.html")
	})
}

// generateID 生成ID
func generateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return "admin_" + hex.EncodeToString(bytes)
}

// generateUserID 生成用户ID
func generateUserID() string {
	return uuid.New().String()
}