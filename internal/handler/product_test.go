package handler

import (
	"bytes"
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
	"project-test-primo/internal/usecase/mocks"
)

func TestProductHandler_CreateProduct(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mockUC := mocks.NewProductUseCase(t)
		expectedProduct := &domain.Product{
			ID:    1,
			Name:  "Test Product",
			Price: 99.99,
		}

		mockUC.EXPECT().
			CreateProduct(mock.Anything, mock.MatchedBy(func(r *domain.CreateProductRequest) bool {
				return r.Name == "Test Product" && r.Price == 99.99
			})).
			Return(expectedProduct, nil).
			Once()

		router := gin.New()
		handler := NewProductHandler(mockUC)
		handler.RegisterRoutes(router)

		body := []byte(`{"name": "Test Product", "price": 99.99}`)
		req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp domain.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.True(t, resp.Successful)
	})

	t.Run("invalid request body - missing name", func(t *testing.T) {
		mockUC := mocks.NewProductUseCase(t)
		router := gin.New()
		handler := NewProductHandler(mockUC)
		handler.RegisterRoutes(router)

		body := []byte(`{"price": 99.99}`)
		req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp domain.Response
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Successful)
		assert.Equal(t, "INVALID_REQUEST", resp.ErrorCode)
	})

	t.Run("invalid request body - missing price", func(t *testing.T) {
		mockUC := mocks.NewProductUseCase(t)
		router := gin.New()
		handler := NewProductHandler(mockUC)
		handler.RegisterRoutes(router)

		body := []byte(`{"name": "Test"}`)
		req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp domain.Response
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Successful)
	})
}

func TestProductHandler_GetProduct(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		mockUC := mocks.NewProductUseCase(t)
		expectedProduct := &domain.Product{
			ID:    1,
			Name:  "Test Product",
			Price: 99.99,
		}

		mockUC.EXPECT().
			GetProduct(mock.Anything, int64(1)).
			Return(expectedProduct, nil).
			Once()

		router := gin.New()
		handler := NewProductHandler(mockUC)
		handler.RegisterRoutes(router)

		req, _ := http.NewRequest("GET", "/product/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp domain.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.True(t, resp.Successful)
	})

	t.Run("product not found", func(t *testing.T) {
		mockUC := mocks.NewProductUseCase(t)
		mockUC.EXPECT().
			GetProduct(mock.Anything, int64(9999)).
			Return(nil, sql.ErrNoRows).
			Once()

		router := gin.New()
		handler := NewProductHandler(mockUC)
		handler.RegisterRoutes(router)

		req, _ := http.NewRequest("GET", "/product/9999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var resp domain.Response
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Successful)
		assert.Equal(t, "NOT_FOUND", resp.ErrorCode)
	})

	t.Run("invalid id format", func(t *testing.T) {
		mockUC := mocks.NewProductUseCase(t)
		router := gin.New()
		handler := NewProductHandler(mockUC)
		handler.RegisterRoutes(router)

		req, _ := http.NewRequest("GET", "/product/invalid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp domain.Response
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Successful)
		assert.Equal(t, "INVALID_ID", resp.ErrorCode)
	})
}

func TestProductHandler_UpdateProduct(t *testing.T) {
	t.Run("successful patch", func(t *testing.T) {
		mockUC := mocks.NewProductUseCase(t)
		mockUC.EXPECT().
			UpdateProduct(mock.Anything, int64(1), mock.Anything).
			Return(nil).
			Once()

		router := gin.New()
		handler := NewProductHandler(mockUC)
		handler.RegisterRoutes(router)

		body := []byte(`{"name": "Updated Product", "price": 149.99}`)
		req, _ := http.NewRequest("PATCH", "/product/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp domain.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.True(t, resp.Successful)
	})

	t.Run("product not found for update", func(t *testing.T) {
		mockUC := mocks.NewProductUseCase(t)
		mockUC.EXPECT().
			UpdateProduct(mock.Anything, int64(9999), mock.Anything).
			Return(sql.ErrNoRows).
			Once()

		router := gin.New()
		handler := NewProductHandler(mockUC)
		handler.RegisterRoutes(router)

		body := []byte(`{"price": 149.99}`)
		req, _ := http.NewRequest("PATCH", "/product/9999", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var resp domain.Response
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Successful)
		assert.Equal(t, "NOT_FOUND", resp.ErrorCode)
	})

	t.Run("invalid id format", func(t *testing.T) {
		mockUC := mocks.NewProductUseCase(t)
		router := gin.New()
		handler := NewProductHandler(mockUC)
		handler.RegisterRoutes(router)

		body := []byte(`{"name": "Updated"}`)
		req, _ := http.NewRequest("PATCH", "/product/invalid", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp domain.Response
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Successful)
		assert.Equal(t, "INVALID_ID", resp.ErrorCode)
	})
}

func TestProductHandler_DeleteProduct(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		mockUC := mocks.NewProductUseCase(t)
		mockUC.EXPECT().
			DeleteProduct(mock.Anything, int64(1)).
			Return(nil).
			Once()

		router := gin.New()
		handler := NewProductHandler(mockUC)
		handler.RegisterRoutes(router)

		req, _ := http.NewRequest("DELETE", "/product/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp domain.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.True(t, resp.Successful)
	})

	t.Run("product not found for delete", func(t *testing.T) {
		mockUC := mocks.NewProductUseCase(t)
		mockUC.EXPECT().
			DeleteProduct(mock.Anything, int64(9999)).
			Return(sql.ErrNoRows).
			Once()

		router := gin.New()
		handler := NewProductHandler(mockUC)
		handler.RegisterRoutes(router)

		req, _ := http.NewRequest("DELETE", "/product/9999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var resp domain.Response
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Successful)
		assert.Equal(t, "NOT_FOUND", resp.ErrorCode)
	})

	t.Run("invalid id format", func(t *testing.T) {
		mockUC := mocks.NewProductUseCase(t)
		router := gin.New()
		handler := NewProductHandler(mockUC)
		handler.RegisterRoutes(router)

		req, _ := http.NewRequest("DELETE", "/product/invalid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp domain.Response
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Successful)
		assert.Equal(t, "INVALID_ID", resp.ErrorCode)
	})
}
