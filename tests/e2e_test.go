package tests

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"project-test-primo/internal/domain"
	"project-test-primo/internal/handler"
	"project-test-primo/internal/repository/mocks"
	"project-test-primo/internal/usecase"
)

func setupE2ETest(t *testing.T) (*mocks.MockProductRepository, *gin.Engine) {
	mockRepo := mocks.NewProductRepository(t)
	productUC := usecase.NewProductUseCase(mockRepo)
	productHandler := handler.NewProductHandler(productUC)

	router := gin.Default()
	productHandler.RegisterRoutes(router)

	return mockRepo, router
}

func TestE2E_CreateAndRetrieveProduct(t *testing.T) {
	mockRepo, router := setupE2ETest(t)

	t.Run("create product and retrieve it", func(t *testing.T) {
		mockRepo.EXPECT().Create(context.Background(), mock.MatchedBy(func(p *domain.Product) bool {
			return p.Name == "Laptop" && p.Price == 1299.99
		})).Return(&domain.Product{
			ID:    1,
			Name:  "Laptop",
			Price: 1299.99,
		}, nil)

		mockRepo.EXPECT().GetByID(context.Background(), int64(1)).Return(&domain.Product{
			ID:    1,
			Name:  "Laptop",
			Price: 1299.99,
		}, nil)
		createBody := []byte(`{
			"name": "Laptop",
			"description": "A powerful laptop",
			"price": 1299.99,
			"sale_price": 999.99
		}`)

		req := httptest.NewRequest("POST", "/product", bytes.NewBuffer(createBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var createResp domain.Response
		err := json.Unmarshal(w.Body.Bytes(), &createResp)
		require.NoError(t, err)
		assert.True(t, createResp.Successful)

		data := createResp.Data.(map[string]interface{})
		productID := int64(data["data1"].(float64))

		getReq := httptest.NewRequest("GET", "/product/1", nil)
		getW := httptest.NewRecorder()
		router.ServeHTTP(getW, getReq)

		assert.Equal(t, http.StatusOK, getW.Code)

		var getResp domain.Response
		err = json.Unmarshal(getW.Body.Bytes(), &getResp)
		require.NoError(t, err)
		assert.True(t, getResp.Successful)
		assert.NotZero(t, productID)
	})
}

func TestE2E_CreateUpdateDelete(t *testing.T) {
	mockRepo, router := setupE2ETest(t)

	t.Run("full CRUD workflow", func(t *testing.T) {
		mockRepo.EXPECT().Create(context.Background(), mock.MatchedBy(func(p *domain.Product) bool {
			return p.Name == "Mouse" && p.Price == 49.99
		})).Return(&domain.Product{
			ID:    1,
			Name:  "Mouse",
			Price: 49.99,
		}, nil)

		mockRepo.EXPECT().Update(context.Background(), int64(1), mock.Anything).Return(nil)

		mockRepo.EXPECT().GetByID(context.Background(), int64(1)).Return(&domain.Product{
			ID:    1,
			Name:  "Wireless Mouse",
			Price: 49.99,
		}, nil)

		mockRepo.EXPECT().Delete(context.Background(), int64(1)).Return(nil)
		createBody := []byte(`{
			"name": "Mouse",
			"price": 49.99
		}`)

		createReq := httptest.NewRequest("POST", "/product", bytes.NewBuffer(createBody))
		createReq.Header.Set("Content-Type", "application/json")
		createW := httptest.NewRecorder()
		router.ServeHTTP(createW, createReq)

		assert.Equal(t, http.StatusOK, createW.Code)

		var createResp domain.Response
		json.Unmarshal(createW.Body.Bytes(), &createResp)
		assert.True(t, createResp.Successful)

		updateBody := []byte(`{
			"name": "Wireless Mouse",
			"sale_price": 39.99
		}`)

		updateReq := httptest.NewRequest("PATCH", "/product/1", bytes.NewBuffer(updateBody))
		updateReq.Header.Set("Content-Type", "application/json")
		updateW := httptest.NewRecorder()
		router.ServeHTTP(updateW, updateReq)

		assert.Equal(t, http.StatusOK, updateW.Code)

		var updateResp domain.Response
		json.Unmarshal(updateW.Body.Bytes(), &updateResp)
		assert.True(t, updateResp.Successful)

		getReq := httptest.NewRequest("GET", "/product/1", nil)
		getW := httptest.NewRecorder()
		router.ServeHTTP(getW, getReq)

		var getResp domain.Response
		json.Unmarshal(getW.Body.Bytes(), &getResp)
		assert.True(t, getResp.Successful)

		deleteReq := httptest.NewRequest("DELETE", "/product/1", nil)
		deleteW := httptest.NewRecorder()
		router.ServeHTTP(deleteW, deleteReq)

		assert.Equal(t, http.StatusOK, deleteW.Code)

		var deleteResp domain.Response
		json.Unmarshal(deleteW.Body.Bytes(), &deleteResp)
		assert.True(t, deleteResp.Successful)

		notFoundReq := httptest.NewRequest("GET", "/product/1", nil)
		notFoundW := httptest.NewRecorder()
		router.ServeHTTP(notFoundW, notFoundReq)

		assert.Equal(t, http.StatusOK, notFoundW.Code)
	})
}

func TestE2E_ValidationScenarios(t *testing.T) {
	mockRepo, router := setupE2ETest(t)

	t.Run("create product without name should fail", func(t *testing.T) {
		body := []byte(`{"price": 99.99}`)
		req := httptest.NewRequest("POST", "/product", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp domain.Response
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Successful)
	})

	t.Run("create product with zero price should fail", func(t *testing.T) {
		body := []byte(`{"name": "Product", "price": 0}`)
		req := httptest.NewRequest("POST", "/product", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp domain.Response
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Successful)
	})

	t.Run("get non-existent product should return not found", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(context.Background(), int64(9999)).Return(nil, sql.ErrNoRows)

		req := httptest.NewRequest("GET", "/product/9999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var resp domain.Response
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Successful)
		assert.Equal(t, "NOT_FOUND", resp.ErrorCode)
	})
}

func TestE2E_PartialUpdate(t *testing.T) {
	mockRepo, router := setupE2ETest(t)

	t.Run("update only specific fields", func(t *testing.T) {
		mockRepo.EXPECT().Create(context.Background(), mock.MatchedBy(func(p *domain.Product) bool {
			return p.Name == "Keyboard" && p.Price == 149.99
		})).Return(&domain.Product{
			ID:    1,
			Name:  "Keyboard",
			Price: 149.99,
		}, nil)

		mockRepo.EXPECT().Update(context.Background(), int64(1), mock.Anything).Return(nil)

		mockRepo.EXPECT().GetByID(context.Background(), int64(1)).Return(&domain.Product{
			ID:    1,
			Name:  "Keyboard",
			Price: 199.99,
		}, nil)
		createBody := []byte(`{
			"name": "Keyboard",
			"description": "Mechanical keyboard",
			"price": 149.99,
			"sale_price": 119.99
		}`)

		createReq := httptest.NewRequest("POST", "/product", bytes.NewBuffer(createBody))
		createReq.Header.Set("Content-Type", "application/json")
		createW := httptest.NewRecorder()
		router.ServeHTTP(createW, createReq)

		assert.Equal(t, http.StatusOK, createW.Code)

		updateBody := []byte(`{"price": 199.99}`)

		updateReq := httptest.NewRequest("PATCH", "/product/1", bytes.NewBuffer(updateBody))
		updateReq.Header.Set("Content-Type", "application/json")
		updateW := httptest.NewRecorder()
		router.ServeHTTP(updateW, updateReq)

		assert.Equal(t, http.StatusOK, updateW.Code)

		var updateResp domain.Response
		json.Unmarshal(updateW.Body.Bytes(), &updateResp)
		assert.True(t, updateResp.Successful)

		getReq := httptest.NewRequest("GET", "/product/1", nil)
		getW := httptest.NewRecorder()
		router.ServeHTTP(getW, getReq)

		var getResp domain.Response
		json.Unmarshal(getW.Body.Bytes(), &getResp)

		data := getResp.Data.(map[string]interface{})
		assert.Equal(t, "Keyboard", data["data2"])
		assert.Equal(t, 199.99, data["price"])
	})
}

func TestE2E_UpdateZeroPrice(t *testing.T) {
	mockRepo, router := setupE2ETest(t)

	t.Run("update price to zero should fail validation", func(t *testing.T) {
		mockRepo.EXPECT().Create(context.Background(), mock.MatchedBy(func(p *domain.Product) bool {
			return p.Name == "Product" && p.Price == 50.0
		})).Return(&domain.Product{
			ID:    1,
			Name:  "Product",
			Price: 50.0,
		}, nil)

		mockRepo.EXPECT().GetByID(context.Background(), int64(1)).Return(&domain.Product{
			ID:    1,
			Name:  "Product",
			Price: 50.0,
		}, nil)

		createBody := []byte(`{"name": "Product", "price": 50}`)
		createReq := httptest.NewRequest("POST", "/product", bytes.NewBuffer(createBody))
		createReq.Header.Set("Content-Type", "application/json")
		createW := httptest.NewRecorder()
		router.ServeHTTP(createW, createReq)

		assert.Equal(t, http.StatusOK, createW.Code)

		updateBody := []byte(`{"price": 0}`)
		updateReq := httptest.NewRequest("PATCH", "/product/1", bytes.NewBuffer(updateBody))
		updateReq.Header.Set("Content-Type", "application/json")
		updateW := httptest.NewRecorder()
		router.ServeHTTP(updateW, updateReq)

		assert.Equal(t, http.StatusBadRequest, updateW.Code)

		var updateResp domain.Response
		json.Unmarshal(updateW.Body.Bytes(), &updateResp)
		assert.False(t, updateResp.Successful)
	})
}
