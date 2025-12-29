package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/config"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/database"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/grpc/client"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/handler"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/middleware"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/repository"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/service"
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
		log.Printf("Warning: Failed to connect to Auth Service gRPC: %v", err)
		log.Printf("Service will start but authentication will not work until gRPC connection is established")
		authClient = nil
	} else {
		defer authClient.Close()
		fmt.Printf("Connected to Auth Service gRPC at %s\n", authGRPCAddr)
	}

	// Initialize Repositories
	db := mongoClient.Database(cfg.MongoDatabase)
	categoryRepo := repository.NewMongoCategoryRepository(db)
	incomeRepo := repository.NewMongoIncomeRepository(db)
	budgetRepo := repository.NewMongoBudgetRepository(db)
	transactionRepo := repository.NewMongoTransactionRepository(db)
	dashboardRepo := repository.NewMongoDashboardRepository(db)
	reportRepo := repository.NewMongoReportRepository(db)

	// Initialize Services
	categorySvc := service.NewCategoryService(categoryRepo)
	incomeSvc := service.NewIncomeService(incomeRepo)
	budgetSvc := service.NewBudgetService(budgetRepo)
	transactionSvc := service.NewTransactionService(transactionRepo, budgetRepo)
	dashboardSvc := service.NewDashboardService(dashboardRepo)
	reportSvc := service.NewReportService(reportRepo, transactionRepo, budgetRepo)

	// Initialize system categories on startup
	fmt.Println("Initializing system categories...")
	if err := categorySvc.InitializeSystemCategories(context.Background()); err != nil {
		log.Printf("Warning: Failed to initialize system categories: %v", err)
	} else {
		fmt.Println("✓ System categories initialized successfully")
	}

	// Initialize Handlers
	categoryHandler := handler.NewCategoryHandler(categorySvc)
	incomeHandler := handler.NewIncomeHandler(incomeSvc)
	budgetHandler := handler.NewBudgetHandler(budgetSvc)
	transactionHandler := handler.NewTransactionHandler(transactionSvc)
	dashboardHandler := handler.NewDashboardHandler(dashboardSvc)
	reportHandler := handler.NewReportHandler(reportSvc)

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
	if authClient != nil {
		api.Use(middleware.AuthMiddleware(authClient))
	} else {
		log.Println("Warning: Running WITHOUT real authentication (gRPC client not available)")
		log.Println("Using demo fallback middleware with hardcoded user ID")
		// Fallback middleware for demo mode - sets demo user ID
		api.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				// Set demo user ID in context
				c.Set("user_id", "695284b2d79e48201abebde7")
				return next(c)
			}
		})
	}

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
	
	// Budget Items
	budgets.POST("/:id/items", budgetHandler.AddBudgetItem)
	budgets.GET("/:id/items/:itemId", budgetHandler.GetBudgetItem)
	budgets.PUT("/:id/items/:itemId", budgetHandler.UpdateBudgetItem)
	budgets.DELETE("/:id/items/:itemId", budgetHandler.DeleteBudgetItem)

	// Transactions
	transactions := api.Group("/transactions")
	transactions.POST("", transactionHandler.Create)
	transactions.GET("", transactionHandler.List)
	transactions.GET("/filtered", transactionHandler.ListWithFilters)
	transactions.GET("/spending-by-category", transactionHandler.GetSpendingByCategory)
	transactions.GET("/monthly-trends", transactionHandler.GetMonthlyTrends)
	transactions.GET("/:id", transactionHandler.GetByID)
	transactions.PUT("/:id", transactionHandler.Update)
	transactions.DELETE("/:id", transactionHandler.Delete)

	// Dashboard
	dashboard := api.Group("/dashboard")
	dashboard.GET("/summary", dashboardHandler.GetSummary)
	dashboard.GET("/recent-transactions", dashboardHandler.GetRecentTransactions)
	dashboard.GET("/top-categories", dashboardHandler.GetTopCategories)
	dashboard.GET("/budget-progress", dashboardHandler.GetBudgetProgress)

	// Reports
	reports := api.Group("/reports")
	reports.GET("/budget-vs-actual", reportHandler.GetBudgetVsActual)
	reports.GET("/spending-by-category", reportHandler.GetSpendingByCategory)
	reports.GET("/monthly-trends", reportHandler.GetMonthlyTrends)

	// Start Server
	fmt.Printf("Financial Service starting on port %s...\n", cfg.Port)
	if err := e.Start(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
