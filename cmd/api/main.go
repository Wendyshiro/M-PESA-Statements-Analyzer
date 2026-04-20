package main

import (
	"log"
	"mpesa-finance/cache"
	"mpesa-finance/config"
	"mpesa-finance/internal/auth"
	"mpesa-finance/internal/database"
	"mpesa-finance/internal/handlers"
	"mpesa-finance/internal/middleware"
	"mpesa-finance/internal/repository"
	"mpesa-finance/queue"
	"mpesa-finance/internal/worker"
	"context"
	"net/http"
	"strings"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	//connect to database
	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Connected to database successfully")
	//Connect to redis
	redisCache, err := cache.NewRedisCache(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisCache.Close()
	log.Println("Connected to Redis successfully")
	//Create job queue
	jobQueue, err := queue.NewJobQueue(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to created job queue: %v", err)

	}
	defer jobQueue.Close()
	log.Println("Job queue initialized")

	//create repositories
	userRepo := repository.NewUserRepository(db)
	jobRepo := repository.NewJobRepository(db)

	//Create and start the worker
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w := worker.NewWorker(jobQueue, jobRepo)
	go w.Start(ctx)
	log.Println("Worker started in background")

	//create services
	authService := auth.NewService(cfg.JWTSecret)

	//Create handlers
	uploadHandler := handlers.NewUploadHandler(cfg.UploadDir, jobRepo, jobQueue)
	authHandler := handlers.NewAuthHandler(authService, userRepo)
	jobHandler := handlers.NewJobHandler(jobRepo)
	healthHandler := handlers.NewHealthHandler(db, redisCache)
	summaryHandler := handlers.NewSummaryHandler(jobRepo)

	//Create router
	mux := http.NewServeMux()
	// Register routes
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request){
		healthHandler.Check(w, r)
	})
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)

	// Protected routes auth required
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/upload", uploadHandler.HandleUpload)
	protectedMux.HandleFunc("/summary/", summaryHandler.GetSummary)
	protectedMux.HandleFunc("/jobs", jobHandler.GetUserJobs)
	protectedMux.HandleFunc("/jobs/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/jobs/") && len(strings.TrimPrefix(r.URL.Path, "/jobs/")) > 0 {
			jobHandler.GetJobStatus(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	mux.Handle("/", middleware.AuthMiddleware(authService)(protectedMux))

	// Wrap with middleware
	var handler http.Handler = mux
	handler = middleware.SecurityHeaders(handler)
	handler = middleware.CORS([]string{"http://localhost:3000"})(handler)

	//Start server

	addr := ":" + cfg.Port
	log.Printf("Starting server on %s", addr)
	log.Printf("Starting database server on %s", addr)
	log.Printf("Starting Redis server on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}