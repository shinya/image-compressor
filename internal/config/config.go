package config

import (
	"os"
)

type Config struct {
	Port           string
	DownloadDir    string
	MaxFileSize    int64
	AllowedFormats []string
	DefaultQuality int
	DefaultWidth   int
	DefaultHeight  int
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		DownloadDir:    getEnv("DOWNLOAD_DIR", "./downloads"),
		MaxFileSize:    10 * 1024 * 1024, // 10MB
		AllowedFormats: []string{".png", ".jpg", ".jpeg", ".gif"},
		DefaultQuality: 80,
		DefaultWidth:   1920,
		DefaultHeight:  1080,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
