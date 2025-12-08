package dto

type StockRequest struct {
	ProductID int64 `json:"product_id" validate:"required"`
	Qty       int   `json:"qty" validate:"required"`
}
