package handler

import (
	"net/http"
	"stock-management/internal/service"
	"stock-management/pkg/dto"
	"stock-management/pkg/response"

	"github.com/goccy/go-json"
)

type Stock interface {
	In(w http.ResponseWriter, r *http.Request)
	Out(w http.ResponseWriter, r *http.Request)
	Movement(w http.ResponseWriter, r *http.Request)
}
type stockHand struct {
	stockService service.Stock
}

func NewStock(stockService service.Stock) Stock {
	return &stockHand{
		stockService: stockService,
	}
}
func (h *stockHand) In(w http.ResponseWriter, r *http.Request) {
	request := new(dto.StockRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		panic(response.Except(400, "failed to parse json"))
	}
	err := h.stockService.In(request)
	if err != nil {
		panic(err)
	}
	response.Success(w, http.StatusOK, "OK", r.URL.Path)
}
func (h *stockHand) Out(w http.ResponseWriter, r *http.Request) {
	request := new(dto.StockRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		panic(response.Except(400, "failed to parse json"))
	}
	err := h.stockService.Out(request)
	if err != nil {
		panic(err)
	}
	response.Success(w, http.StatusOK, "OK", r.URL.Path)
}
func (h *stockHand) Movement(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	size := r.URL.Query().Get("size")
	if page == "" {
		page = "1"
	}
	if size == "" {
		size = "10"
	}
	result, err := h.stockService.Movement(page, size)
	if err != nil {
		panic(err)
	}
	response.Success(w, http.StatusOK, result, r.URL.Path)
}
