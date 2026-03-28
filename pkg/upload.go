package pkg

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// UploadConfig holds upload configuration
type UploadConfig struct {
	MaxFileSize   int64    // in bytes
	AllowedTypes  []string // e.g., []string{"image/jpeg", "image/png", "image/jpg"}
	UploadDir     string   // e.g., "uploads/products"
}

// DefaultImageUploadConfig returns default config for image uploads
func DefaultImageUploadConfig() UploadConfig {
	return UploadConfig{
		MaxFileSize:  5 * 1024 * 1024, // 5MB
		AllowedTypes: []string{"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp"},
		UploadDir:    "uploads/products",
	}
}

// ValidateFile validates uploaded file
func (c *UploadConfig) ValidateFile(file *multipart.FileHeader) error {
	// Check file size
	if file.Size > c.MaxFileSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", c.MaxFileSize)
	}

	// Check file type
	contentType := file.Header.Get("Content-Type")
	isAllowed := false
	for _, allowedType := range c.AllowedTypes {
		if contentType == allowedType {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return fmt.Errorf("file type %s is not allowed. Allowed types: %v", contentType, c.AllowedTypes)
	}

	return nil
}

// SaveFile saves uploaded file and returns the file path
func (c *UploadConfig) SaveFile(file *multipart.FileHeader) (string, error) {
	// Validate file
	if err := c.ValidateFile(file); err != nil {
		return "", err
	}

	// Create upload directory if not exists
	if err := os.MkdirAll(c.UploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	filePath := filepath.Join(c.UploadDir, filename)

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	// Return relative path (for URL generation)
	return filePath, nil
}

// DeleteFile deletes a file from filesystem
func DeleteFile(filePath string) error {
	if filePath == "" {
		return nil
	}

	// Only delete if it's a local file (not URL)
	if strings.HasPrefix(filePath, "http://") || strings.HasPrefix(filePath, "https://") {
		return nil
	}

	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}

// GetFileURL converts file path to URL
func GetFileURL(filePath string, baseURL string) string {
	if filePath == "" {
		return ""
	}

	// If already a URL, return as is
	if strings.HasPrefix(filePath, "http://") || strings.HasPrefix(filePath, "https://") {
		return filePath
	}

	// Convert to URL
	return fmt.Sprintf("%s/%s", strings.TrimSuffix(baseURL, "/"), strings.TrimPrefix(filePath, "/"))
}
