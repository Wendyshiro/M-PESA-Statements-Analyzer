package utils

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// ExtractTextFromPDF extracts text content from a PDF using Tesseract OCR
func ExtractTextFromPDF(filePath string) (string, error) {
	// Create a temporary file for the output
	tempDir := filepath.Dir(filePath)
	tempOutput := filepath.Join(tempDir, "output")
	
	// Add .txt extension to the output file
	outputFile := tempOutput + ".txt"
	
	// Use pdftotext to extract text from PDF
	cmd := exec.Command("pdftotext", "-layout", filePath, outputFile)
	
	// Capture command output for better error reporting
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to extract text with pdftotext: %v - %s", err, stderr.String())
	}

	// Read the extracted text
	content, err := os.ReadFile(outputFile)
	if err != nil {
		return "", fmt.Errorf("failed to read extracted text from %s: %v", outputFile, err)
	}

	// Clean up the temporary file
	err = os.Remove(outputFile)
	if err != nil {
		// Log the error but don't fail the function because of it
		log.Printf("Warning: failed to remove temporary file %s: %v", outputFile, err)
	}

	text := string(content)
	if text == "" {
		return "", fmt.Errorf("no text content found in PDF")
	}

	return text, nil
}