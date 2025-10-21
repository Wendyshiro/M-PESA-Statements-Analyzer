package utils

import "fmt"

// LogInfo logs an informational message
func LogInfo(message string) {
	fmt.Printf("[INFO] %s\n", message)
}

// LogError logs an error message with optional error details
func LogError(message string, err error) {
	if err != nil {
		fmt.Printf("[ERROR] %s: %v\n", message, err)
	} else {
		fmt.Printf("[ERROR] %s\n", message)
	}
}
