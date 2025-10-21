package main

import (
	"mpesa-finance/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
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

	r.Run(":8080")
}
