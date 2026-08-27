package usecase

import (
	"context"
	"project-test-primo/internal/domain"
	"project-test-primo/internal/repository"
)

type ProductUseCase interface {
	CreateProduct(ctx context.Context, req *domain.CreateProductRequest) (*domain.Product, error)
	GetProduct(ctx context.Context, id int64) (*domain.Product, error)
	UpdateProduct(ctx context.Context, id int64, updates map[string]interface{}) error
	DeleteProduct(ctx context.Context, id int64) error
}

type productUseCase struct {
	repo repository.ProductRepository
}

func NewProductUseCase(repo repository.ProductRepository) ProductUseCase {
	return &productUseCase{repo: repo}
}

func (uc *productUseCase) CreateProduct(ctx context.Context, req *domain.CreateProductRequest) (*domain.Product, error) {
	if req.Name == "" {
		return nil, ErrInvalidProductName
	}
	if req.Price <= 0 {
		return nil, ErrInvalidPrice
	}

	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		SalePrice:   req.SalePrice,
	}

	return uc.repo.Create(ctx, product)
}

func (uc *productUseCase) GetProduct(ctx context.Context, id int64) (*domain.Product, error) {
	if id <= 0 {
		return nil, ErrInvalidProductID
	}
	return uc.repo.GetByID(ctx, id)
}

func (uc *productUseCase) UpdateProduct(ctx context.Context, id int64, updates map[string]interface{}) error {
	if id <= 0 {
		return ErrInvalidProductID
	}

	existing, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if name, ok := updates["name"]; ok {
		nameStr, ok := name.(string)
		if !ok || nameStr == "" {
			return ErrInvalidProductName
		}
		existing.Name = nameStr
	}

	if description, ok := updates["description"]; ok {
		if description != nil {
			descStr, ok := description.(string)
			if ok {
				existing.Description = &descStr
			}
		}
	}

	if price, ok := updates["price"]; ok {
		priceFloat, ok := price.(float64)
		if !ok || priceFloat <= 0 {
			return ErrInvalidPrice
		}
		existing.Price = priceFloat
	}

	if salePrice, ok := updates["sale_price"]; ok {
		if salePrice != nil {
			salePriceFloat, ok := salePrice.(float64)
			if ok {
				existing.SalePrice = &salePriceFloat
			}
		}
	}

	return uc.repo.Update(ctx, id, existing)
}

func (uc *productUseCase) DeleteProduct(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidProductID
	}
	return uc.repo.Delete(ctx, id)
}
