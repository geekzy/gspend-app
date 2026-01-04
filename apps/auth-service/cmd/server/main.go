package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/geekzy/gspend-app/apps/auth-service/internal/config"
	"github.com/geekzy/gspend-app/apps/auth-service/internal/database"
	"github.com/geekzy/gspend-app/apps/auth-service/internal/grpc"
	"github.com/geekzy/gspend-app/apps/auth-service/internal/handler"
	"github.com/geekzy/gspend-app/apps/auth-service/internal/middleware"
	"github.com/geekzy/gspend-app/apps/auth-service/internal/repository"
	"github.com/geekzy/gspend-app/apps/auth-service/internal/service"
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
	mongoClient, err := database.NewMongoClient(cfg.MongoDB.URI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.TODO())
	fmt.Println("Connected to MongoDB!")

	// Initialize Redis
	redisClient, err := database.NewRedisClient(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()
	fmt.Println("Connected to Redis!")

	// Initialize Repository
	userRepo := repository.NewMongoUserRepository(mongoClient.Database(cfg.MongoDB.Database))

	// Initialize Email Service
	emailService := service.NewEmailService(&cfg)
	if emailService.IsEnabled() {
		fmt.Printf("Email Service enabled (SMTP: %s:%d)\n", cfg.SMTP.Host, cfg.SMTP.Port)
	} else {
		fmt.Println("Email Service disabled (set smtp.enabled=true to enable)")
	}

	// Initialize Service
	authService := service.NewAuthService(userRepo, &cfg, emailService)

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
			"status":      "ok",
			"service":     "auth-service",
			"env":         cfg.AppEnv,
			"smtpEnabled": cfg.SMTP.Enabled,
		})
	})
	// Public routes
	v1.POST("/register", authHandler.Register)
	v1.POST("/login", authHandler.Login)
	v1.POST("/verify-email", authHandler.VerifyEmail)
	v1.POST("/resend-verification", authHandler.ResendVerification)
	v1.POST("/forgot-password", authHandler.ForgotPassword)
	v1.POST("/reset-password", authHandler.ResetPassword)
	// Protected routes
	v1.GET("/me", authHandler.GetProfile, middleware.AuthMiddleware(&cfg))
	v1.PUT("/me", authHandler.UpdateProfile, middleware.AuthMiddleware(&cfg))
	v1.POST("/change-password", authHandler.ChangePassword, middleware.AuthMiddleware(&cfg))

	// Start gRPC Server
	authGRPCService := grpc.NewAuthGRPCService(authService)
	grpcPort, err := strconv.Atoi(cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Invalid gRPC port: %v", err)
	}
	grpcServer := grpc.NewGRPCServer(authGRPCService, grpcPort)

	go func() {
		fmt.Printf("gRPC Server starting on port %s...\n", cfg.GRPCPort)
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
