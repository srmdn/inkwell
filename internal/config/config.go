package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	ContentDir  string
	JWTSecret   string
	AdminEmail  string
	AdminPasswd string

	// Theme
	ThemeDir      string
	ThemeBuildCmd string
	ThemeService  string // optional: systemd service to restart after build
}

func Load(envFile string) (*Config, error) {
	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			return nil, fmt.Errorf("loading env file: %w", err)
		}
	}

	port := getEnv("PORT", "8090")
	if _, err := strconv.Atoi(port); err != nil {
		return nil, fmt.Errorf("PORT must be a number, got %q", port)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return &Config{
		Port:          port,
		DatabaseURL:   getEnv("DATABASE_URL", "data/inkwell.db"),
		ContentDir:    getEnv("CONTENT_DIR", "content/blog"),
		JWTSecret:     jwtSecret,
		AdminEmail:    os.Getenv("ADMIN_EMAIL"),
		AdminPasswd:   os.Getenv("ADMIN_PASSWORD"),
		ThemeDir:      getEnv("THEME_DIR", "theme"),
		ThemeBuildCmd: getEnv("THEME_BUILD_CMD", "npm run build"),
		ThemeService:  os.Getenv("THEME_SERVICE"),
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
