package middleware

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

const (
	MaxFileSize = 10 * 1024 * 1024 // 10MB
)

// ValidateFileUpload checks if uploaded file is valid
func ValidateFileUpload(file multipart.File, header *multipart.FileHeader) error {
	// Check file size
	if header.Size > MaxFileSize {
		return fmt.Errorf("file size %d exceeds maximum allowed size of %d bytes", 
			header.Size, MaxFileSize)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".pdf" {
		return fmt.Errorf("only PDF files are allowed, got %s", ext)
	}

	// Verify it's actually a PDF by checking magic bytes
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Reset file pointer to beginning
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to reset file pointer: %w", err)
	}

	// Check for PDF signature
	contentType := http.DetectContentType(buffer)
	if contentType != "application/pdf" {
		// Double-check with PDF magic bytes
		if !isPDF(buffer) {
			return fmt.Errorf("file is not a valid PDF (detected type: %s)", contentType)
		}
	}

	return nil
}

// isPDF checks for PDF magic bytes
func isPDF(data []byte) bool {
	// PDF files start with %PDF-
	return len(data) >= 5 && string(data[0:5]) == "%PDF-"
}

// SanitizeFilename removes dangerous characters from filename
func SanitizeFilename(filename string) string {
	// Remove path separators and other dangerous characters
	filename = strings.ReplaceAll(filename, "/", "")
	filename = strings.ReplaceAll(filename, "\\", "")
	filename = strings.ReplaceAll(filename, "..", "")
	filename = strings.TrimSpace(filename)
	
	// Limit length
	if len(filename) > 255 {
		ext := filepath.Ext(filename)
		filename = filename[:255-len(ext)] + ext
	}
	
	return filename
}