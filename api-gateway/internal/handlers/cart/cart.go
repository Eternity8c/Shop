package handler_cart

import (
	client_cart "api-geteway/internal/client/cart"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Cart interface {
	AddItem(ctx context.Context, userID, productID string, quantity int) error
}

type CartHandler struct {
	client client_cart.Client
	log    *slog.Logger
}

func NewCartHandler(c client_cart.Client, l *slog.Logger) *CartHandler {
	return &CartHandler{
		client: c,
		log:    l,
	}
}

func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	const op = "cart_handler.AddItem"

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// userID — из контекста (его туда положил auth-middleware гейтвея)
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// productID — из пути через роутер
	productID := mux.Vars(r)["productID"]
	if productID == "" {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// quantity — из тела
	quantity := 1
	if r.Body != nil {
		var req struct {
			Quantity int `json:"quantity"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil && req.Quantity > 0 {
			quantity = req.Quantity
		}
	}

	if err := h.client.AddItem(ctx, userID, productID, quantity); err != nil {
		h.log.Error("failed to add item", "op", op, "error", err)

		switch {
		case strings.Contains(err.Error(), "unauthorized"):
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		case strings.Contains(err.Error(), "bad request"):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "Internal error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message":    "Product added to cart",
		"product_id": productID,
		"quantity":   quantity,
	})
}
