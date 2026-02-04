package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"mpesa-finance/internal/middleware"

	"github.com/google/uuid"
)

type UploadHandler struct {
	uploadDir string
}

func NewUploadHandler(uploadDir string) *UploadHandler {
	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	return &UploadHandler{
		uploadDir: uploadDir,
	}
}

type UploadResponse struct {
	JobID    string `json:"job_id"`
	Message  string `json:"message"`
	Filename string `json:"filename"`
}

func (h *UploadHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	// Only allow POST
	if r.Method != http.MethodPost {
		respondError(w, "Method not allowed", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (32MB max memory)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		respondError(w, "Failed to parse form data", "INVALID_FORM", http.StatusBadRequest)
		return
	}

	// Get file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		respondError(w, "No file provided", "NO_FILE", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file
	if err := middleware.ValidateFileUpload(file, header); err != nil {
		respondError(w, err.Error(), "INVALID_FILE", http.StatusBadRequest)
		return
	}

	// Generate unique filename
	jobID := uuid.New().String()
	sanitizedName := middleware.SanitizeFilename(header.Filename)
	filename := fmt.Sprintf("%s_%s", jobID, sanitizedName)
	filepath := filepath.Join(h.uploadDir, filename)

	// Save file
	dst, err := os.Create(filepath)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		respondError(w, "Failed to save file", "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		log.Printf("Failed to write file: %v", err)
		respondError(w, "Failed to save file", "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}

	// TODO: Queue job for processing (we'll implement this later)
	log.Printf("File uploaded successfully: %s (job: %s)", sanitizedName, jobID)

	// Respond with success
	response := UploadResponse{
		JobID:    jobID,
		Message:  "File uploaded successfully and queued for processing",
		Filename: sanitizedName,
	}

	respondJSON(w, response, http.StatusAccepted)
}
