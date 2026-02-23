package router

import (
	authHandlers "api-geteway/internal/handlers/auth"
	productHandler "api-geteway/internal/handlers/product"
	"api-geteway/internal/middleware"
	"net/http"
)

type Router struct {
	authHandler    *authHandlers.AuthHandler
	productHandler *productHandler.ProductHandler
}

func New(ah *authHandlers.AuthHandler, ph *productHandler.ProductHandler) *Router {
	return &Router{
		authHandler:    ah,
		productHandler: ph,
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
}
