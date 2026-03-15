package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	NewRelic NewRelicConfig
	OpenAI   OpenAIConfig
	SMTP     SMTPConfig
	Storage  StorageConfig
}

type ServerConfig struct {
	Port         string
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Environment  string
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type JWTConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

type NewRelicConfig struct {
	LicenseKey string
	AppName    string
	Enabled    bool
}

type OpenAIConfig struct {
	APIKey string
	Model  string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	Enabled  bool
}

type StorageConfig struct {
	Type         string // "local" or "s3"
	LocalPath    string
	S3Bucket     string
	S3Region     string
	S3AccessKey  string
	S3SecretKey  string
	MaxFileSize  int64 // in bytes
	AllowedTypes []string
}

var AppConfig *Config

func Load() error {
	// Load .env file if it exists (not required in production)
	_ = godotenv.Load()

	AppConfig = &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 15*time.Second),
			Environment:  getEnv("ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "123456"),
			DBName:          getEnv("DB_NAME", "hrms"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		JWT: JWTConfig{
			AccessTokenSecret:  getEnv("JWT_ACCESS_SECRET", "your-access-secret-key-change-in-production"),
			RefreshTokenSecret: getEnv("JWT_REFRESH_SECRET", "your-refresh-secret-key-change-in-production"),
			AccessTokenExpiry:  getDurationEnv("JWT_ACCESS_EXPIRY", 15*time.Minute),
			RefreshTokenExpiry: getDurationEnv("JWT_REFRESH_EXPIRY", 7*24*time.Hour),
		},
		NewRelic: NewRelicConfig{
			LicenseKey: getEnv("NEW_RELIC_LICENSE_KEY", ""),
			AppName:    getEnv("NEW_RELIC_APP_NAME", "HRMS"),
			Enabled:    getBoolEnv("NEW_RELIC_ENABLED", false),
		},
		OpenAI: OpenAIConfig{
			APIKey: getEnv("OPENAI_API_KEY", ""),
			Model:  getEnv("OPENAI_MODEL", "gpt-3.5-turbo"),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", ""),
			Port:     getIntEnv("SMTP_PORT", 587),
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", "noreply@hrms.com"),
			Enabled:  getBoolEnv("SMTP_ENABLED", false),
		},
		Storage: StorageConfig{
			Type:         getEnv("STORAGE_TYPE", "local"),
			LocalPath:    getEnv("STORAGE_LOCAL_PATH", "./uploads"),
			S3Bucket:     getEnv("S3_BUCKET", ""),
			S3Region:     getEnv("S3_REGION", "us-east-1"),
			S3AccessKey:  getEnv("S3_ACCESS_KEY", ""),
			S3SecretKey:  getEnv("S3_SECRET_KEY", ""),
			MaxFileSize:  getInt64Env("STORAGE_MAX_FILE_SIZE", 10*1024*1024), // 10MB
			AllowedTypes: []string{".pdf", ".doc", ".docx", ".jpg", ".jpeg", ".png"},
		},
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getInt64Env(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}
