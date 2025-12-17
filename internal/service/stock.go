package service

import (
	"fmt"
	"math"
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

type Stock interface {
	In(request *dto.StockRequest) error
	Out(request *dto.StockRequest) error
	Movement(page, size string) (*dto.WebPagination[[]dto.MovementReposne], error)
}
type stockServ struct {
	db                 *gorm.DB
	validator          *validator.Validate
	logger             zerolog.Logger
	productRepository  repository.Product
	movementRepository repository.Movement
}

func NewStock(
	db *gorm.DB,
	validator *validator.Validate,
	logger zerolog.Logger,
	productRepository repository.Product,
	movementRepository repository.Movement,
) Stock {
	return &stockServ{
		db:                 db,
		validator:          validator,
		logger:             logger,
		productRepository:  productRepository,
		movementRepository: movementRepository,
	}
}

func (s *stockServ) In(request *dto.StockRequest) error {
	if err := s.validator.Struct(request); err != nil {
		s.logger.Warn().Msgf("validation failed: %v", err)
		return err
	}
	product, err := s.productRepository.FindByID(request.ProductID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.logger.Warn().Msgf("product %d not found", request.ProductID)
			return response.Except(404, "product not found")
		}
		s.logger.Error().Msg("failed find by id product to database")
		return err
	}
	stock := product.Stock + request.Qty
	tx := s.db.Begin()
	if err := s.productRepository.UpdateStock(tx, product.ID, stock); err != nil {
		tx.Rollback()
		s.logger.Error().Msg("failed update stock product to database")
		return err
	}
	movement := &entity.Movement{
		ProductID: product.ID,
		Type:      enum.TypeIn,
		Qty:       request.Qty,
	}
	if err := s.movementRepository.Create(tx, movement); err != nil {
		tx.Rollback()
		s.logger.Error().Msgf("failed create movement to database: %v", err)
		return err
	}
	tx.Commit()
	s.logger.Info().Str("product_id", fmt.Sprintf("%d", product.ID)).Msg("stock in success")
	return nil
}
func (s *stockServ) Out(request *dto.StockRequest) error {
	if err := s.validator.Struct(request); err != nil {
		s.logger.Warn().Msgf("validation failed: %v", err)
		return err
	}
	product, err := s.productRepository.FindByID(request.ProductID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.logger.Warn().Msgf("product %d not found", request.ProductID)
			return response.Except(404, "product not found")
		}
		s.logger.Error().Msg("failed find by id product to database")
		return err
	}
	if product.Stock < 1 {
		s.logger.Warn().Msgf("stock %d empty", product.Stock)
		return response.Except(400, "stock empty")
	}
	stock := product.Stock - request.Qty
	tx := s.db.Begin()
	if err := s.productRepository.UpdateStock(tx, product.ID, stock); err != nil {
		tx.Rollback()
		s.logger.Error().Msg("failed update stock product to database")
		return err
	}
	movement := &entity.Movement{
		ProductID: product.ID,
		Type:      enum.TypeIn,
		Qty:       request.Qty,
	}
	if err := s.movementRepository.Create(tx, movement); err != nil {
		tx.Rollback()
		s.logger.Error().Msgf("failed create movement to database: %v", err)
		return err
	}
	tx.Commit()
	s.logger.Info().Str("product_id", fmt.Sprintf("%d", product.ID)).Msg("stock out success")
	return nil
}
func (s *stockServ) Movement(page, size string) (*dto.WebPagination[[]dto.MovementReposne], error) {
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
	movements, err := s.movementRepository.Find(newPage, newSize)
	if err != nil {
		s.logger.Error().Msg("failed find movements to database")
		return nil, err
	}
	totalElement, err := s.movementRepository.Count()
	if err != nil {
		s.logger.Error().Msg("failed count movement to database")
		return nil, err
	}
	moveResp := make([]dto.MovementReposne, 0, len(movements))
	if len(movements) != 0 {
		for _, x := range movements {
			moveResp = append(moveResp, dto.MovementReposne{
				ID:        x.ID,
				ProductID: x.ProductID,
				Type:      x.Type,
				Qty:       x.Qty,
				CreatedAt: x.CreatedAt,
				Product: &dto.ProductResposne{
					ID:        x.Product.ID,
					Sku:       x.Product.Sku,
					Name:      x.Product.Name,
					Price:     x.Product.Price,
					Stock:     x.Product.Stock,
					CreatedAt: x.Product.CreatedAt,
					UpdatedAt: x.Product.UpdatedAt,
				},
			})
		}
	}
	totalPage := math.Ceil(float64(totalElement) / float64(newSize))
	resp := &dto.WebPagination[[]dto.MovementReposne]{
		Content:      moveResp,
		Page:         newPage,
		Size:         newSize,
		TotalPage:    int(totalPage),
		TotalElement: int(totalElement),
	}
	s.logger.Info().Str("total_element", fmt.Sprintf("%d", totalElement)).Msg("stock movements success")
	return resp, nil
}
