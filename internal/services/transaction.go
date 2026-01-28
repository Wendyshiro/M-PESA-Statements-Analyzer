package models

// Transaction represents a single M-Pesa transaction
type Transaction struct {
	ReceiptNo         string  `json:"receipt_no"`
	CompletionTime    string  `json:"completion_time"`
	Details           string  `json:"details"`
	TransactionStatus string  `json:"transaction_status"`
	PaidIn            float64 `json:"paid_in"`
	Withdrawn         float64 `json:"withdrawn"`
	Balance           float64 `json:"balance"`
}
