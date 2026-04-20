package models
import "time"

type JobStatus string

const (
	JobStatusQueued  JobStatus = "queued"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed  JobStatus = "failed"
)

type Job struct {
	ID   string `json:"id"`
	UserID  string `json:"user_id"`
	FilePath  string `json:"file_path"`
	OriginalFilename string `json:"original_filename"`
	Status JobStatus `json:"status"`
	ErrorMessage string `json:"error_message, omitempty"`
	PDFPassword  string `json:"pdf_password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at, omitempty"`
}