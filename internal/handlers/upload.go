package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"mpesa-finance/utils"

	"github.com/gin-gonic/gin"
)

// Helper function to get first n characters of a string
func firstN(s string, n int) string {
	if len(s) < n {
		return s
	}
	return s[:n]
}


func UploadPDFHandler(c *gin.Context) {
	// Ensure uploads directory exists
	if err := os.MkdirAll("uploads", 0755); err != nil {
		utils.LogError("Failed to create uploads directory:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process file"})
		return
	}

	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		utils.LogError("No file uploaded or error reading file:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please select a valid PDF file"})
		return
	}

	// Validate file type
	if filepath.Ext(file.Filename) != ".pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only PDF files are allowed"})
		return
	}

	// Create a unique filename to prevent collisions
	timestamp := time.Now().Format("20060102-150405")
	dst := filepath.Join("uploads", fmt.Sprintf("%s_%s", timestamp, filepath.Base(file.Filename)))

	// Save the uploaded file
	if err := c.SaveUploadedFile(file, dst); err != nil {
		utils.LogError("Failed to save uploaded file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Log file details
	fileInfo, _ := os.Stat(dst)
	utils.LogInfo(fmt.Sprintf("Processing PDF: %s (Size: %d bytes)", dst, fileInfo.Size()))

	// Extract text from PDF
	text, err := utils.ExtractTextFromPDF(dst)
	if err != nil {
		utils.LogError("Failed to extract PDF text:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process PDF. Please ensure it's a valid M-Pesa statement.",
		})
		return
	}

	// Log first 200 characters of extracted text for debugging
	utils.LogInfo(fmt.Sprintf("Extracted text (first 200 chars): %s", firstN(text, 200)))

	// Parse transactions from text
	transactions, err := utils.ParseTransactionsFromText(text)
	if err != nil {
		utils.LogError("Failed to parse transactions:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse transaction data from the PDF",
		})
		return
	}

	if len(transactions) == 0 {
		utils.LogError("No transactions found in the PDF", nil)
		c.JSON(http.StatusOK, gin.H{
			"message": "No transactions found in the uploaded statement",
			"summary": gin.H{
				"categories": map[string]float64{},
			},
		})
		return
	}

	// Save transactions to output.txt
	outputFile := "output.txt"
	if err := os.WriteFile(outputFile, []byte(text), 0644); err != nil {
		utils.LogError("Failed to save transactions:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save transaction data",
		})
		return
	}

	// Generate summary
	summary := utils.AnalyzeTransactions(transactions)

	// Log success
	utils.LogInfo(fmt.Sprintf("Successfully processed %d transactions", len(transactions)))

	// Return success response with summary data
	c.JSON(http.StatusOK, gin.H{
		"message":      "Statement processed successfully",
		"transactions": len(transactions),
		"summary": gin.H{
			"categories": summary.CategoryBreakdown,
		},
	})
}
