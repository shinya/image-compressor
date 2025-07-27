package api

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"image-compressor/internal/config"
	"image-compressor/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	imageService *service.ImageService
	config       *config.Config
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		imageService: service.NewImageService(cfg),
		config:       cfg,
	}
}

func (h *Handler) CompressImage(c *gin.Context) {
	// ファイルを取得
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "No image file provided",
		})
		return
	}

	// ファイルの検証
	if err := h.imageService.ValidateFile(file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// 圧縮オプションを取得
	options := &service.CompressionOptions{
		Quality: h.config.DefaultQuality,
		Width:   h.config.DefaultWidth,
		Height:  h.config.DefaultHeight,
	}

	if qualityStr := c.PostForm("quality"); qualityStr != "" {
		if quality, err := strconv.Atoi(qualityStr); err == nil && quality >= 1 && quality <= 100 {
			options.Quality = quality
		}
	}

	if widthStr := c.PostForm("width"); widthStr != "" {
		if width, err := strconv.Atoi(widthStr); err == nil && width > 0 {
			options.Width = width
		}
	}

	if heightStr := c.PostForm("height"); heightStr != "" {
		if height, err := strconv.Atoi(heightStr); err == nil && height > 0 {
			options.Height = height
		}
	}

	// 画像を圧縮
	result, err := h.imageService.CompressImage(file, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to compress image: " + err.Error(),
		})
		return
	}

	// 出力ファイル名を生成
	originalName := filepath.Base(file.Filename)
	ext := filepath.Ext(originalName)
	baseName := strings.TrimSuffix(originalName, ext)
	outputFileName := baseName + "_compressed.webp"

	// レスポンスを返す
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"original_filename":   file.Filename,
			"compressed_filename": outputFileName,
			"original_size":       result.OriginalSize,
			"compressed_size":     result.CompressedSize,
			"compression_ratio":   result.CompressionRatio,
			"processing_time":     result.ProcessingTime,
			"download_url":        "/api/download/" + outputFileName,
		},
	})
}

func (h *Handler) DownloadFile(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "No filename provided",
		})
		return
	}

	filePath := filepath.Join(h.config.DownloadDir, filename)

	// ファイルの存在確認
	if _, err := h.imageService.GetFileSize(filePath); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "File not found",
		})
		return
	}

	// ファイルをダウンロード
	c.File(filePath)
}

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
