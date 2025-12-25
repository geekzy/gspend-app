package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/geekzy/gspend-app/apps/auth-service/internal/config"
	"github.com/geekzy/gspend-app/apps/auth-service/internal/grpc"
	"github.com/geekzy/gspend-app/apps/auth-service/internal/handler"
	"github.com/geekzy/gspend-app/apps/auth-service/internal/middleware"
	"github.com/geekzy/gspend-app/apps/auth-service/internal/repository"
	"github.com/geekzy/gspend-app/apps/auth-service/internal/service"
	"github.com/geekzy/gspend-app/apps/auth-service/pkg/database"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
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

	// Initialize Repository
	userRepo := repository.NewMongoUserRepository(mongoClient.Database(cfg.MongoDatabase))

	// Initialize Service
	authService := service.NewAuthService(userRepo, &cfg)

	// Initialize Handler
	authHandler := handler.NewAuthHandler(authService)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// Routes
	v1 := e.Group("/api/v1/auth")
	v1.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "ok",
			"service": "auth-service",
			"env":     cfg.AppEnv,
		})
	})
	v1.POST("/register", authHandler.Register)
	v1.POST("/login", authHandler.Login)
	v1.GET("/me", authHandler.GetProfile, middleware.AuthMiddleware(&cfg))

	// Start gRPC Server
	authGRPCService := grpc.NewAuthGRPCService(authService)
	grpcServer := grpc.NewGRPCServer(authGRPCService, 9091) // Hardcoded 9091 for now, can move to config

	go func() {
		if err := grpcServer.Start(); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	// Start Server
	fmt.Printf("Auth Service starting on port %s...\n", cfg.Port)
	if err := e.Start(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
