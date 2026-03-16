package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"mpesa-finance/internal/models"

	"github.com/redis/go-redis/v9"
)

const (
	QueueKey = "job_queue"
)

type JobQueue struct {
	client *redis.Client
}

func NewJobQueue(redisAddr string) (*JobQueue, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
	//TEST CONNECTION
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis:%w", err)
	}
	return &JobQueue{client: client}, nil
}

// enqueue adds a job to the queue
func (q *JobQueue) Enqueue(ctx context.Context, job *models.Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("Failed to marshal job: %w", err)

	}
	//add to redis list
	return q.client.RPush(ctx, QueueKey, data).Err()
}

// Dequeu removes and return the next job from the queue
func (q *JobQueue) Dequeue(ctx context.Context, timeout time.Duration) (*models.Job, error) {
	//blpop blocks until a job is available or timeout
	result, err := q.client.BLPop(ctx, timeout, QueueKey).Result()
	if err == redis.Nil {
		return nil, nil //timeout no job available
	}
	if err != nil {
		return nil, err
	}
	//result is [queueName, data]
	if len(result) < 2 {
		return nil, fmt.Errorf("unexpected result from queue")
	}
	var job models.Job
	if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job: %w", err)
	}
	return &job, nil

}

//size returns the number of jobs in the queue

func (q *JobQueue) Size(ctx context.Context) (int64, error) {
	return q.client.LLen(ctx, QueueKey).Result()

}

// clear removes all jobs from the queue
func (q *JobQueue) Clear(ctx context.Context) error {
	return q.client.Del(ctx, QueueKey).Err()
}

func (q *JobQueue) Close() error {
	return q.client.Close()
}
