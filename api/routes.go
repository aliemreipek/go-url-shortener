package api

import (
	"os"
	"strings"
	"time"

	"github.com/aliemreipek/go-url-shortener/database"
	"github.com/aliemreipek/go-url-shortener/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// request represents the incoming JSON payload for creating a new short URL
type request struct {
	URL         string `json:"url"`
	CustomShort string `json:"short"`
	Expiry      uint64 `json:"expiry"` // Expiry in hours
}

// response represents the JSON payload sent back to the client
type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

// ResolveURL handles the redirection from short URL to original URL
// GET /:url
func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")

	// 1. Check Redis (Cache) first for speed
	rdb := database.Rdb
	val, err := rdb.Get(database.Ctx, url).Result()

	// If err is nil, it means the key exists in Redis!
	if err == nil {
		return c.Redirect(val, 301)
	}

	// If err is NOT nil, it could be redis.Nil (not found) or a connection error.
	// In both cases, we fall back to the Database.

	// 2. If not in Redis, check PostgreSQL (Database)
	var urlModel model.Url
	// Find the short URL in the database
	if err := database.DB.Where("url = ?", url).First(&urlModel).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Short URL not found"})
	}

	// 3. Increment Click Count (Async)
	go func() {
		urlModel.Clicked++
		database.DB.Save(&urlModel)
	}()

	// 4. Save to Redis for future requests (Cache Aside Pattern)
	// We cache it for 24 hours to reduce DB load
	rdb.Set(database.Ctx, url, urlModel.Redirect, 24*time.Hour)

	return c.Redirect(urlModel.Redirect, 301)
}

// ShortenURL handles creating a new short URL
// POST /api/v1
func ShortenURL(c *fiber.Ctx) error {
	body := new(request)

	// Parse incoming JSON
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Helper: Enforce HTTP prefix (SECURE FIX)
	// We default to HTTPS for security best practices
	if !strings.HasPrefix(body.URL, "http") {
		body.URL = "https://" + body.URL
	}

	// 1. Check for infinite loop (prevent shortening the domain itself)
	// In production, use os.Getenv("DOMAIN")
	if strings.Contains(body.URL, "localhost:3000") {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You cannot hack this system :)"})
	}

	var id string

	// 2. Custom Short URL Logic
	if body.CustomShort != "" {
		id = body.CustomShort
		// Check if this custom alias is already taken in DB
		var temp model.Url
		if err := database.DB.Where("url = ?", id).First(&temp).Error; err == nil {
			// If no error, it means we found one -> Forbidden
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "URL custom alias is already in use"})
		}
	} else {
		// 3. Generate Random Short URL (UUID)
		// We take the first 6 characters of a UUID for brevity
		id = uuid.New().String()[:6]
	}

	// 4. Save to Database
	newURL := model.Url{
		URL:       id,
		Redirect:  body.URL,
		Random:    body.CustomShort == "",
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&newURL).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to connect to server"})
	}

	// 5. Save to Redis (Optional: Write-Through Cache)
	// We save it to Redis immediately so it's ready for the first click
	// Default expiry is 24 hours
	database.Rdb.Set(database.Ctx, id, body.URL, 24*time.Hour)

	// Prepare Response
	resp := response{
		URL:             body.URL,
		CustomShort:     os.Getenv("DOMAIN") + "/" + id,
		Expiry:          24 * time.Hour,
		XRateRemaining:  10, // Placeholder for rate limiting
		XRateLimitReset: 30, // Placeholder
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
