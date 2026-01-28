package handlers

import (
	"log"
	"mpesa-finance/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var AICategorizer *utils.AICategorizer

func init() {
	//initialize the AI Categorizer with the api key from environment variable
	var err error
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Printf("WARNING: OPENAI_API_KEY environment variable not set. AI Categorization will not work")

	} else {
		AICategorizer, err = utils.NewAICategorizer(apiKey)
		if err != nil {
			log.Printf("Failed to initialize AI categorizer: %v", err)

		}
	}
}

// aicategorizeHandler handles the AI Categorization of transactions
func AICategorizeHandler(c *gin.Context) {
	//get the transaction description from query parameters
	description := c.Query("description")
	if description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "description query parameter is required"})
		return
	}
	if AICategorizer == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "AI Categorization service is not available",
			"details": "AI Categorization service is not available",
		})
		return
	}
	//Use the AI categorizer
	category, confidence, err := AICategorizer.Categorize(description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to categorize transaction",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"category":   category,
		"confidence": confidence,
	})
}
