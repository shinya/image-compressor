package service

import (
	"fmt"
	"image"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"image-compressor/internal/config"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

type ImageService struct {
	config *config.Config
}

type CompressionOptions struct {
	Quality int `json:"quality"`
	Width   int `json:"width"`
	Height  int `json:"height"`
}

type CompressionResult struct {
	OriginalSize     int64   `json:"original_size"`
	CompressedSize   int64   `json:"compressed_size"`
	CompressionRatio float64 `json:"compression_ratio"`
	ProcessingTime   float64 `json:"processing_time"`
}

func NewImageService(cfg *config.Config) *ImageService {
	return &ImageService{
		config: cfg,
	}
}

func (s *ImageService) CompressImage(file *multipart.FileHeader, options *CompressionOptions) (*CompressionResult, error) {
	startTime := time.Now()

	log.Printf("Starting compression for file: %s (size: %d bytes)", file.Filename, file.Size)

	// ファイルを開く
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// 画像をデコード
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// オプションが指定されていない場合はデフォルト値を使用
	if options == nil {
		options = &CompressionOptions{
			Quality: s.config.DefaultQuality,
			Width:   s.config.DefaultWidth,
			Height:  s.config.DefaultHeight,
		}
	}

	log.Printf("Compression options - Quality: %d, Width: %d, Height: %d",
		options.Quality, options.Width, options.Height)

	// リサイズ処理
	if options.Width > 0 || options.Height > 0 {
		img = imaging.Fit(img, options.Width, options.Height, imaging.Lanczos)
		log.Printf("Image resized to fit: %dx%d", options.Width, options.Height)
	}

	// 出力ファイル名を生成
	originalName := filepath.Base(file.Filename)
	ext := filepath.Ext(originalName)
	baseName := strings.TrimSuffix(originalName, ext)
	outputFileName := fmt.Sprintf("%s_compressed.webp", baseName)
	outputPath := filepath.Join(s.config.DownloadDir, outputFileName)

	// 出力ディレクトリを作成
	if err := os.MkdirAll(s.config.DownloadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create download directory: %w", err)
	}

	// WebPファイルを作成
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// WebPとしてエンコード
	webpOptions := &webp.Options{
		Lossless: false,
		Quality:  float32(options.Quality),
	}

	if err := webp.Encode(outputFile, img, webpOptions); err != nil {
		return nil, fmt.Errorf("failed to encode WebP: %w", err)
	}

	// 圧縮後のファイルサイズを取得
	compressedSize, err := s.GetFileSize(outputPath)
	if err != nil {
		log.Printf("Warning: Could not get compressed file size: %v", err)
		compressedSize = 0
	}

	processingTime := time.Since(startTime).Seconds()
	compressionRatio := float64(compressedSize) / float64(file.Size) * 100

	result := &CompressionResult{
		OriginalSize:     file.Size,
		CompressedSize:   compressedSize,
		CompressionRatio: compressionRatio,
		ProcessingTime:   processingTime,
	}

	log.Printf("Compression completed - Original: %d bytes, Compressed: %d bytes, Ratio: %.2f%%, Time: %.2fs",
		result.OriginalSize, result.CompressedSize, result.CompressionRatio, result.ProcessingTime)

	return result, nil
}

func (s *ImageService) ValidateFile(file *multipart.FileHeader) error {
	// ファイルサイズチェック
	if file.Size > s.config.MaxFileSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", s.config.MaxFileSize)
	}

	// ファイル形式チェック
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := false
	for _, format := range s.config.AllowedFormats {
		if ext == format {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("file format not allowed. Allowed formats: %v", s.config.AllowedFormats)
	}

	return nil
}

func (s *ImageService) GetFileSize(filePath string) (int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}
