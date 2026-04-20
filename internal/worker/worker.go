package worker

import (
	"context"
	"log"
	"time"

	"mpesa-finance/internal/models"
	"mpesa-finance/internal/repository"
	"mpesa-finance/internal/services"
	"mpesa-finance/queue"
	"strings"
)

type Worker struct {
	jobQueue *queue.JobQueue
	jobRepo  *repository.JobRepository
}

func NewWorker(jobQueue *queue.JobQueue, jobRepo *repository.JobRepository) *Worker {
	return &Worker{
		jobQueue: jobQueue,
		jobRepo:  jobRepo,
	}
}

// Start begins the worker loop. It runs until the context is cancelled.
func (w *Worker) Start(ctx context.Context) {
	log.Println("Worker started, waiting for jobs...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Worker stopping...")
			return
		default:
		}

		job, err := w.jobQueue.Dequeue(ctx, 5*time.Second)
		if err != nil {
			log.Printf("Worker: error dequeuing job: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if job == nil {
			continue
		}

		log.Printf("Worker: picked up job %s (file: %s)", job.ID, job.OriginalFilename)
		w.processJob(ctx, job)
	}
}

// processJob handles a single job end-to-end
func (w *Worker) processJob(ctx context.Context, job *models.Job) {
	err := w.jobRepo.UpdateStatus(ctx, job.ID, models.JobStatusProcessing, "")
	if err != nil {
		log.Printf("Worker: failed to mark job %s as processing: %v", job.ID, err)
		return
	}

	log.Printf("Worker: extracting text from %s", job.FilePath)
	text, err := services.ExtractTextFromPDF(job.FilePath, job.PDFPassword)
	if err != nil {
		log.Printf("Worker: failed to extract text for job %s: %v", job.ID, err)
		if strings.Contains(err.Error(), "Incorrect password"){
			w.failJob(ctx, job.ID, "PDF is password protected. Please re-upload with correct password.")
		}
		w.failJob(ctx, job.ID, "Failed to extract text from PDF: "+err.Error())
		return
	}

	log.Printf("Worker: parsing transactions for job %s", job.ID)
	transactions, err := services.ParseTransactionsFromText(text)
	if err != nil {
		log.Printf("Worker: failed to parse transactions for job %s: %v", job.ID, err)
		w.failJob(ctx, job.ID, "Failed to parse transactions: "+err.Error())
		return
	}

	log.Printf("Worker: job %s completed — parsed %d transactions", job.ID, len(transactions))

	err = w.jobRepo.UpdateStatus(ctx, job.ID, models.JobStatusCompleted, "")
	if err != nil {
		log.Printf("Worker: failed to mark job %s as completed: %v", job.ID, err)
	}
}

// failJob marks a job as failed with an error message
func (w *Worker) failJob(ctx context.Context, jobID string, errMsg string) {
	err := w.jobRepo.UpdateStatus(ctx, jobID, models.JobStatusFailed, errMsg)
	if err != nil {
		log.Printf("Worker: failed to mark job %s as failed: %v", jobID, err)
	}
}