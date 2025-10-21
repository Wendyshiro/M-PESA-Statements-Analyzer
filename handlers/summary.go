package handlers

import (
	"fmt"
	"io/ioutil"
	"mpesa-finance/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// FirstN returns the first n characters of a string
func FirstN(s string, n int) string {
	if len(s) < n {
		return s
	}
	return s[:n]
}

func SummaryHandler(c *gin.Context) {
	// Check if output.txt exists
	outputFile := "output.txt"
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		utils.LogInfo("No output.txt found, returning empty summary")
		c.JSON(http.StatusOK, gin.H{
			"message": "No transaction data available. Please upload a statement first.",
			"summary": gin.H{
				"categories": map[string]float64{},
			},
		})
		return
	}

	// Read the transaction data file
	data, err := ioutil.ReadFile(outputFile)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to read transaction data: %v", err)
		utils.LogError(errMsg, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to read transaction data",
			"details": err.Error(),
		})
		return
	}

	// Check if file is empty
	if len(data) == 0 {
		utils.LogInfo("Output file is empty")
		c.JSON(http.StatusOK, gin.H{
			"message": "No transaction data available in the statement",
			"summary": gin.H{
				"categories": map[string]float64{},
			},
		})
		return
	}

	// Parse text into structured transactions
	transactions, err := utils.ParseTransactionsFromText(string(data))
	if err != nil {
		utils.LogError("Failed to parse transaction data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to parse transaction data",
			"details": err.Error(),
		})
		return
	}

	if len(transactions) == 0 {
		utils.LogInfo("No transactions found in the parsed data")
		c.JSON(http.StatusOK, gin.H{
			"message": "No transactions found in the statement",
			"summary": gin.H{
				"categories": map[string]float64{},
			},
		})
		return
	}

	// Log first few transactions for debugging
	utils.LogInfo(fmt.Sprintf("Processing %d transactions", len(transactions)))
	if len(transactions) > 0 {
		utils.LogInfo("Sample transaction: " + fmt.Sprintf("%+v", transactions[0]))
	}

	// Analyze transactions
	summary := utils.AnalyzeTransactions(transactions)

	// Log summary for debugging
	utils.LogInfo("Generated category breakdown:")
	for category, amount := range summary.CategoryBreakdown {
		utils.LogInfo(fmt.Sprintf("%s: %.2f", category, amount))
	}

	// Format the response
	categoryBreakdown := make(map[string]float64)
	for category, amount := range summary.CategoryBreakdown {
		// Convert from cents to shillings if needed
		categoryBreakdown[category] = amount / 100.0
	}

	// Return JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction summary retrieved successfully",
		"summary": gin.H{
			"categories":      categoryBreakdown,
			"total_income":    summary.TotalIncome / 100.0,
			"total_expenses":  summary.TotalExpenses / 100.0,
			"net_balance":     (summary.TotalIncome - summary.TotalExpenses) / 100.0,
		},
		"total_transactions": len(transactions),
	})
}
