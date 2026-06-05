package util

import (
	"path/filepath"
	"strings"
)

// Whitelist of allowed image MIME types
var allowedImageTypes = map[string]bool{
	"image/jpeg":    true,
	"image/jpg":     true,
	"image/png":     true,
	"image/gif":     true,
	"image/webp":    true,
	"image/bmp":     true,
	"image/tiff":    true,
	"image/svg+xml": true,
}

// Whitelist of allowed image file extensions
var allowedImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
	".bmp":  true,
	".tiff": true,
	".tif":  true,
	".svg":  true,
}

// ValidateImageType checks if the file is a valid image type
func ValidateImageType(filename, contentType string) bool {
	// Check MIME type
	if contentType != "" && allowedImageTypes[contentType] {
		return true
	}

	// Check file extension as fallback
	ext := strings.ToLower(filepath.Ext(filename))
	return allowedImageExtensions[ext]
}
