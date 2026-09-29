package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Env         string
	Port        string
	DatabaseURL string
	CORSOrigin  string
	MaxUploadMB int

	Auth struct {
		JWTAccessSecret    string
		JWTRefreshSecret   string
		AccessTokenMinutes int
		RefreshTokenDays   int
	}

	Email struct {
		Host     string
		Port     int
		Username string
		Password string
		From     string
	}

	Cloudinary struct {
		CloudName string
		APIKey    string
		APISecret string
	}
}

func MustLoad() (Config, error) {
	_ = godotenv.Load()
	var cfg Config

	cfg.Env = getEnv("APP_ENV", "development")
	cfg.Port = getEnv("PORT", "5000")
	cfg.DatabaseURL = mustEnv("DATABASE_URL")
	cfg.CORSOrigin = getEnv("CORS_ORIGIN", "http://localhost:3000")
	cfg.MaxUploadMB = mustInt("MAX_UPLOAD_MB", 5)

	cfg.Auth.JWTAccessSecret = mustEnv("JWT_ACCESS_SECRET")
	cfg.Auth.JWTRefreshSecret = mustEnv("JWT_REFRESH_SECRET")
	cfg.Auth.AccessTokenMinutes = mustInt("ACCESS_TOKEN_MINUTES", 15)
	cfg.Auth.RefreshTokenDays = mustInt("REFRESH_TOKEN_DAYS", 7)

	cfg.Email.Host = getEnv("SMTP_HOST", "smtp.gmail.com")
	cfg.Email.Port = mustInt("SMTP_PORT", 587)
	cfg.Email.Username = mustEnv("SMTP_USERNAME")
	cfg.Email.Password = mustEnv("SMTP_PASSWORD")
	cfg.Email.From = mustEnv("SMTP_FROM")

	cfg.Cloudinary.CloudName = mustEnv("CLOUDINARY_CLOUD_NAME")
	cfg.Cloudinary.APIKey = mustEnv("CLOUDINARY_API_KEY")
	cfg.Cloudinary.APISecret = mustEnv("CLOUDINARY_API_SECRET")

	return cfg, nil
}

func mustEnv(key string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		panic(fmt.Sprintf("environment variable %s is empty", key))
	}
	return value
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func mustInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("environment variable %s is invalid: %s", key, err))
	}
	return v
}

// func mustBool(key string, fallback bool) bool {
// 	value := strings.TrimSpace(os.Getenv(key))
// 	if value == "" {
// 		return fallback
// 	}
// 	v, err := strconv.ParseBool(value)
// 	if err != nil {
// 		panic(fmt.Sprintf("environment variable %s is invalid: %s", key, err))
// 	}
// 	return v
// }

// func mustDuration(key, fallback string) time.Duration {
// 	value := getEnv(key, fallback)

// 	d, err := time.ParseDuration(value)
// 	if err != nil {
// 		panic(fmt.Sprintf("environment variable %s is invalid: %s", key, err))
// 	}
// 	return d
// }
