package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"image-compressor/internal/config"
	"image-compressor/internal/service"

	"github.com/gin-gonic/gin"
)

// Handler APIハンドラー
type Handler struct {
	imageService *service.ImageService
	config       *config.Config
}

// NewHandler 新しいハンドラーを作成
func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		imageService: service.NewImageService(cfg),
		config:       cfg,
	}
}

// CompressImage 画像圧縮API
func (h *Handler) CompressImage(c *gin.Context) {
	// ファイルの取得
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "No image file provided",
		})
		return
	}
	defer file.Close()

	// ファイルの検証
	if err := h.imageService.ValidateFile(file, header); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// 圧縮オプションの取得
	quality := h.config.DefaultQuality
	if qualityStr := c.PostForm("quality"); qualityStr != "" {
		if q, err := strconv.Atoi(qualityStr); err == nil && q >= 0 && q <= 100 {
			quality = q
		}
	}

	width := h.config.DefaultWidth
	if widthStr := c.PostForm("width"); widthStr != "" {
		if w, err := strconv.Atoi(widthStr); err == nil && w > 0 {
			width = w
		}
	}

	height := h.config.DefaultHeight
	if heightStr := c.PostForm("height"); heightStr != "" {
		if h, err := strconv.Atoi(heightStr); err == nil && h > 0 {
			height = h
		}
	}

	// 画像の圧縮
	result, err := h.imageService.CompressImage(file, quality, width, height)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Failed to compress image: %v", err),
		})
		return
	}

	// 成功レスポンス
	c.JSON(http.StatusOK, gin.H{
		"success":           true,
		"message":           "Image compressed successfully",
		"original_size":     result.OriginalSize,
		"compressed_size":   result.CompressedSize,
		"compression_ratio": result.CompressionRatio,
		"processing_time":   result.ProcessingTime,
		"output_file":       result.OutputFileName,
		"download_url":      fmt.Sprintf("/api/download/%s", result.OutputFileName),
	})
}

// DownloadFile ファイルダウンロードAPI
func (h *Handler) DownloadFile(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "No filename provided",
		})
		return
	}

	// セキュリティチェック（パストラバーサル攻撃を防ぐ）
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid filename",
		})
		return
	}

	filePath := fmt.Sprintf("%s/%s", h.config.DownloadDir, filename)
	c.File(filePath)
}

// SetupRoutes ルートの設定
func SetupRoutes(r *gin.Engine) {
	cfg := config.Load()
	handler := NewHandler(cfg)

	// APIルート
	api := r.Group("/api")
	{
		api.POST("/compress", handler.CompressImage)
		api.GET("/download/:filename", handler.DownloadFile)
	}
}
