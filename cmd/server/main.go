package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"image-compressor/internal/api"
	"image-compressor/internal/config"

	"github.com/gin-gonic/gin"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	// バージョン情報の表示
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.Parse()

	if showVersion {
		fmt.Printf("Image Compressor v%s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Build Date: %s\n", date)
		return
	}

	// 設定を読み込み
	cfg := config.Load()

	// Ginの設定
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// CORS設定
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	})

	// 静的ファイルの提供
	r.Static("/static", "./web/static")
	r.LoadHTMLGlob("web/static/*.html")

	// ルート設定
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// バージョン情報API
	r.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version": version,
			"commit":  commit,
			"date":    date,
		})
	})

	// APIルート
	api.SetupRoutes(r)

	// サーバー起動
	log.Printf("Image Compressor v%s starting on port %s", version, cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
