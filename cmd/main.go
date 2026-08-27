package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "project-test-primo/docs"
	"project-test-primo/internal/db"
	"project-test-primo/internal/handler"
	"project-test-primo/internal/repository"
	"project-test-primo/internal/usecase"
)

// @title Product API
// @version 1.0
// @description A simple Product API with clean architecture using GORM
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	godotenv.Load()

	database, err := db.NewPostgresDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.InitSchema(database); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	productRepo := repository.NewProductRepository(database)
	productUC := usecase.NewProductUseCase(productRepo)
	productHandler := handler.NewProductHandler(productUC)

	router := gin.Default()

	router.GET("/api-docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	productHandler.RegisterRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
