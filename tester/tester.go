package tester

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// URLs to be tested
var urls = []string{
	"http://resume.connorisseur.com",
	"http://resume.connorisseur.com/public/json/about.json",
	"http://resume.connorisseur.com/public/json/contacts.json",
	"http://resume.connorisseur.com/public/json/portfolio.json",
	"http://resume.connorisseur.com/public/json/resume.json",
}

// Logger for logging the results
var logger *log.Logger

func init() {
	// Initialize the logger
	file, err := os.OpenFile("tester.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		os.Exit(1)
	}
	logger = log.New(file, "TESTER: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// FetchURL fetches the content from the given URL
func FetchURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code for %s: %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body for %s: %v", url, err)
	}

	return body, nil
}

// ValidateJSON validates if the byte slice is a valid JSON
func ValidateJSON(data []byte) error {
	var js json.RawMessage
	if err := json.Unmarshal(data, &js); err != nil {
		return fmt.Errorf("invalid JSON: %v", err)
	}
	return nil
}

// TestWebsite performs the tests on the website
func TestWebsite() error {
	for _, url := range urls {
		data, err := FetchURL(url)
		if err != nil {
			return fmt.Errorf("failed to fetch URL %s: %v", url, err)
		}

		if url != "http://resume.connorisseur.com" {
			if err := ValidateJSON(data); err != nil {
				return fmt.Errorf("validation failed for %s: %v", url, err)
			}
		}
	}
	return nil
}

// StartPeriodicTesting starts the testing with an initial delay and periodic execution
func StartPeriodicTesting() {
	// Initial delay of 30 seconds
	time.Sleep(30 * time.Second)
	if err := TestWebsite(); err != nil {
		logger.Println("Website test failed:", err)
	} else {
		logger.Println("Website test passed")
	}

	// Periodic execution every 10 minutes
	ticker := time.NewTicker(60 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if err := TestWebsite(); err != nil {
			logger.Println("Website test failed:", err)
		} else {
			logger.Println("Website test passed")
		}
	}
}
