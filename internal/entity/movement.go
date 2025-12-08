package entity

import (
	"stock-management/pkg/enum"
	"time"
)

type Movement struct {
	ID        int64
	ProductID int64
	Type      enum.Type
	Qty       int
	CreatedAt time.Time
	Product   *Product
}
