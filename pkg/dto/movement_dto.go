package dto

import (
	"stock-management/pkg/enum"
	"time"
)

// RESPONSE
type MovementReposne struct {
	ID        int64             `json:"id"`
	ProductID int64             `json:"product_id"`
	Type      enum.Type         `json:"type"`
	Qty       int               `json:"qty"`
	CreatedAt time.Time         `json:"created_at"`
	Movements []MovementReposne `json:"movements"`
}
