package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"project-test-primo/internal/domain"
	"project-test-primo/internal/usecase"
)

type ProductHandler struct {
	uc usecase.ProductUseCase
}

func NewProductHandler(uc usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{uc: uc}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the input payload
// @Tags products
// @Accept json
// @Produce json
// @Param input body domain.CreateProductRequest true "Create product request"
// @Success 200 {object} domain.Response
// @Failure 400 {object} domain.Response
// @Router /product [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req domain.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{
			Successful: false,
			ErrorCode:  "INVALID_REQUEST",
			Data:       nil,
		})
		return
	}

	product, err := h.uc.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{
			Successful: false,
			ErrorCode:  "CREATE_FAILED",
			Data:       nil,
		})
		return
	}

	c.JSON(http.StatusOK, domain.Response{
		Successful: true,
		ErrorCode:  "",
		Data: domain.ProductData{
			ID:        product.ID,
			Name:      product.Name,
			Price:     product.Price,
			SalePrice: product.SalePrice,
		},
	})
}

// GetProduct godoc
// @Summary Get a product by ID
// @Description Get a product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} domain.Response
// @Failure 404 {object} domain.Response
// @Router /product/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{
			Successful: false,
			ErrorCode:  "INVALID_ID",
			Data:       nil,
		})
		return
	}

	product, err := h.uc.GetProduct(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, domain.Response{
				Successful: false,
				ErrorCode:  "NOT_FOUND",
				Data:       nil,
			})
		} else {
			c.JSON(http.StatusInternalServerError, domain.Response{
				Successful: false,
				ErrorCode:  "INTERNAL_ERROR",
				Data:       nil,
			})
		}
		return
	}

	c.JSON(http.StatusOK, domain.Response{
		Successful: true,
		ErrorCode:  "",
		Data: domain.ProductData{
			ID:        product.ID,
			Name:      product.Name,
			Price:     product.Price,
			SalePrice: product.SalePrice,
		},
	})
}

// UpdateProduct godoc
// @Summary Update a product by ID
// @Description Update a product with partial fields
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param input body object true "Update product request (supports: name, description, price, sale_price)"
// @Success 200 {object} domain.Response
// @Failure 400 {object} domain.Response
// @Router /product/{id} [patch]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{
			Successful: false,
			ErrorCode:  "INVALID_ID",
			Data:       nil,
		})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{
			Successful: false,
			ErrorCode:  "INVALID_REQUEST",
			Data:       nil,
		})
		return
	}

	err = h.uc.UpdateProduct(c.Request.Context(), id, req)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, domain.Response{
				Successful: false,
				ErrorCode:  "NOT_FOUND",
				Data:       nil,
			})
		} else {
			c.JSON(http.StatusBadRequest, domain.Response{
				Successful: false,
				ErrorCode:  "UPDATE_FAILED",
				Data:       nil,
			})
		}
		return
	}

	c.JSON(http.StatusOK, domain.Response{
		Successful: true,
		ErrorCode:  "",
		Data:       nil,
	})
}

// DeleteProduct godoc
// @Summary Delete a product by ID
// @Description Delete a product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} domain.Response
// @Failure 404 {object} domain.Response
// @Router /product/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{
			Successful: false,
			ErrorCode:  "INVALID_ID",
			Data:       nil,
		})
		return
	}

	err = h.uc.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, domain.Response{
				Successful: false,
				ErrorCode:  "NOT_FOUND",
				Data:       nil,
			})
		} else {
			c.JSON(http.StatusInternalServerError, domain.Response{
				Successful: false,
				ErrorCode:  "DELETE_FAILED",
				Data:       nil,
			})
		}
		return
	}

	c.JSON(http.StatusOK, domain.Response{
		Successful: true,
		ErrorCode:  "",
		Data:       nil,
	})
}

func (h *ProductHandler) RegisterRoutes(r *gin.Engine) {
	products := r.Group("/")
	{
		products.POST("product", h.CreateProduct)
		products.GET("product/:id", h.GetProduct)
		products.PATCH("product/:id", h.UpdateProduct)
		products.DELETE("product/:id", h.DeleteProduct)
	}
}
