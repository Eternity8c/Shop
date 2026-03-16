package router

import (
	"cart-service/internal/handler"
	"net/http"
)

type Router struct {
	cartHandler *handler.CartHandler
}

func New(ch *handler.CartHandler) *Router {
	return &Router{
		cartHandler: ch,
	}
}

func (r *Router) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /cart", r.cartHandler.AddToCart)
}
