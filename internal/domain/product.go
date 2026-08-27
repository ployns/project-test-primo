package domain

import "time"

type Product struct {
	ID          int64      `gorm:"primaryKey" json:"id"`
	Name        string     `gorm:"column:name;type:varchar(255)" json:"name"`
	Description *string    `gorm:"column:description;type:text" json:"description"`
	Price       float64    `gorm:"column:price;type:double precision" json:"price"`
	SalePrice   *float64   `gorm:"column:sale_price;type:double precision" json:"sale_price"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Product) TableName() string {
	return "products"
}

type CreateProductRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description *string  `json:"description"`
	Price       float64  `json:"price" binding:"required,gt=0"`
	SalePrice   *float64 `json:"sale_price"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	SalePrice   *float64 `json:"sale_price"`
}

type Response struct {
	Successful bool        `json:"successful"`
	ErrorCode  string      `json:"error_code"`
	Data       interface{} `json:"data"`
}

type ProductData struct {
	ID        int64   `json:"data1"`
	Name      string  `json:"data2"`
	Price     float64 `json:"price,omitempty"`
	SalePrice *float64 `json:"sale_price,omitempty"`
}
