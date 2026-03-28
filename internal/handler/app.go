package handler

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	"github.com/mfb27/luban/internal/admin"
	"github.com/mfb27/luban/internal/config"
	"github.com/mfb27/luban/internal/middleware"
	"github.com/mfb27/luban/internal/model"
	"github.com/mfb27/luban/internal/response"
	"github.com/mfb27/luban/internal/storage"
	"github.com/mfb27/luban/internal/zhipu"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type App struct {
	Engine           *gin.Engine
	adminAuthService *admin.AdminAuthService

	cfg   *config.Config
	log   *zap.Logger
	db    *gorm.DB
	redis *redis.Client
	minio *storage.MinIO
	zhipu *zhipu.Client

	minioBucket     string
	minioPublicBase string
}

type AppDeps struct {
	Cfg   *config.Config
	Log   *zap.Logger
	DB    *gorm.DB
	Redis *redis.Client
	MinIO *storage.MinIO
}

func NewApp(deps AppDeps) (*App, error) {
	if err := deps.DB.AutoMigrate(&model.Session{}, &model.Message{}, &model.User{}, &model.Model{}, &model.Attachment{}, &model.Admin{}); err != nil {
		return nil, err
	}

	// Migrate existing sessions to add user_id
	if err := migrateSessions(deps.DB); err != nil {
		deps.Log.Error("Failed to migrate sessions", zap.Error(err))
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.WithRequestID(deps.Log))
	engine.Use(middleware.RequestLogger(deps.Log))

	// Use simple CORS middleware for development
	engine.Use(middleware.SimpleCORSMiddleware())

	// Create ZhipuAI client
	zhipuClient := zhipu.NewClient(&zhipu.Config{
		APIKey:  deps.Cfg.ZhipuAI.APIKey,
		BaseURL: deps.Cfg.ZhipuAI.BaseURL,
	})

	// Create admin auth service
	adminAuthService := admin.NewAdminAuthService(deps.DB, deps.Cfg.Admin.SecretKey)

	app := &App{
		Engine:           engine,
		adminAuthService: adminAuthService,
		cfg:              deps.Cfg,
		log:              deps.Log,
		db:               deps.DB,
		redis:            deps.Redis,
		minio:            deps.MinIO,
		zhipu:            zhipuClient,
		minioBucket:      deps.MinIO.Bucket,
		minioPublicBase:  deps.MinIO.PublicBaseURL,
	}

	app.registerRoutes()
	app.seedAdminIfNeeded()

	return app, nil
}

func (a *App) registerRoutes() {
	// Health check endpoint
	a.Engine.GET("/api/health", func(c *gin.Context) {
		response.NewResponseHelper(c).Success(gin.H{"ok": true})
	})

	// Public routes
	api := a.Engine.Group("/api")
	{
		// Auth routes
		api.POST("/auth/register", a.register)
		api.POST("/auth/login", a.login)

		// Public routes (no authentication required)
		api.GET("/models", a.getModels)

		// Chat route (public)
		api.POST("/chat", a.postChat)

		// Protected routes (require authentication)
		auth := api.Group("/user")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.GET("/me", a.getUser)
		}

		protected := api.Group("/", middleware.AuthMiddleware())
		{
			protected.GET("/sessions", a.listSessions)
			protected.POST("/sessions", a.createSession)
			protected.GET("/sessions/:id/messages", a.listMessages)
			protected.POST("/upload", a.upload)
		}
	}

	// Register admin routes
	a.registerAdminRoutes()

	// // static
	// staticDir := a.cfg.Web.StaticDir
	// if staticDir == "" {
	// 	staticDir = "./frontend"
	// }

	// // Serve main frontend
	// a.Engine.StaticFile("/", "./frontend/index.html")
	// a.Engine.Static("/assets", "./frontend/assets")
}

// migrateSessions adds user_id to existing sessions
func migrateSessions(db *gorm.DB) error {
	// Check if user_id column exists
	var hasUserID bool
	db.Raw("SELECT COUNT(*) FROM information_schema.columns WHERE table_name = ? AND column_name = ?", "sessions", "user_id").Scan(&hasUserID)

	if !hasUserID {
		// Add user_id column
		if err := db.Exec("ALTER TABLE sessions ADD COLUMN user_id VARCHAR(36)").Error; err != nil {
			return err
		}

		// Get all users
		var users []model.User
		if err := db.Find(&users).Error; err != nil {
			return err
		}

		// Keep sessions without user_id for anonymous users
	}

	return nil
}

// seedAdminIfNeeded 如果没有管理员则创建初始管理员
func (a *App) seedAdminIfNeeded() {
	var count int64
	if err := a.db.Model(&model.Admin{}).Count(&count).Error; err != nil {
		a.log.Error("Failed to check admin count", zap.Error(err))
		return
	}

	if count == 0 {
		// 创建初始管理员
		admin := model.Admin{
			ID:     generateAdminID(),
			Name:   a.cfg.Admin.InitialAdmin.Name,
			Email:  a.cfg.Admin.InitialAdmin.Email,
			Status: "active",
		}

		// 加密密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(a.cfg.Admin.InitialAdmin.Password), bcrypt.DefaultCost)
		if err != nil {
			a.log.Error("Failed to hash admin password", zap.Error(err))
			return
		}
		admin.Password = string(hashedPassword)

		if err := a.db.Create(&admin).Error; err != nil {
			a.log.Error("Failed to create initial admin", zap.Error(err))
		} else {
			a.log.Info("Created initial admin user", zap.String("email", admin.Email))
		}
	}
}

// generateAdminID 生成管理员ID
func generateAdminID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return "admin_" + hex.EncodeToString(bytes)
}
