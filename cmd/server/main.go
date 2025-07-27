package main

import (
	"log"
	"net/http"

	"image-compressor/internal/api"
	"image-compressor/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {
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

	// APIルート
	api.SetupRoutes(r)

	// サーバー起動
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
