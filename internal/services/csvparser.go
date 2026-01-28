package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"mpesa-finance/internal/models"
)

// ParseCSV reads an M-Pesa CSV file and returns a list of Transaction structs
func ParseCSV(filePath string) ([]models.Transaction, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // handle variable columns

	var transactions []models.Transaction
	line := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV: %v", err)
		}

		line++
		// Skip header or malformed rows
		if line == 1 || len(record) < 7 {
			continue
		}

		// Clean each cell
		for i := range record {
			record[i] = strings.TrimSpace(record[i])
		}

		// Convert numbers safely (handle commas and blanks)
		parseFloat := func(s string) float64 {
			s = strings.ReplaceAll(s, ",", "")
			if s == "" {
				return 0
			}
			f, _ := strconv.ParseFloat(s, 64)
			return f
		}

		transaction := models.Transaction{
			ReceiptNo:         record[0],
			CompletionTime:    record[1],
			Details:           record[2],
			TransactionStatus: record[3],
			PaidIn:            parseFloat(record[4]),
			Withdrawn:         parseFloat(record[5]),
			Balance:           parseFloat(record[6]),
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
