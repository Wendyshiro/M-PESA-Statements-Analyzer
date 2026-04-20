package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"path/filepath"

	"mpesa-finance/internal/middleware"
	"mpesa-finance/internal/models"
	"mpesa-finance/internal/repository"
	"mpesa-finance/queue"

	"github.com/google/uuid"
)

type UploadHandler struct {
	uploadDir string
	jobRepo *repository.JobRepository
	jobQueue *queue.JobQueue
}

func NewUploadHandler(uploadDir string, jobRepo *repository.JobRepository, jobQueue *queue.JobQueue) *UploadHandler {

    if err := os.MkdirAll(uploadDir, 0755); err != nil {
        log.Fatalf("Failed to create upload directory: %v", err)
    }

    return &UploadHandler{
        uploadDir: uploadDir,
        jobRepo:   jobRepo,
        jobQueue:  jobQueue,
    }
}

type UploadResponse struct {
	JobID    string `json:"job_id"`
	Message  string `json:"message"`
	Filename string `json:"filename"`
	Status string `json:"status"`
}

func (h *UploadHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	// Only allow POST
	if r.Method != http.MethodPost {
		respondError(w, "Method not allowed", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
		return
	}
	//get user from context
	claims, ok := middleware.GetClaims(r)
	if !ok {
		respondError(w,"Unauthorized", "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}

	// Parse multipart form (32MB max memory)
	
    if err := r.ParseMultipartForm(32 << 20); err != nil {
    log.Printf("ParseMultipartForm error: %v", err)
    log.Printf("Content-Type: %s", r.Header.Get("Content-Type"))
    respondError(w, "Failed to parse form data: "+err.Error(), "INVALID_FORM", http.StatusBadRequest)
    return
}
	//get optional PDF Password
	pdfPassword := r.FormValue("password")

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
	filePath := filepath.Join(h.uploadDir, filename)

	// Save file
	dst, err := os.Create(filePath)
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
	//create job in database
	job := &models.Job{
		ID: jobID,
		UserID: claims.UserID,
		FilePath: filePath,
		OriginalFilename: sanitizedName,
		Status: models.JobStatusQueued,
		PDFPassword: pdfPassword,
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.jobRepo.Create(ctx, job); err != nil {
		log.Printf("Failed to create job: %v", err)
		respondError(w, "Failed to create job", "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}
	//add job to queue for background processing
	if err := h.jobQueue.Enqueue(ctx, job); err != nil{
		log.Printf("Failed to enqueue job: %v", err)
		//job is in db but not queued - we could have a cleanup process for this
		respondError(w, "Failed to queue job", "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}
	log.Printf("Job created:%s (user: %s, file: %s)", jobID, claims.UserID, sanitizedName)
	response := UploadResponse {
		JobID: jobID,
		Message: "File uploaded successfully and queued for processing",
		Filename: sanitizedName,
		Status: string(models.JobStatusQueued),
	}
	respondJSON(w, response, http.StatusAccepted)

}

