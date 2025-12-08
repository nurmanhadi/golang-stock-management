package service

import (
	"fmt"
	"math"
	"net/http"
	"stock-management/internal/entity"
	"stock-management/internal/repository"
	"stock-management/pkg/dto"
	"stock-management/pkg/enum"
	"stock-management/pkg/response"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Product interface {
	Create(request *dto.ProductAddRequest) error
	GetAll(page, size string) (*dto.WebPagination[[]dto.ProductResposne], error)
}
type productServ struct {
	validator          *validator.Validate
	logger             zerolog.Logger
	db                 *gorm.DB
	productRepository  repository.Product
	movementRepository repository.Movement
}

func NewProduct(
	validator *validator.Validate,
	logger zerolog.Logger,
	db *gorm.DB,
	productRepository repository.Product,
	movementRepository repository.Movement,
) Product {
	return &productServ{
		validator:          validator,
		logger:             logger,
		db:                 db,
		productRepository:  productRepository,
		movementRepository: movementRepository,
	}
}
func (s *productServ) Create(request *dto.ProductAddRequest) error {
	if err := s.validator.Struct(request); err != nil {
		s.logger.Warn().Msgf("validation failed: %v", err)
		return err
	}
	totalProduct, err := s.productRepository.CountBySku(request.Sku)
	if err != nil {
		s.logger.Error().Msgf("failed count by sku to database: %v", err)
		return err
	}
	if totalProduct > 0 {
		s.logger.Warn().Msgf("sku %s already exists", request.Sku)
		return response.Except(http.StatusConflict, "sku already exists")
	}
	tx := s.db.Begin()
	product := &entity.Product{
		Sku:   request.Sku,
		Name:  request.Name,
		Price: request.Price,
		Stock: request.Stock,
	}
	productID, err := s.productRepository.Create(tx, product)
	if err != nil {
		tx.Rollback()
		s.logger.Error().Msgf("failed create product to database: %v", err)
		return err
	}
	movement := &entity.Movement{
		ProductID: productID,
		Type:      enum.TypeIn,
		Qty:       product.Stock,
	}
	if err := s.movementRepository.Create(tx, movement); err != nil {
		tx.Rollback()
		s.logger.Error().Msgf("failed create movement to database: %v", err)
		return err
	}
	tx.Commit()
	s.logger.Info().Str("product_id", fmt.Sprintf("%d", productID)).Msg("create product success")
	return nil
}
func (s *productServ) GetAll(page, size string) (*dto.WebPagination[[]dto.ProductResposne], error) {
	newPage, err := strconv.Atoi(page)
	if err != nil {
		s.logger.Error().Msgf("failed parse page: %s to int", page)
		return nil, err
	}
	newSize, err := strconv.Atoi(size)
	if err != nil {
		s.logger.Error().Msgf("failed parse size: %s to int", size)
		return nil, err
	}
	products, err := s.productRepository.Find(newPage, newSize)
	if err != nil {
		s.logger.Error().Msgf("failed find products to database: %v", err)
		return nil, err
	}
	totalProduct, err := s.productRepository.Count()
	if err != nil {
		s.logger.Error().Msgf("failed count product to database: %s", err)
		return nil, err
	}
	productRes := make([]dto.ProductResposne, 0, len(products))
	if len(products) != 0 {
		for _, x := range products {
			productRes = append(productRes, dto.ProductResposne{
				ID:        x.ID,
				Sku:       x.Sku,
				Name:      x.Name,
				Price:     x.Price,
				Stock:     x.Stock,
				CreatedAt: x.CreatedAt,
				UpdatedAt: x.UpdatedAt,
			})
		}
	}
	totalPage := math.Ceil(float64(totalProduct) / float64(newSize))
	resp := &dto.WebPagination[[]dto.ProductResposne]{
		Content:      productRes,
		Page:         newPage,
		Size:         newSize,
		TotalPage:    int(totalPage),
		TotalElement: int(totalProduct),
	}
	s.logger.Info().Str("total_element", fmt.Sprintf("%d", totalProduct)).Msg("get all product success")
	return resp, nil
}
