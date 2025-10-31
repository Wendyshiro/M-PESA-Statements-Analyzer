package main

import (
	"log"
	"mpesa-finance/handlers"
	"mpesa-finance/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
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
	r.POST("/upload", handlers.UploadPDFHandler)
	r.GET("/summary", handlers.SummaryHandler)
	//ai cateorizer
	r.GET("/ai-categorize", handlers.AICategorizeHandler)

	r.Run(":8081")
}
