package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"mpesa-finance/internal/repository"
	"mpesa-finance/internal/services"
)

type SummaryHandler struct {
	jobRepo *repository.JobRepository
}

func NewSummaryHandler(jobRepo *repository.JobRepository) *SummaryHandler {
	return &SummaryHandler{jobRepo: jobRepo}
}

func (h *SummaryHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	// Extract job ID from URL: /summary/{jobId}
	jobID := strings.TrimPrefix(r.URL.Path, "/summary/")
	if jobID == "" {
		respondError(w, "Job ID required", "BAD_REQUEST", http.StatusBadRequest)
		return
	}

	// Get job from database
	job, err := h.jobRepo.GetByID(r.Context(), jobID)
	if err != nil {
		respondError(w, "Job not found", "NOT_FOUND", http.StatusNotFound)
		return
	}

	// Job must be completed
	if job.Status != "completed" {
		respondError(w, "Job is not completed yet", "JOB_NOT_COMPLETE", http.StatusBadRequest)
		return
	}

	// Resolve absolute path if needed
	filePath := job.FilePath
	if !filepath.IsAbs(filePath) {
		wd, err := os.Getwd()
		if err != nil {
			respondError(w, "Failed to resolve file path", "INTERNAL_ERROR", http.StatusInternalServerError)
			return
		}
		filePath = filepath.Join(wd, filePath)
	}

	// Extract text from the PDF
	log.Printf("Summary: extracting from path=%s", filePath)
	text, err := services.ExtractTextFromPDF(filePath, job.PDFPassword)
	if err != nil {
		log.Printf("Summary: extraction error: %v", err)
		respondError(w, "Failed to read statement: "+err.Error(), "EXTRACTION_ERROR", http.StatusInternalServerError)
		return
	}
	log.Printf("Summary: extracted %d chars of text", len(text))

	// Parse transactions
	transactions, err := services.ParseTransactionsFromText(text)
	if err != nil {
		respondError(w, "Failed to parse transactions", "PARSE_ERROR", http.StatusInternalServerError)
		return
	}
	log.Printf("Summary: parsed %d transactions", len(transactions))

	// Analyze
	summary := services.AnalyzeTransactions(transactions)

	respondJSON(w, map[string]interface{}{
		"message": "Summary retrieved successfully",
		"summary": map[string]interface{}{
			"categories":     summary.CategoryBreakdown,
			"total_income":   summary.TotalIncome,
			"total_expenses": summary.TotalExpenses,
			"net_balance":    summary.NetBalanceChange,
		},
		"total_transactions": len(transactions),
	}, http.StatusOK)
}