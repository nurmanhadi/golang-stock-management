package config

import (
	"stock-management/delivery/rest/handler"
	"stock-management/delivery/rest/middleware"
	"stock-management/delivery/rest/routes"
	"stock-management/internal/repository"
	"stock-management/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Bootstrap struct {
	DB        *gorm.DB
	Logger    zerolog.Logger
	Validator *validator.Validate
	Router    *chi.Mux
}

func Initialize(deps *Bootstrap) {
	// publisher

	// cache

	// repository
	productRepo := repository.NewProduct(deps.DB)
	movementRepo := repository.NewMovement(deps.DB)

	// service
	productServ := service.NewProduct(deps.Validator, deps.Logger, deps.DB, productRepo, movementRepo)
	stockServ := service.NewStock(deps.DB, deps.Validator, deps.Logger, productRepo, movementRepo)

	// handler
	productHand := handler.NewProduct(productServ)
	stockHand := handler.NewStock(stockServ)

	// middleware
	deps.Router.Use(middleware.ErrorHandler)

	// routes
	r := routes.Router{
		Router:         deps.Router,
		ProductHandler: productHand,
		StockHandler:   stockHand,
	}
	r.New()

	// subcriber
}
