package handler

import (
	"net/http"
	"stock-management/internal/service"
	"stock-management/pkg/dto"
	"stock-management/pkg/response"

	"github.com/goccy/go-json"
)

type Product interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
}
type productHand struct {
	productService service.Product
}

func NewProduct(productService service.Product) Product {
	return &productHand{productService: productService}
}
func (h *productHand) Create(w http.ResponseWriter, r *http.Request) {
	request := new(dto.ProductAddRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		panic(response.Except(400, "failed to parse json"))
	}
	err := h.productService.Create(request)
	if err != nil {
		panic(err)
	}
	response.Success(w, http.StatusCreated, "OK", r.URL.Path)
}
func (h *productHand) GetAll(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	size := r.URL.Query().Get("size")
	result, err := h.productService.GetAll(page, size)
	if err != nil {
		panic(err)
	}
	response.Success(w, http.StatusOK, result, r.URL.Path)
}
