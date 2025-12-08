package routes

import (
	"stock-management/delivery/rest/handler"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	Router         *chi.Mux
	ProductHandler handler.Product
}

func (r *Router) New() {
	r.Router.Route("/products", func(product chi.Router) {
		product.Post("/", r.ProductHandler.Create)
		product.Get("/", r.ProductHandler.GetAll)
	})
}
