package services

import (
	"fmt"
	"strconv"
	"strings"

	"mpesa-finance/internal/models"
)

// ParseTransactionsFromText parses extracted PDF text into Transaction structs
func ParseTransactionsFromText(text string) ([]models.Transaction, error) {
	fmt.Println("=== PARSING TRANSACTIONS ===")
	fmt.Printf("Input text (first 500 chars):\n%.500s\n...\n", text)

	lines := strings.Split(text, "\n")
	var transactions []models.Transaction
	startParsing := false
	var headers []string

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Detect the header line to start parsing
		if strings.Contains(line, "Receipt") && 
		   strings.Contains(line, "Time") && 
		   strings.Contains(line, "Details") {
			startParsing = true
			headers = splitBySpaces(line)
			fmt.Printf("Found headers at line %d: %v\n", i, headers)
			continue
		}

		if !startParsing {
			continue
		}

		// Skip page numbers and other non-transaction lines
		if strings.HasPrefix(line, "Page ") || 
		   strings.HasPrefix(line, "Balance B/") ||
		   strings.HasPrefix(line, "M-PESA") ||
		   strings.Contains(line, "Statement") {
			continue
		}

		// Try to parse the line as a transaction
		fields := splitBySpaces(line)
		fmt.Printf("Line %d: %s\n", i, line)
		fmt.Printf("Fields (%d): %v\n", len(fields), fields)

		// Try to find transaction patterns
		if len(fields) >= 7 {
			// Try to parse as a standard transaction line
			t, err := parseTransactionLine(fields)
			if err == nil {
				transactions = append(transactions, t)
				continue
			}
		}

		// If we get here, the line might be a continuation of the previous transaction
		if len(transactions) > 0 && len(fields) > 0 {
			// Append to the details of the last transaction
			lastIdx := len(transactions) - 1
			transactions[lastIdx].Details += " " + strings.Join(fields, " ")
		}
	}

	fmt.Printf("Successfully parsed %d transactions\n", len(transactions))
	return transactions, nil
}

// parseTransactionLine attempts to parse a line into a Transaction
func parseTransactionLine(fields []string) (models.Transaction, error) {
	// Expected format: [ReceiptNo] [Date] [Time] [Details...] [PaidIn] [Withdrawn] [Balance]
	if len(fields) < 7 {
		return models.Transaction{}, fmt.Errorf("not enough fields")
	}

	// Find the split between details and amounts
	// Amounts are the last 3 fields
	amountStart := len(fields) - 3
	
	// Extract amounts
	paidIn := parseFloat(fields[amountStart])
	withdrawn := parseFloat(fields[amountStart+1])
	balance := parseFloat(fields[amountStart+2])

	// The rest is receipt number, date, time, and details
	receiptNo := fields[0]
	completionTime := fields[1] + " " + fields[2]
	details := strings.Join(fields[3:amountStart], " ")

	return models.Transaction{
		ReceiptNo:         receiptNo,
		CompletionTime:    completionTime,
		Details:           details,
		TransactionStatus: "Completed",
		PaidIn:            paidIn,
		Withdrawn:         withdrawn,
		Balance:           balance,
	}, nil
}

// parseFloat safely converts a string to float64
func parseFloat(s string) float64 {
	s = strings.ReplaceAll(s, ",", "")
	if s == "" {
		return 0
	}
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

// splitBySpaces splits a line into fields using multiple spaces as separator
func splitBySpaces(line string) []string {
	// First, normalize multiple spaces to single spaces
	line = strings.Join(strings.Fields(line), " ")
	// Then split on spaces
	return strings.Fields(line)
}
