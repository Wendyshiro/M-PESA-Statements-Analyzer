package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"mpesa-finance/internal/middleware"
	"mpesa-finance/internal/repository"
)

type JobHandler struct {
	jobRepo *repository.JobRepository
}

func NewJobHandler(jobRepo *repository.JobRepository) *JobHandler {
	return &JobHandler{jobRepo: jobRepo}
}

type JobStatusResponse struct {
	JobID            string  `json:"job_id"`
	Status           string  `json:"status"`
	OriginalFilename string  `json:"original_filename"`
	ErrorMessage     string  `json:"error_message,omitempty"`
	CreatedAt        string  `json:"created_at"`
	CompletedAt      *string `json:"completed_at,omitempty"`
}

func (h *JobHandler) GetJobStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, "Method not allowed", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
		return
	}
	//get user from context

	claims, ok := middleware.GetClaims(r)
	if !ok {
		respondError(w, "Unauthorized", "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}
	//get id from url path
	path := strings.TrimPrefix(r.URL.Path, "/jobs/")
	if path == "" || path == r.URL.Path {
		respondError(w, "Job ID required", "INVALID_REQUEST", http.StatusBadRequest)
		return
	}
	jobID := path
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	//get job from database
	job, err := h.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		respondError(w, "Job not found", "NOT_FOUND", http.StatusNotFound)
		return
	}
	//verify job belongs to user
	if job.UserID != claims.UserID {
		respondError(w, "Access denied", "FORBIDDEN", http.StatusForbidden)
		return
	}

	var completedAt *string
	if job.CompletedAt != nil {
		completed := job.CompletedAt.Format(time.RFC3339)
		completedAt = &completed
	}
	response := JobStatusResponse{
		JobID:            job.ID,
		Status:           string(job.Status),
		OriginalFilename: job.OriginalFilename,
		ErrorMessage:     job.ErrorMessage,
		CreatedAt:        job.CreatedAt.Format(time.RFC3339),
		CompletedAt:      completedAt,
	}
	respondJSON(w, response, http.StatusOK)
}

func (h *JobHandler) GetUserJobs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, "Method not allowed", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
		return
	}

	// Get user from context
	claims, ok := middleware.GetClaims(r)
	if !ok {
		respondError(w, "Unauthorized", "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Get user's jobs (limit to last 50)
	jobs, err := h.jobRepo.GetByUserID(ctx, claims.UserID, 50)
	if err != nil {
		respondError(w, "Failed to fetch jobs", "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}

	var response []JobStatusResponse
	for _, job := range jobs {
		var completedAt *string
		if job.CompletedAt != nil {
			completed := job.CompletedAt.Format(time.RFC3339)
			completedAt = &completed
		}

		response = append(response, JobStatusResponse{
			JobID:            job.ID,
			Status:           string(job.Status),
			OriginalFilename: job.OriginalFilename,
			ErrorMessage:     job.ErrorMessage,
			CreatedAt:        job.CreatedAt.Format(time.RFC3339),
			CompletedAt:      completedAt,
		})
	}

	respondJSON(w, response, http.StatusOK)
}
