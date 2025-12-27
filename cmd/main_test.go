package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aliemreipek/go-url-shortener/api"
	"github.com/aliemreipek/go-url-shortener/database"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// TestShortenURL tests the link shortening functionality of the system.
func TestShortenURL(t *testing.T) {

	// Test runs inside /cmd, so we need to look one level up.
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file for test: %v", err)
	}
	// --------------------------------------------

	// 1. Setup Infrastructure (Database connection is required)
	database.ConnectDB()

	// 2. Initialize Fiber App for testing
	app := fiber.New()
	app.Post("/api/v1", api.ShortenURL)

	// 3. Prepare Test Scenario
	randomShort := "test_" + uuid.New().String()[:5]

	requestBody := map[string]string{
		"url":   "https://github.com/aliemreipek",
		"short": randomShort,
	}
	jsonBody, _ := json.Marshal(requestBody)

	// 4. Create Request
	req := httptest.NewRequest(
		"POST",
		"/api/v1",
		bytes.NewBuffer(jsonBody),
	)
	req.Header.Set("Content-Type", "application/json")

	// 5. Send Request
	resp, err := app.Test(req, -1)

	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}

	// 6. Assertions
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
