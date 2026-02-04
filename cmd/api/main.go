package main

import (
	"log"
	"net/http"
	"os"

	"mpesa-finance/config"
	"mpesa-finance/internal/auth"
	"mpesa-finance/internal/handlers"
	"mpesa-finance/internal/middleware"
	utils "mpesa-finance/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	//Create services
	authService := auth.NewService(cfg.JWTSecret)

	//Create handlers
	uploadHandler := handlers.NewUploadHandler(cfg.UploadDir)
	authHandler := handlers.NewAuthHandler(authService)
	//Create router
	mux := http.NewServeMux()
	// Register routes

	mux.HandleFunc("/upload", uploadHandler.HandleUpload)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)

	// Protected routes auth required
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/upload", uploadHandler.HandleUpload)

	// Wrap with middleware
	var handler http.Handler = mux
	handler = middleware.SecurityHeaders(handler)
	handler = middleware.CORS([]string{"http://localhost:3000"})(handler)

	//Start server

	addr := ":" + cfg.Port
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	//load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey != "" {
		var err error
		handlers.AICategorizer, err = utils.NewAICategorizer(apiKey)
		if err != nil {
			log.Printf("Warning: Failed to initialize AI categorizer: %v", err)
		} else {
			log.Println("AI categorizer initialized successfully")
		}
	} else {
		log.Println("Warning: OPENAI_API_KEY not set. AI categorization will be disabled.")
	}
	r.MaxMultipartMemory = 8 << 20 // 8 MB

	// Serve static files from the dashboard directory
	r.Static("/static", "./dashboard/static")
	r.StaticFile("/favicon.ico", "./dashboard/favicon.ico")

	// Handle dashboard route - serve index.html for all dashboard routes
	r.GET("/dashboard", func(c *gin.Context) {
		c.File("./dashboard/index.html")
	})

	// Serve index.html at root
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/dashboard")
	})

	// API endpoints

	r.POST("/upload", func(c *gin.Context) {
		uploadHandler.HandleUpload(c.Writer, c.Request)
	})
	r.GET("/summary", handlers.SummaryHandler)
	//ai cateorizer
	r.GET("/ai-categorize", handlers.AICategorizeHandler)

	r.Run(":8081")
}
