package repository

import (
	"context"
	"stock-management/internal/entity"

	"gorm.io/gorm"
)

type Movement interface {
	Create(db *gorm.DB, movement *entity.Movement) error
	Find(page, size int) ([]entity.Movement, error)
	Count() (int64, error)
}
type movementRepo struct {
	db *gorm.DB
}

func NewMovement(db *gorm.DB) Movement {
	return &movementRepo{db: db}
}
func (r *movementRepo) Create(db *gorm.DB, movement *entity.Movement) error {
	ctx := context.Background()
	err := gorm.G[entity.Movement](db).Create(ctx, movement)
	if err != nil {
		return nil
	}
	return nil
}
func (r *movementRepo) Find(page, size int) ([]entity.Movement, error) {
	ctx := context.Background()
	products, err := gorm.G[entity.Movement](r.db).
		Offset((page-1)*size).
		Limit(size).
		Preload("Product", nil).
		Find(ctx)
	if err != nil {
		return nil, err
	}
	return products, nil
}
func (r *movementRepo) Count() (int64, error) {
	ctx := context.Background()
	total, err := gorm.G[entity.Movement](r.db).Count(ctx, "id")
	if err != nil {
		return 0, err
	}
	return total, nil
}
