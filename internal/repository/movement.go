package repository

import (
	"context"
	"stock-management/internal/entity"

	"gorm.io/gorm"
)

type Movement interface {
	Create(db *gorm.DB, movement *entity.Movement) error
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
