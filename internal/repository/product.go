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
	FindByID(id int64) (entity.Product, error)
	UpdateStock(db *gorm.DB, productID int64, stock int) error
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
		return 0, err
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
		return nil, err
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
func (r *productRepo) FindByID(id int64) (entity.Product, error) {
	ctx := context.Background()
	product, err := gorm.G[entity.Product](r.db).First(ctx)
	if err != nil {
		return entity.Product{}, err
	}
	return product, nil
}
func (r *productRepo) UpdateStock(db *gorm.DB, productID int64, stock int) error {
	ctx := context.Background()
	_, err := gorm.G[entity.Product](db).Where("id = ?", productID).Update(ctx, "stock", stock)
	if err != nil {
		return err
	}
	return nil
}
