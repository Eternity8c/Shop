package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"product-service/internal/dto"
	"product-service/internal/models"
	"product-service/internal/storage"
	"time"
)

type Handler struct {
	log            *slog.Logger
	productService ProductService
}

type ProductService interface {
	GetAllProducts(ctx context.Context) ([]models.Product, error)
	GetProductByName(ctx context.Context, name string) (models.Product, error)
	CreatedProduct(ctx context.Context, dto dto.Product) (int64, error)
}

var (
	ErrMethodNotAllowed = errors.New("method not allowed")
)

func New(ctx context.Context, log *slog.Logger, ps ProductService) *Handler {
	return &Handler{
		productService: ps,
		log:            log,
	}
}

func (h *Handler) AllProducts(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.AllProducts"

	if r.Method != http.MethodGet {
		h.log.Error("%s: %w", op, ErrMethodNotAllowed)
		http.Error(w, ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	products, err := h.productService.GetAllProducts(ctx)
	if err != nil {
		h.log.Error("%s: %w", op, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(products); err != nil {
		h.log.Error("%s: %w", op, err)
		return
	}
}

func (h *Handler) ProductByName(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.ProductByName"

	if r.Method != http.MethodGet {
		h.log.Error("%s: %w", op, ErrMethodNotAllowed)
		http.Error(w, ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	name := r.URL.Query().Get("name")

	products, err := h.productService.GetProductByName(ctx, name)
	if err != nil {
		if errors.Is(err, storage.ErrProductNotFound) {
			h.log.Info("%s: %w", op, storage.ErrProductNotFound)
			http.Error(w, "product not found", http.StatusNotFound)
			return
		}
		h.log.Error("%s: %w", op, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(products); err != nil {
		h.log.Error("%s: %w", op, err)
		return
	}
}

func (h *Handler) CreatedProduct(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.CreatedProduct"

	if r.Method != http.MethodPost {
		h.log.Error("%s: %w", op, ErrMethodNotAllowed)
		http.Error(w, ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}

	var dto dto.Product

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		h.log.Error("%s: %w", op, err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := h.productService.CreatedProduct(ctx, dto)
	if err != nil {
		h.log.Error("%s: %w", op, err)
		http.Error(w, "interanl server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(id); err != nil {
		h.log.Error("%s: %w", op, err)
		return
	}
}
