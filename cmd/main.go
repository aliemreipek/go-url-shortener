package main

import (
	"log"
	"os"

	"github.com/aliemreipek/go-url-shortener/api"
	"github.com/aliemreipek/go-url-shortener/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func setupRoutes(app *fiber.App) {
	// Group endpoints: /api/v1
	// POST /api/v1 -> Create Short URL
	app.Post("/api/v1", api.ShortenURL)

	// GET /:url -> Redirect Logic
	app.Get("/:url", api.ResolveURL)
}

func main() {
	// 1. Initialize Database & Redis
	// .env file is loaded inside this function, so we don't need to load it again here.
	database.ConnectDB()

	// 2. Initialize Fiber App
	app := fiber.New()

	// Middleware: Logger (Shows requests in terminal)
	app.Use(logger.New())

	// 3. Setup Routes
	setupRoutes(app)

	// 4. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s...", port)
	log.Fatal(app.Listen(":" + port))
}
