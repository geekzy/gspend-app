package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/imam/gspend-app/apps/financial-service/internal/config"
	"github.com/imam/gspend-app/apps/financial-service/pkg/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize MongoDB
	mongoClient, err := database.NewMongoClient(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(nil)
	fmt.Println("Connected to MongoDB!")

	// Initialize Redis
	redisClient, err := database.NewRedisClient(cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()
	fmt.Println("Connected to Redis!")

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health Check
	e.GET("/api/v1/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "ok",
			"service": "financial-service",
			"env":     cfg.AppEnv,
		})
	})

	// Start Server
	fmt.Printf("Financial Service starting on port %s...\n", cfg.Port)
	if err := e.Start(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
