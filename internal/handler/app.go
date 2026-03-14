package handler

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/mfb27/luban/internal/config"
	"github.com/mfb27/luban/internal/middleware"
	"github.com/mfb27/luban/internal/model"
	"github.com/mfb27/luban/internal/storage"
	"github.com/mfb27/luban/internal/zhipu"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	Engine *gin.Engine

	cfg    *config.Config
	log    *zap.Logger
	db     *gorm.DB
	redis  *redis.Client
	minio  *storage.MinIO
	zhipu  *zhipu.Client

	minioBucket      string
	minioPublicBase  string
}

type AppDeps struct {
	Cfg   *config.Config
	Log   *zap.Logger
	DB    *gorm.DB
	Redis *redis.Client
	MinIO *storage.MinIO
}

func NewApp(deps AppDeps) (*App, error) {
	if err := deps.DB.AutoMigrate(&model.Session{}, &model.Message{}, &model.User{}, &model.Model{}, &model.Attachment{}); err != nil {
		return nil, err
	}

	// Migrate existing sessions to add user_id
	if err := migrateSessions(deps.DB); err != nil {
		deps.Log.Error("Failed to migrate sessions", zap.Error(err))
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.WithRequestID(deps.Log))
	engine.Use(gin.Logger())

	// Create ZhipuAI client
	zhipuClient := zhipu.NewClient(&zhipu.Config{
		APIKey:  deps.Cfg.ZhipuAI.APIKey,
		BaseURL: deps.Cfg.ZhipuAI.BaseURL,
	})

	app := &App{
		Engine:          engine,
		cfg:             deps.Cfg,
		log:             deps.Log,
		db:              deps.DB,
		redis:           deps.Redis,
		minio:           deps.MinIO,
		zhipu:           zhipuClient,
		minioBucket:     deps.MinIO.Bucket,
		minioPublicBase: deps.MinIO.PublicBaseURL,
	}

	app.registerRoutes()
	app.seedIfNeeded()

	return app, nil
}

func (a *App) registerRoutes() {
	a.Engine.GET("/api/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	// Public routes
	api := a.Engine.Group("/api")
	{
		// Auth routes
		api.POST("/auth/register", a.register)
		api.POST("/auth/login", a.login)

		// Public routes (no authentication required)
		api.GET("/models", a.getModels)

		// Protected routes (require authentication)
		auth := api.Group("/user")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.GET("/me", a.getUser)
		}

		// Chat route (public)
		api.POST("/chat", a.postChat)

		// Protected routes
		protected := api.Group("/", middleware.AuthMiddleware())
		{
			protected.GET("/sessions", a.listSessions)
			protected.POST("/sessions", a.createSession)
			protected.GET("/sessions/:id/messages", a.listMessages)
			protected.POST("/upload", a.upload)
		}
	}

	// static
	staticDir := a.cfg.Web.StaticDir
	if staticDir == "" {
		staticDir = "./web"
	}
	a.Engine.Static("/assets", filepath.Join(staticDir, "assets"))
	a.Engine.StaticFile("/", filepath.Join(staticDir, "index.html"))
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

