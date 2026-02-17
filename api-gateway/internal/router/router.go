package router

import (
	authHandlers "api-geteway/internal/handlers/auth"
	"net/http"
)

type Router struct {
	authHandler *authHandlers.AuthHandler
}

func New(ah *authHandlers.AuthHandler) *Router {
	return &Router{
		authHandler: ah,
	}
}

func (r *Router) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /login", r.authHandler.Login)
	mux.HandleFunc("POST /register", r.authHandler.Register)
}
