package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CartHandler struct {
	log     *slog.Logger
	service CartService
}

type CartService interface {
	Add(ctx context.Context, userID string, productID int, quantity int) (bool, error)
}

var (
	ErrMethodNotAllowed = errors.New("method not allowed")
)

func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.Cart.AddToCart"

	if r.Method != http.MethodPost {
		h.log.Error("method not allowed", "op", op, "method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		h.log.Error("user not authenticated", "op", op)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		h.log.Error("invalid path", "op", op, "path", r.URL.Path)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	productIDStr := pathParts[3]
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		h.log.Error("invalid product ID", "op", op, "product_id", productIDStr, "error", err)
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Quantity int `json:"quantity"`
	}

	quantity := 1
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil && req.Quantity > 0 {
			quantity = req.Quantity
		}
	}

	success, err := h.service.Add(ctx, userID, productID, quantity)
	if err != nil {
		h.log.Error("failed to add product", "op", op, "error", err)

		switch {
		case strings.Contains(err.Error(), "out of stock"):
			http.Error(w, "Product out of stock", http.StatusBadRequest)
		case strings.Contains(err.Error(), "insufficient stock"):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "Failed to add product to cart", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"success":    success,
		"message":    "Product added to cart",
		"product_id": productID,
		"quantity":   quantity,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error("failed to encode response", "op", op, "error", err)
	}
}
