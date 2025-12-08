package dto

import "time"

// REQUEST
type ProductAddRequest struct {
	Sku   string  `json:"sku" validate:"required,max=20"`
	Name  string  `json:"name" validate:"required,max=100"`
	Price float64 `json:"price" validate:"required"`
	Stock int     `json:"stock" validate:"required"`
}

// RESPONSE
type ProductResposne struct {
	ID        int64     `json:"id"`
	Sku       string    `json:"sku"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
