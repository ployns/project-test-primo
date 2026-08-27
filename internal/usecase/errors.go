package usecase

import "errors"

var (
	ErrInvalidProductID   = errors.New("invalid product id")
	ErrInvalidProductName = errors.New("product name cannot be empty")
	ErrInvalidPrice       = errors.New("price must be greater than 0")
	ErrProductNotFound    = errors.New("product not found")
)
