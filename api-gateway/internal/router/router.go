package router

import (
	authHandlers "api-geteway/internal/handlers/auth"
	handler_cart "api-geteway/internal/handlers/cart"
	productHandler "api-geteway/internal/handlers/product"
	"api-geteway/internal/middleware"
	"net/http"
)

type Router struct {
	authHandler    *authHandlers.AuthHandler
	productHandler *productHandler.ProductHandler
	cartHandler    *handler_cart.CartHandler
}

func New(ah *authHandlers.AuthHandler, ph *productHandler.ProductHandler, cart_handler *handler_cart.CartHandler) *Router {
	return &Router{
		authHandler:    ah,
		productHandler: ph,
		cartHandler:    cart_handler,
	}
}

func (r *Router) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /login", r.authHandler.Login)
	mux.HandleFunc("POST /register", r.authHandler.Register)

	// Product routes (публичные)
	mux.HandleFunc("GET /products", r.productHandler.AllProducts)
	mux.HandleFunc("GET /products/search", r.productHandler.ProductByName)

	// Product routes (только для admin)
	mux.Handle("POST /products", middleware.AdminOnly(http.HandlerFunc(r.productHandler.CreateProduct)))

	mux.HandleFunc("POST /cart", r.cartHandler.AddItem)
}
