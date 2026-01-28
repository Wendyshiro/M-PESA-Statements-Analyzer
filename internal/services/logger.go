package services

import (
	"fmt"
	"log"
	"os"
	"time"
)

var errorLogFile *os.File
var errorLogger *log.Logger

// InitLogger initializes the error logger with a file
func InitLogger() error {
	logDir := "logs"

	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// Create or append to error log file with timestamp
	logFileName := fmt.Sprintf("%s/errors_%s.log", logDir, time.Now().Format("2006-01-02"))
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	errorLogFile = file
	errorLogger = log.New(file, "[ERROR] ", log.LstdFlags|log.Lshortfile)
	return nil
}

// CloseLogger closes the error log file
func CloseLogger() error {
	if errorLogFile != nil {
		return errorLogFile.Close()
	}
	return nil
}

// LogInfo logs an informational message to console
func LogInfo(message string) {
	fmt.Printf("[INFO] %s\n", message)
}

// LogError logs an error message to both console and file
func LogError(message string, err error) {
	if err != nil {
		output := fmt.Sprintf("%s: %v", message, err)
		fmt.Printf("[ERROR] %s\n", output)
		if errorLogger != nil {
			errorLogger.Println(output)
		}
	} else {
		fmt.Printf("[ERROR] %s\n", message)
		if errorLogger != nil {
			errorLogger.Println(message)
		}
	}
}

// LogDebug logs a debug message to console
func LogDebug(message string) {
	fmt.Printf("[DEBUG] %s\n", message)
}
