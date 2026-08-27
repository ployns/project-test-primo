package repository

import (
	"context"

	"gorm.io/gorm"

	"project-test-primo/internal/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) (*domain.Product, error)
	GetByID(ctx context.Context, id int64) (*domain.Product, error)
	Update(ctx context.Context, id int64, product *domain.Product) error
	Delete(ctx context.Context, id int64) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	result := r.db.WithContext(ctx).Create(product)
	return product, result.Error
}

func (r *productRepository) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	product := &domain.Product{}
	result := r.db.WithContext(ctx).First(product, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (r *productRepository) Update(ctx context.Context, id int64, product *domain.Product) error {
	result := r.db.WithContext(ctx).Model(&domain.Product{ID: id}).Updates(product)
	return result.Error
}

func (r *productRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&domain.Product{}, id)
	return result.Error
}
