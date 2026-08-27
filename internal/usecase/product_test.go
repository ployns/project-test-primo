package usecase

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"project-test-primo/internal/domain"
	"project-test-primo/internal/repository/mocks"
)

func TestProductUseCase_CreateProduct(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mockRepo := mocks.NewProductRepository(t)
		expectedProduct := &domain.Product{
			ID:    1,
			Name:  "Test Product",
			Price: 99.99,
		}

		mockRepo.EXPECT().
			Create(mock.Anything, mock.MatchedBy(func(p *domain.Product) bool {
				return p.Name == "Test Product" && p.Price == 99.99
			})).
			Return(expectedProduct, nil).
			Once()

		uc := NewProductUseCase(mockRepo)
		product, err := uc.CreateProduct(context.Background(), &domain.CreateProductRequest{
			Name:  "Test Product",
			Price: 99.99,
		})

		require.NoError(t, err)
		assert.Equal(t, expectedProduct.ID, product.ID)
		assert.Equal(t, "Test Product", product.Name)
		assert.Equal(t, 99.99, product.Price)
	})

	t.Run("validation - empty name", func(t *testing.T) {
		mockRepo := mocks.NewProductRepository(t)

		uc := NewProductUseCase(mockRepo)
		_, err := uc.CreateProduct(context.Background(), &domain.CreateProductRequest{
			Name:  "",
			Price: 99.99,
		})

		require.Error(t, err)
		assert.Equal(t, ErrInvalidProductName, err)
	})

	t.Run("validation - zero price", func(t *testing.T) {
		mockRepo := mocks.NewProductRepository(t)

		uc := NewProductUseCase(mockRepo)
		_, err := uc.CreateProduct(context.Background(), &domain.CreateProductRequest{
			Name:  "Test",
			Price: 0,
		})

		require.Error(t, err)
		assert.Equal(t, ErrInvalidPrice, err)
	})

	t.Run("validation - negative price", func(t *testing.T) {
		mockRepo := mocks.NewProductRepository(t)

		uc := NewProductUseCase(mockRepo)
		_, err := uc.CreateProduct(context.Background(), &domain.CreateProductRequest{
			Name:  "Test",
			Price: -10.0,
		})

		require.Error(t, err)
		assert.Equal(t, ErrInvalidPrice, err)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := mocks.NewProductRepository(t)
		mockRepo.EXPECT().
			Create(mock.Anything, mock.Anything).
			Return(nil, sql.ErrConnDone).
			Once()

		uc := NewProductUseCase(mockRepo)
		_, err := uc.CreateProduct(context.Background(), &domain.CreateProductRequest{
			Name:  "Test",
			Price: 99.99,
		})

		require.Error(t, err)
		assert.Equal(t, sql.ErrConnDone, err)
	})
}

func TestProductUseCase_GetProduct(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		mockRepo := mocks.NewProductRepository(t)
		expectedProduct := &domain.Product{
			ID:    1,
			Name:  "Test Product",
			Price: 99.99,
		}

		mockRepo.EXPECT().
			GetByID(mock.Anything, int64(1)).
			Return(expectedProduct, nil).
			Once()

		uc := NewProductUseCase(mockRepo)
		product, err := uc.GetProduct(context.Background(), 1)

		require.NoError(t, err)
		assert.Equal(t, expectedProduct, product)
	})

	t.Run("validation - invalid id", func(t *testing.T) {
		mockRepo := mocks.NewProductRepository(t)

		uc := NewProductUseCase(mockRepo)
		_, err := uc.GetProduct(context.Background(), -1)

		require.Error(t, err)
		assert.Equal(t, ErrInvalidProductID, err)
	})

	t.Run("validation - zero id", func(t *testing.T) {
		mockRepo := mocks.NewProductRepository(t)

		uc := NewProductUseCase(mockRepo)
		_, err := uc.GetProduct(context.Background(), 0)

		require.Error(t, err)
		assert.Equal(t, ErrInvalidProductID, err)
	})

	t.Run("product not found", func(t *testing.T) {
		mockRepo := mocks.NewProductRepository(t)
		mockRepo.EXPECT().
			GetByID(mock.Anything, int64(9999)).
			Return(nil, sql.ErrNoRows).
			Once()

		uc := NewProductUseCase(mockRepo)
		_, err := uc.GetProduct(context.Background(), 9999)

		require.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
	})
}

func TestProductUseCase_UpdateProduct(t *testing.T) {
	t.Run("successful update with all fields", func(t *testing.T) {
		mockRepo := mocks.NewProductRepository(t)
		existingProduct := &domain.Product{
			ID:    1,
			Name:  "Original",
			Price: 99.99,
		}
		newPrice := 149.99

		mockRepo.EXPECT().
			GetByID(mock.Anything, int64(1)).
			Return(existingProduct, nil).
			Once()
		mockRepo.EXPECT().
			Update(mock.Anything, int64(1), mock.MatchedBy(func(p *domain.Product) bool {
				return p.Name == "Updated" && p.Price == newPrice
			})).
			Return(nil).
			Once()

		uc := NewProductUseCase(mockRepo)
		err := uc.UpdateProduct(context.Background(), 1, &domain.UpdateProductRequest{
			Name:  ptrString("Updated"),
			Price: &newPrice,
		})

		require.NoError(t, err)
	})

	t.Run("successful update with partial fields", func(t *testing.T) {
		mockRepo := mocks.NewProductRepository(t)
		existingProduct := &domain.Product{
			ID:    1,
			Name:  "Original",
			Price: 99.99,
		}
		newPrice := 149.99

		mockRepo.EXPECT().
			GetByID(mock.Anything, int64(1)).
			Return(existingProduct, nil).
			Once()
		mockRepo.EXPECT().
			Update(mock.Anything, int64(1), mock.Anything).
			Return(nil).
			Once()

		uc := NewProductUseCase(mockRepo)
		err := uc.UpdateProduct(context.Background(), 1, &domain.UpdateProductRequest{
			Price: &newPrice,
		})

		require.NoError(t, err)
	})

	t.Run("validation - invalid id", func(t *testing.T) {
		mockRepo := mocks.NewProductRepository(t)

		uc := NewProductUseCase(mockRepo)
		err := uc.UpdateProduct(context.Background(), -1, &domain.UpdateProductRequest{})

		require.Error(t, err)
		assert.Equal(t, ErrInvalidProductID, err)
	})

	t.Run("product not found", func(t *testing.T) {
		mockRepo := mocks.NewProductRepository(t)
		mockRepo.EXPECT().
			GetByID(mock.Anything, int64(9999)).
			Return(nil, sql.ErrNoRows).
			Once()

		uc := NewProductUseCase(mockRepo)
		err := uc.UpdateProduct(context.Background(), 9999, &domain.UpdateProductRequest{})

		require.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("validation - empty name on update", func(t *testing.T) {
		mockRepo := mocks.NewProductRepository(t)
		existingProduct := &domain.Product{
			ID:    1,
			Name:  "Original",
			Price: 99.99,
		}

		mockRepo.EXPECT().
			GetByID(mock.Anything, int64(1)).
			Return(existingProduct, nil).
			Once()

		uc := NewProductUseCase(mockRepo)
		err := uc.UpdateProduct(context.Background(), 1, &domain.UpdateProductRequest{
			Name: ptrString(""),
		})

		require.Error(t, err)
		assert.Equal(t, ErrInvalidProductName, err)
	})

	t.Run("validation - negative price on update", func(t *testing.T) {
		mockRepo := mocks.NewProductRepository(t)
		existingProduct := &domain.Product{
			ID:    1,
			Name:  "Original",
			Price: 99.99,
		}
		negPrice := -50.0

		mockRepo.EXPECT().
			GetByID(mock.Anything, int64(1)).
			Return(existingProduct, nil).
			Once()

		uc := NewProductUseCase(mockRepo)
		err := uc.UpdateProduct(context.Background(), 1, &domain.UpdateProductRequest{
			Price: &negPrice,
		})

		require.Error(t, err)
		assert.Equal(t, ErrInvalidPrice, err)
	})
}

func TestProductUseCase_DeleteProduct(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		mockRepo := mocks.NewProductRepository(t)
		mockRepo.EXPECT().
			Delete(mock.Anything, int64(1)).
			Return(nil).
			Once()

		uc := NewProductUseCase(mockRepo)
		err := uc.DeleteProduct(context.Background(), 1)

		require.NoError(t, err)
	})

	t.Run("validation - invalid id", func(t *testing.T) {
		mockRepo := mocks.NewProductRepository(t)

		uc := NewProductUseCase(mockRepo)
		err := uc.DeleteProduct(context.Background(), -1)

		require.Error(t, err)
		assert.Equal(t, ErrInvalidProductID, err)
	})

	t.Run("product not found", func(t *testing.T) {
		mockRepo := mocks.NewProductRepository(t)
		mockRepo.EXPECT().
			Delete(mock.Anything, int64(9999)).
			Return(sql.ErrNoRows).
			Once()

		uc := NewProductUseCase(mockRepo)
		err := uc.DeleteProduct(context.Background(), 9999)

		require.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
	})
}

// Helper function for pointer to string
func ptrString(s string) *string {
	return &s
}
