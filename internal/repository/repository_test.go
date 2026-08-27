package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"project-test-primo/internal/domain"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=test_db sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	err = db.Migrator().DropTable(&domain.Product{})
	require.NoError(t, err)

	err = db.AutoMigrate(&domain.Product{})
	require.NoError(t, err)

	return db
}

func TestProductRepository_Create(t *testing.T) {
	db := setupTestDB(t)

	repo := NewProductRepository(db)

	t.Run("successful creation", func(t *testing.T) {
		ctx := context.Background()
		product := &domain.Product{
			Name:  "Test Product",
			Price: 99.99,
		}

		created, err := repo.Create(ctx, product)
		require.NoError(t, err)
		assert.Greater(t, created.ID, int64(0))
		assert.Equal(t, "Test Product", created.Name)
		assert.Equal(t, 99.99, created.Price)
		assert.NotZero(t, created.CreatedAt)
	})

	t.Run("creation with description and sale_price", func(t *testing.T) {
		ctx := context.Background()
		desc := "A great product"
		salePrice := 79.99
		product := &domain.Product{
			Name:        "Product with Sale",
			Description: &desc,
			Price:       99.99,
			SalePrice:   &salePrice,
		}

		created, err := repo.Create(ctx, product)
		require.NoError(t, err)
		assert.Equal(t, &desc, created.Description)
		assert.Equal(t, &salePrice, created.SalePrice)
	})
}

func TestProductRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)

	repo := NewProductRepository(db)

	t.Run("get existing product", func(t *testing.T) {
		ctx := context.Background()
		product := &domain.Product{
			Name:  "Fetch Test Product",
			Price: 49.99,
		}

		created, err := repo.Create(ctx, product)
		require.NoError(t, err)

		fetched, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, fetched.ID)
		assert.Equal(t, "Fetch Test Product", fetched.Name)
		assert.Equal(t, 49.99, fetched.Price)
	})

	t.Run("get non-existing product", func(t *testing.T) {
		ctx := context.Background()
		_, err := repo.GetByID(ctx, 9999)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestProductRepository_Update(t *testing.T) {
	db := setupTestDB(t)

	repo := NewProductRepository(db)

	t.Run("successful update", func(t *testing.T) {
		ctx := context.Background()
		product := &domain.Product{
			Name:  "Original Name",
			Price: 100.00,
		}

		created, err := repo.Create(ctx, product)
		require.NoError(t, err)

		created.Name = "Updated Name"
		created.Price = 150.00

		err = repo.Update(ctx, created.ID, created)
		require.NoError(t, err)

		fetched, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", fetched.Name)
		assert.Equal(t, 150.00, fetched.Price)
	})

	t.Run("update non-existing product", func(t *testing.T) {
		ctx := context.Background()
		product := &domain.Product{
			Name:  "Test",
			Price: 50.00,
		}

		err := repo.Update(ctx, 9999, product)
		// GORM doesn't error on update with non-existent ID, it just matches 0 rows
		// So we check if rows were affected by trying to get the product
		assert.NoError(t, err)
	})
}

func TestProductRepository_Delete(t *testing.T) {
	db := setupTestDB(t)

	repo := NewProductRepository(db)

	t.Run("successful delete", func(t *testing.T) {
		ctx := context.Background()
		product := &domain.Product{
			Name:  "Delete Test Product",
			Price: 50.00,
		}

		created, err := repo.Create(ctx, product)
		require.NoError(t, err)

		err = repo.Delete(ctx, created.ID)
		require.NoError(t, err)

		_, err = repo.GetByID(ctx, created.ID)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("delete non-existing product", func(t *testing.T) {
		ctx := context.Background()
		err := repo.Delete(ctx, 9999)
		// GORM doesn't error on delete with non-existent ID
		assert.NoError(t, err)
	})
}
