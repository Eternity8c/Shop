package productHandler

import (
	"api-geteway/internal/client/product"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type ProductClient interface {
	AllProducts(ctx context.Context) ([]product.Product, error)
	ProductByName(ctx context.Context, name string) (product.Product, error)
	CreateProduct(ctx context.Context, dto product.CreateProductDTO) (int64, error)
}

type ProductHandler struct {
	client ProductClient
	log    *slog.Logger
}

func New(c ProductClient, log *slog.Logger) *ProductHandler {
	return &ProductHandler{client: c, log: log}
}

func (h *ProductHandler) AllProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	products, err := h.client.AllProducts(ctx)
	if err != nil {
		h.log.Error("AllProducts", "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) ProductByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	p, err := h.client.ProductByName(ctx, name)
	if err != nil {
		h.log.Error("ProductByName", "err", err)
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var dto product.CreateProductDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id, err := h.client.CreateProduct(ctx, dto)
	if err != nil {
		h.log.Error("CreateProduct", "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(id)
}
