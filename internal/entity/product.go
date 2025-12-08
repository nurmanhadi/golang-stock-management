package entity

import "time"

type Product struct {
	ID        int64
	Sku       string
	Name      string
	Price     float64
	Stock     int
	CreatedAt time.Time
	UpdatedAt time.Time
	Movements []Movement
}
