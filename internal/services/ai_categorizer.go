package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type AICategorizer struct {
	client      *openai.Client
	model       string
	categories  []string
	temperature float32
}

func NewAICategorizer(apiKey string) (*AICategorizer, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}
	return &AICategorizer{
		client:      openai.NewClient(apiKey),
		model:       "gpt-3.5-turbo",
		temperature: 0.3,
		categories: []string{
			"Airtime & Data",
			"Shopping",
			"Utilities",
			"Merchant Payments",
			"Withdrawals",
			"Send Money",
			"Safaricom Services",
			"Other Expenses",
		},
	}, nil

}

func (c *AICategorizer) Categorize(transaction string) (string, float32, error) {
	// prepare the prompt for the AI
	prompt := fmt.Sprintf(`Categorize the following M-Pesa transaction into one of these categories: %s

Transaction: %s

Return a JSON object with:
{
  "category": "The most appropriate category",
  "confidence": 0.0 to 1.0,
  "reason": "Brief explanation"
}`,
		strings.Join(c.categories, ", "),
		transaction,
	)
	//make the api call
	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: c.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a helpful financial assistant that catergorizes M-PESA transactions.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: c.temperature,
		},
	)
	//response handling
	if err != nil {
		return "", 0, fmt.Errorf("AI categorization failed: %v", err)
	}
	//parse the response
	var result struct {
		Category   string  `json:"category"`
		Confidence float32 `json:"confidence"`
		Reason     string  `json:"reason"`
	}
	//extract JSON from the response
	content := resp.Choices[0].Message.Content
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return "", 0, fmt.Errorf("failed to parse AI response: %v", err)
	}
	//validate the category is in our list
	for _, cat := range c.categories {
		if strings.EqualFold(cat, result.Category) {
			return cat, result.Confidence, nil
		}
	}
	return "Other Expenses", 1.0, nil

}

func (c *AICategorizer) CategorizeWithFallback(details string) string {
	if details == "" {
		return "Uncategorized"
	}
	// first try ai categorization
	category, confidence, err := c.Categorize(details)
	if err != nil {
		log.Printf("AI categorization failed, falling back to rules: %v", err)
		return categorizeTransaction(details)
	}
	//if confidence is low, fall back to rule base
	const confidenceThresfold = 0.7
	if confidence < confidenceThresfold {
		log.Printf("Low Confidence (%.2f) for category '%s', falling back to rules")
		return categorizeTransaction(details)

	}
	return category
}
