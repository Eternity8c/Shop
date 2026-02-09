package authHandlers

import (
	"api-geteway/internal/auth"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
)

type AuthHandler struct {
	authClient *auth.Client
	log        *slog.Logger
}

type registerReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthHandler(c *auth.Client, l *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authClient: c,
		log:        l,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// 1) json.Decoder -> registerReq
	// 2) ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second); defer cancel()
	// 3) resp, err := h.client.Register(ctx, req.Email, req.Password, req.FullName)
	// 4) если err -> h.log.Error(...); http.Error(...)
	// 5) b, _ := protojson.Marshal(resp); w.Header().Set("Content-Type", "application/json"); w.Write(b)
	const op = "Handlers.Auth.Register"

	var rr registerReq

	if err := json.NewDecoder(r.Body).Decode(&rr); err != nil {
		h.log.Error("%s: %w", op, err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	resp, err := h.authClient.Register(ctx, rr.Email, rr.Password, rr.FullName)
	if err != nil {
		h.log.Error("%s: %w", op, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	b, _ := protojson.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")

	w.Write(b)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// TODO: аналогично Register, но вызвать h.client.Login и вернуть LoginResponse
	const op = "Handler.Auth.Login"

	var lr loginReq

	if err := json.NewDecoder(r.Body).Decode(&lr); err != nil {
		h.log.Error("%s: %w", op, err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	resp, err := h.authClient.Login(ctx, lr.Email, lr.Password)
	if err != nil {
		h.log.Error("%s: %w", op, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	b, _ := protojson.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")

	w.Write(b)
}
