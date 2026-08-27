// +build wireinject

package di

import (
	"github.com/google/wire"
	"gorm.io/gorm"
	"project-test-primo/internal/db"
	"project-test-primo/internal/handler"
	"project-test-primo/internal/repository"
	"project-test-primo/internal/usecase"
)

func InitializeProductHandler(postgres *gorm.DB) *handler.ProductHandler {
	wire.Build(
		repository.NewProductRepository,
		usecase.NewProductUseCase,
		handler.NewProductHandler,
	)
	panic("this code is replaced by wire injection")
}

func InitializeDatabase() (*gorm.DB, error) {
	return db.NewPostgresDB()
}

func InitializeApp() (*gorm.DB, *handler.ProductHandler, error) {
	database, err := InitializeDatabase()
	if err != nil {
		return nil, nil, err
	}

	if err := db.InitSchema(database); err != nil {
		return nil, nil, err
	}

	productHandler := InitializeProductHandler(database)

	return database, productHandler, nil
}
