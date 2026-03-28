package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	Log     LogConfig     `mapstructure:"log"`
	MySQL   MySQLConfig   `mapstructure:"mysql"`
	Redis   RedisConfig   `mapstructure:"redis"`
	MinIO   MinIOConfig   `mapstructure:"minio"`
	ZhipuAI ZhipuAIConfig `mapstructure:"zhipuai"`
	Admin   AdminConfig   `mapstructure:"admin"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

type MySQLConfig struct {
	DSN string `mapstructure:"dsn"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type MinIOConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	Bucket          string `mapstructure:"bucket"`
	UseSSL          bool   `mapstructure:"use_ssl"`
	PublicBaseURL   string `mapstructure:"public_base_url"`
}

type ZhipuAIConfig struct {
	APIKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
}

type AdminConfig struct {
	SecretKey    string             `mapstructure:"secret_key"`
	InitialAdmin InitialAdminConfig `mapstructure:"initial_admin"`
}

type InitialAdminConfig struct {
	Name     string `mapstructure:"name"`
	Email    string `mapstructure:"email"`
	Password string `mapstructure:"password"`
}

func Load() (*Config, error) {
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("web.static_dir", "./web")
	viper.SetDefault("log.level", "info")
	viper.SetDefault("redis.addr", "127.0.0.1:6379")
	viper.SetDefault("redis.db", 0)

	// Allow env var override like LUBAN_MYSQL_DSN, LUBAN_MINIO_ACCESS_KEY_ID...
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if cfg.MySQL.DSN == "" {
		return nil, fmt.Errorf("mysql.dsn is required (or set LUBAN_MYSQL_DSN)")
	}
	if cfg.MinIO.Endpoint == "" {
		return nil, fmt.Errorf("minio.endpoint is required (or set LUBAN_MINIO_ENDPOINT)")
	}
	if cfg.MinIO.AccessKeyID == "" || cfg.MinIO.SecretAccessKey == "" {
		return nil, fmt.Errorf("minio access keys are required (minio.access_key_id/minio.secret_access_key)")
	}
	if cfg.MinIO.Bucket == "" {
		return nil, fmt.Errorf("minio.bucket is required")
	}
	if cfg.MinIO.PublicBaseURL == "" {
		// default to endpoint for simple dev usage
		scheme := "http"
		if cfg.MinIO.UseSSL {
			scheme = "https"
		}
		cfg.MinIO.PublicBaseURL = fmt.Sprintf("%s://%s", scheme, cfg.MinIO.Endpoint)
	}

	if cfg.ZhipuAI.APIKey == "" {
		return nil, fmt.Errorf("zhipuai.api_key is required (or set LUBAN_ZHIPUAI_API_KEY)")
	}

	return &cfg, nil
}
