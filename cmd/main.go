package main

import (
	"context"
	"log"
	"time"

	"github.com/anggitaputra13/mcs-auth/config"
	"github.com/anggitaputra13/mcs-auth/docs"
	"github.com/anggitaputra13/mcs-auth/internal/interfaces/api/handlers"
	"github.com/anggitaputra13/mcs-auth/internal/interfaces/api/routes"
	"github.com/anggitaputra13/mcs-auth/internal/interfaces/postgres"
	auth_usecase "github.com/anggitaputra13/mcs-auth/internal/usecases/auth"
	"github.com/anggitaputra13/mcs-auth/pkg/auth"
	"github.com/anggitaputra13/mcs-auth/pkg/database"
	"github.com/anggitaputra13/mcs-auth/pkg/redis"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/anggitaputra13/mcs-auth/internal/domain/entities"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @title Authentication Service API
// @version 1.0
// @description This is a microservice for user authentication.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /api/v1
// @schemes http
func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize PostgreSQL with retry logic
	var db *gorm.DB
	for i := 0; i < 5; i++ {
		db, err = database.NewPostgresDB(database.Config{
			Host:     cfg.DBHost,
			Port:     cfg.DBPort,
			User:     cfg.DBUser,
			Password: cfg.DBPassword,
			DBName:   cfg.DBName,
			SSLMode:  cfg.DBSSLMode,
		})
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database (attempt %d): %v", i+1, err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to database after retries: %v", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize Redis with retry logic
	var redisClient *redis.Client
	for i := 0; i < 5; i++ {
		redisClient = redis.NewRedisClient(cfg.RedisURL)
		if _, err := redisClient.Ping(context.Background()).Result(); err == nil {
			break
		}
		log.Printf("Failed to connect to Redis (attempt %d): %v", i+1, err)
		time.Sleep(5 * time.Second)
	}
	defer redisClient.Close()

	// Initialize JWT
	jwt := auth.NewJWT(cfg.JWTSecret, time.Duration(cfg.JWTExpiration)*time.Hour, redisClient)

	// Initialize repository
	userRepo := postgres.NewUserRepository(db)

	// Initialize use case
	authUseCase := auth_usecase.NewAuthUseCase(userRepo, jwt)

	// Initialize handler
	authHandler := handlers.NewAuthHandler(authUseCase, jwt)

	// Setup router
	router := gin.Default()

	// Programmatically set swagger info
	docs.SwaggerInfo.Title = "Auth Service API"
	docs.SwaggerInfo.Description = "Authentication Microservice"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = cfg.Host // Or your domain
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Add Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/swagger/doc.json"), // Exact path to your JSON
		ginSwagger.DefaultModelsExpandDepth(-1)))

	// Setup routes
	routes.SetupAuthRoutes(router, authHandler, jwt)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func runMigrations(db *gorm.DB) error {
	// AutoMigrate will create tables if they don't exist
	err := db.AutoMigrate(&entities.User{})
	if err != nil {
		return err
	}
	return nil
}
