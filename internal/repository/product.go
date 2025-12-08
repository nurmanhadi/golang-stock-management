package repository

import (
	"context"
	"stock-management/internal/entity"

	"gorm.io/gorm"
)

type Product interface {
	Create(db *gorm.DB, product *entity.Product) (int64, error)
	CountBySku(sku string) (int64, error)
	Find(page, size int) ([]entity.Product, error)
	Count() (int64, error)
}
type productRepo struct {
	db *gorm.DB
}

func NewProduct(db *gorm.DB) Product {
	return &productRepo{
		db: db,
	}
}

func (r *productRepo) Create(db *gorm.DB, product *entity.Product) (int64, error) {
	ctx := context.Background()
	err := gorm.G[entity.Product](db).Create(ctx, product)
	if err != nil {
		return 0, nil
	}
	return product.ID, nil
}
func (r *productRepo) CountBySku(sku string) (int64, error) {
	ctx := context.Background()
	total, err := gorm.G[entity.Product](r.db).Where("sku = ?", sku).Count(ctx, "id")
	if err != nil {
		return 0, err
	}
	return total, nil
}
func (r *productRepo) Find(page, size int) ([]entity.Product, error) {
	ctx := context.Background()
	products, err := gorm.G[entity.Product](r.db).
		Offset((page - 1) * size).
		Limit(size).
		Find(ctx)
	if err != nil {
		return nil, nil
	}
	return products, nil
}
func (r *productRepo) Count() (int64, error) {
	ctx := context.Background()
	total, err := gorm.G[entity.Product](r.db).Count(ctx, "id")
	if err != nil {
		return 0, err
	}
	return total, nil
}
