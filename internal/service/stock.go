package service

import (
	"fmt"
	"stock-management/internal/entity"
	"stock-management/internal/repository"
	"stock-management/pkg/dto"
	"stock-management/pkg/enum"
	"stock-management/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Stock interface {
	In(request *dto.StockRequest) error
	Out(request *dto.StockRequest) error
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
