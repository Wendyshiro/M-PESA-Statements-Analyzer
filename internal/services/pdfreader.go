package services

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// ExtractTextFromPDF extracts text from a PDF, unlocking it first if a password is provided
func ExtractTextFromPDF(filePath string, password string) (string, error) {
	targetPath := filePath

	// If a password is provided, unlock the PDF first using qpdf
	if password != "" {
		unlockedPath := filePath + "_unlocked.pdf"

		cmd := exec.Command("qpdf", "--password="+password, "--decrypt", filePath, unlockedPath)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to unlock PDF (wrong password?): %v - %s", err, stderr.String())
		}

		// Clean up the unlocked file after we're done
		defer os.Remove(unlockedPath)
		targetPath = unlockedPath
	}

	// Extract text from the (possibly unlocked) PDF
	tempDir := filepath.Dir(targetPath)
	tempOutput := filepath.Join(tempDir, "output")
	outputFile := tempOutput + ".txt"

	cmd := exec.Command("pdftotext", "-layout", targetPath, outputFile)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to extract text with pdftotext: %v - %s", err, stderr.String())
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		return "", fmt.Errorf("failed to read extracted text from %s: %v", outputFile, err)
	}

	if err := os.Remove(outputFile); err != nil {
		log.Printf("Warning: failed to remove temporary file %s: %v", outputFile, err)
	}

	text := string(content)
	if text == "" {
		return "", fmt.Errorf("no text content found in PDF")
	}

	return text, nil
}