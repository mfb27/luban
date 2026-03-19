package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mfb27/luban/internal/cache"
	"github.com/mfb27/luban/internal/config"
	"github.com/mfb27/luban/internal/handler"
	"github.com/mfb27/luban/internal/logger"
	"github.com/mfb27/luban/internal/storage"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// NewGormLogger 创建一个使用zap的GORM日志记录器
func NewGormLogger(log *zap.Logger) gormlogger.Interface {
	return gormlogger.New(
		&gormLoggerWriter{log: log},
		gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  gormlogger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
}

// gormLoggerWriter 实现GORM日志写入器，使用zap
type gormLoggerWriter struct {
	log *zap.Logger
}

func (l *gormLoggerWriter) Printf(format string, args ...interface{}) {
	l.log.Info(fmt.Sprintf(format, args...))
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start luban HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		log, err := logger.New(cfg.Log)
		if err != nil {
			return err
		}
		defer func() { _ = log.Sync() }()

		db, err := gorm.Open(mysql.Open(cfg.MySQL.DSN), &gorm.Config{
			Logger: NewGormLogger(log),
		})
		if err != nil {
			return fmt.Errorf("open mysql: %w", err)
		}

		rdb, err := cache.New(cfg.Redis)
		if err != nil {
			return fmt.Errorf("init redis: %w", err)
		}
		defer func() { _ = rdb.Close() }()

		minioClient, err := storage.NewMinIO(cfg.MinIO)
		if err != nil {
			return fmt.Errorf("init minio: %w", err)
		}

		app, err := handler.NewApp(handler.AppDeps{
			Cfg:   cfg,
			Log:   log,
			DB:    db,
			Redis: rdb,
			MinIO: minioClient,
		})
		if err != nil {
			return err
		}

		srv := &http.Server{
			Addr:              fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
			Handler:           app.Engine,
			ReadHeaderTimeout: 10 * time.Second,
		}

		go func() {
			log.Info("http server starting", zap.String("addr", srv.Addr))
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal("http server failed", zap.Error(err))
			}
		}()

		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch

		log.Info("shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(ctx)
	},
}

