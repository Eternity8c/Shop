package router

import (
	"net/http"
	"product-service/internal/handler"
)

type Router struct {
	handler *handler.Handler
}

func New(h *handler.Handler) *Router {
	return &Router{
		handler: h,
	}
}

func (r *Router) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /Products", r.handler.AllProducts)
	mux.HandleFunc("GET /Product/search", r.handler.ProductByName)
	mux.HandleFunc("GET /Product/{ID}", r.handler.ProductByID)
	mux.HandleFunc("POST /created", r.handler.CreatedProduct)
}
