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

// CompressionResult 圧縮結果の詳細情報
type CompressionResult struct {
	OriginalSize     int64   `json:"original_size"`
	CompressedSize   int64   `json:"compressed_size"`
	CompressionRatio float64 `json:"compression_ratio"`
	ProcessingTime   float64 `json:"processing_time"`
	OutputFileName   string  `json:"output_file_name"`
}

// ImageService 画像処理サービス
type ImageService struct {
	config *config.Config
}

// NewImageService 新しいImageServiceを作成
func NewImageService(cfg *config.Config) *ImageService {
	return &ImageService{
		config: cfg,
	}
}

// CompressImage 画像をWebP形式に圧縮
func (s *ImageService) CompressImage(file multipart.File, quality int, width, height int) (*CompressionResult, error) {
	startTime := time.Now()

	log.Println("Starting image compression...")

	// 画像をデコード
	img, _, err := image.Decode(file)
	if err != nil {
		log.Printf("Failed to decode image: %v", err)
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	log.Println("Image decoded successfully")

	// リサイズ処理
	if width > 0 || height > 0 {
		log.Printf("Resizing image to %dx%d", width, height)
		img = imaging.Fit(img, width, height, imaging.Lanczos)
		log.Println("Image resized successfully")
	}

	// 出力ファイル名を生成
	outputFileName := fmt.Sprintf("compressed_%d.webp", time.Now().Unix())
	outputPath := filepath.Join(s.config.DownloadDir, outputFileName)

	// 出力ファイルを作成
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Printf("Failed to create output file: %v", err)
		return nil, fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// WebPエンコーダーの設定
	options := &webp.Options{
		Lossless: false,
		Quality:  float32(quality),
	}

	// WebP形式でエンコード
	log.Println("Encoding to WebP format...")
	if err := webp.Encode(outputFile, img, options); err != nil {
		log.Printf("Failed to encode WebP: %v", err)
		return nil, fmt.Errorf("failed to encode WebP: %w", err)
	}
	log.Println("WebP encoding completed successfully")

	// ファイルサイズを取得
	outputFile.Seek(0, 0)
	compressedSize, err := outputFile.Seek(0, 2)
	if err != nil {
		log.Printf("Failed to get compressed file size: %v", err)
		return nil, fmt.Errorf("failed to get compressed file size: %w", err)
	}

	// 元のファイルサイズを取得
	file.Seek(0, 0)
	originalSize, err := file.Seek(0, 2)
	if err != nil {
		log.Printf("Failed to get original file size: %v", err)
		return nil, fmt.Errorf("failed to get original file size: %w", err)
	}

	// 圧縮率を計算
	compressionRatio := float64(originalSize-compressedSize) / float64(originalSize) * 100

	processingTime := time.Since(startTime).Seconds()

	log.Printf("Compression completed: %d bytes -> %d bytes (%.2f%% reduction, %.2fs)",
		originalSize, compressedSize, compressionRatio, processingTime)

	return &CompressionResult{
		OriginalSize:     originalSize,
		CompressedSize:   compressedSize,
		CompressionRatio: compressionRatio,
		ProcessingTime:   processingTime,
		OutputFileName:   outputFileName,
	}, nil
}

// ValidateFile ファイルの検証
func (s *ImageService) ValidateFile(file multipart.File, header *multipart.FileHeader) error {
	// ファイルサイズの確認
	if header.Size > s.config.MaxFileSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", s.config.MaxFileSize)
	}

	// ファイル形式の確認
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !s.isAllowedFormat(ext) {
		return fmt.Errorf("unsupported file format. Allowed formats: %s", strings.Join(s.config.AllowedFormats, ", "))
	}

	return nil
}

// GetFileSize ファイルサイズを取得
func (s *ImageService) GetFileSize(file multipart.File) (int64, error) {
	file.Seek(0, 0)
	size, err := file.Seek(0, 2)
	if err != nil {
		return 0, err
	}
	file.Seek(0, 0)
	return size, nil
}

// isAllowedFormat 許可された形式かどうかを確認
func (s *ImageService) isAllowedFormat(ext string) bool {
	for _, format := range s.config.AllowedFormats {
		if ext == format {
			return true
		}
	}
	return false
}
