package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aliemreipek/go-url-shortener/model"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB holds the PostgreSQL connection instance
var DB *gorm.DB

// Ctx is a global context for Redis operations
var Ctx = context.Background()

// Rdb holds the Redis client instance
var Rdb *redis.Client

// ConnectDB initializes connections for both PostgreSQL and Redis
func ConnectDB() {
	// 1. Load environment variables
	// We don't use log.Fatal here because in Docker/Production, .env file might not exist.
	if err := godotenv.Load(); err != nil {
		fmt.Println("⚠️  No .env file found, relying on system env variables")
	}

	// 2. PostgreSQL Connection Setup
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Istanbul",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Connect to Postgres using GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Show SQL queries in terminal
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	log.Println("✅ Connected to PostgreSQL successfully")

	// Set the global DB instance
	DB = db

	// 3. Auto-Migration (Create Tables)
	log.Println("Running Migrations...")
	err = db.AutoMigrate(&model.Url{})
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}

	// 4. Redis Connection Setup
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	Rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword, // No password set
		DB:       0,             // Use default DB
	})

	// Ping Redis to verify connection
	// We use a timeout context to avoid hanging if Redis is down
	_, cancel := context.WithTimeout(Ctx, 5*time.Second)
	defer cancel()

	pong, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis: ", err)
	}

	log.Println("✅ Connected to Redis successfully:", pong)
}
