package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/config"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/grpc/client"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/handler"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/middleware"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/repository"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/service"
	"github.com/geekzy/gspend-app/apps/financial-service/pkg/database"
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

	// Initialize gRPC Client for Auth Service
	authGRPCAddr := cfg.AuthServiceGRPCAddr
	if authGRPCAddr == "" {
		authGRPCAddr = "localhost:9091" // Fallback for local development
	}
	
	authClient, err := client.NewAuthGRPCClient(authGRPCAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Auth Service gRPC: %v", err)
	}
	defer authClient.Close()
	fmt.Printf("Connected to Auth Service gRPC at %s\n", authGRPCAddr)

	// Initialize Repositories
	db := mongoClient.Database(cfg.MongoDatabase)
	categoryRepo := repository.NewMongoCategoryRepository(db)
	incomeRepo := repository.NewMongoIncomeRepository(db)
	budgetRepo := repository.NewMongoBudgetRepository(db)
	transactionRepo := repository.NewMongoTransactionRepository(db)

	// Initialize Services
	categorySvc := service.NewCategoryService(categoryRepo)
	incomeSvc := service.NewIncomeService(incomeRepo)
	budgetSvc := service.NewBudgetService(budgetRepo)
	transactionSvc := service.NewTransactionService(transactionRepo, budgetRepo)

	// Initialize Handlers
	categoryHandler := handler.NewCategoryHandler(categorySvc)
	incomeHandler := handler.NewIncomeHandler(incomeSvc)
	budgetHandler := handler.NewBudgetHandler(budgetSvc)
	transactionHandler := handler.NewTransactionHandler(transactionSvc)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// Health Check
	e.GET("/api/v1/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "ok",
			"service": "financial-service",
			"env":     cfg.AppEnv,
		})
	})

	// Protected Routes
	api := e.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(authClient))

	// Categories
	categories := api.Group("/categories")
	categories.POST("", categoryHandler.Create)
	categories.GET("", categoryHandler.List)
	categories.PUT("/:id", categoryHandler.Update)
	categories.DELETE("/:id", categoryHandler.Delete)

	// Incomes
	incomes := api.Group("/incomes")
	incomes.POST("", incomeHandler.Create)
	incomes.GET("", incomeHandler.List)
	incomes.PUT("/:id", incomeHandler.Update)
	incomes.DELETE("/:id", incomeHandler.Delete)

	// Budgets
	budgets := api.Group("/budgets")
	budgets.POST("", budgetHandler.Create)
	budgets.GET("", budgetHandler.List)
	budgets.GET("/active", budgetHandler.GetActive)
	budgets.PUT("/:id", budgetHandler.Update)
	budgets.DELETE("/:id", budgetHandler.Delete)

	// Transactions
	transactions := api.Group("/transactions")
	transactions.POST("", transactionHandler.Create)
	transactions.GET("", transactionHandler.List)
	transactions.GET("/:id", transactionHandler.GetByID)
	transactions.PUT("/:id", transactionHandler.Update)
	transactions.DELETE("/:id", transactionHandler.Delete)

	// Start Server
	fmt.Printf("Financial Service starting on port %s...\n", cfg.Port)
	if err := e.Start(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
