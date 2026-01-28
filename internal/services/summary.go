package services

import (
	"fmt"
	"mpesa-finance/internal/models"
	"strings"
)

type Summary struct {
	TotalIncome       float64            `json:"total_income"`
	TotalExpenses     float64            `json:"total_expenses"`
	NetBalanceChange  float64            `json:"net_balance_change"`
	CategoryBreakdown map[string]float64 `json:"categories"`
}

// AnalyzeTransactions generates a summary of income, expenses and categories
func AnalyzeTransactions(transactions []models.Transaction) Summary {
	summary := Summary{
		CategoryBreakdown: make(map[string]float64),
	}

	// Log category distribution for debugging
	logTransactionCategories(transactions)

	for _, t := range transactions {
		// Track total income vs expenses
		if t.PaidIn > 0 {
			summary.TotalIncome += t.PaidIn
		}
		if t.Withdrawn > 0 {
			summary.TotalExpenses += t.Withdrawn
		}

		// Categorize based on keywords in Details
		category := categorizeTransaction(t.Details)
		amount := t.PaidIn
		if amount == 0 {
			amount = t.Withdrawn
		}
		summary.CategoryBreakdown[category] += amount

		// Log large or unusual transactions
		if amount > 10000 { // Log transactions over 10,000 KES
			fmt.Printf("Large transaction: %.2f KES - %s (%s)\n", amount, t.Details, category)
		}
	}

	summary.NetBalanceChange = summary.TotalIncome - summary.TotalExpenses

	// Log summary statistics
	fmt.Printf("\n=== Transaction Summary ===\n")
	fmt.Printf("Total Income: %.2f KES\n", summary.TotalIncome)
	fmt.Printf("Total Expenses: %.2f KES\n", summary.TotalExpenses)
	fmt.Printf("Net Change: %.2f KES\n", summary.NetBalanceChange)
	fmt.Printf("Categories: %+v\n", summary.CategoryBreakdown)

	return summary
}

// categorizeTransaction identifies category from the Details text
func categorizeTransaction(details string) string {
	if details == "" {
		return "Uncategorized"
	}

	details = strings.ToLower(details)

	// Common M-Pesa patterns
	switch {
	//airtime
	case strings.Contains(details, "airtime"),
		strings.Contains(details, "data"),
		strings.Contains(details, "bundle"):
		return "Airtime & Data"
		//shopping
	case strings.Contains(details, "till"),
		strings.Contains(details, "pos"),
		strings.Contains(details, "shop"),
		strings.Contains(details, "naivas"),
		strings.Contains(details, "quickmart"):
		return "Shopping"

	//utilities
	case strings.Contains(details, "bill"),
		strings.Contains(details, "utility"),
		strings.Contains(details, "water"),
		strings.Contains(details, "electricity"),
		strings.Contains(details, "internet"),
		strings.Contains(details, "pharmaceuticals"),
		strings.Contains(details, "gas"):
		return "Utilities"

	//food
	case strings.Contains(details, "restaurant"),
		strings.Contains(details, "cafe"),
		strings.Contains(details, "java"),
		strings.Contains(details, "cj"),
		strings.Contains(details, "mint & salt"):
		return "Food & Dining"

	case strings.Contains(details, "till"),
		strings.Contains(details, "pay bill"),
		strings.Contains(details, "merchant"),
		strings.Contains(details, "pos"):
		return "Merchant Payments"

	case strings.Contains(details, "received from"),
		strings.Contains(details, "from "),
		strings.Contains(details, "deposit"),
		strings.Contains(details, "absa"),
		strings.Contains(details, "sent by"):
		return "Money Received"

	case strings.Contains(details, "withdraw"),
		strings.Contains(details, "agent"),
		strings.Contains(details, "atm"):
		return "Cash Withdrawals"

	case strings.Contains(details, "send money"),
		strings.Contains(details, "sent to"),
		strings.Contains(details, "to "),
		strings.Contains(details, "transfer"):
		return "Send Money"

	case strings.Contains(details, "m-shwari"),
		strings.Contains(details, "fuliza"),
		strings.Contains(details, "loan"),
		strings.Contains(details, "save"):
		return "Loans & Savings"

	case strings.Contains(details, "safaricom"),
		strings.Contains(details, "mpesa"):
		return "Safaricom Services"

	case strings.Contains(details, "bill"),
		strings.Contains(details, "payment"):
		return "Bills & Utilities"

	default:
		return "Other Expenses"
	}
}

// logTransactionCategories helps with debugging by showing category distribution
func logTransactionCategories(transactions []models.Transaction) {
	categoryCount := make(map[string]int)

	for _, t := range transactions {
		category := categorizeTransaction(t.Details)
		categoryCount[category]++
	}

	fmt.Println("\n=== Transaction Category Distribution ===")
	for cat, count := range categoryCount {
		fmt.Printf("%s: %d transactions\n", cat, count)
	}
}
