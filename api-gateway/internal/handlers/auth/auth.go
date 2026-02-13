package authHandlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	authproto "shop/auth-service/api/gen/go/api/proto"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
)

type AuthHandler struct {
	auth Auth
	log  *slog.Logger
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

type Auth interface {
	Register(ctx context.Context, email, password string, fullName string) (*authproto.RegisterResponse, error)
	Login(ctx context.Context, email, password string) (*authproto.LoginResponse, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

func NewAuthHandler(a Auth, l *slog.Logger) *AuthHandler {
	return &AuthHandler{
		auth: a,
		log:  l,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.Auth.Register"

	var rr registerReq

	if err := json.NewDecoder(r.Body).Decode(&rr); err != nil {
		h.log.Error("%s: %w", op, err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	resp, err := h.auth.Register(ctx, rr.Email, rr.Password, rr.FullName)
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
	const op = "Handler.Auth.Login"

	var lr loginReq

	if err := json.NewDecoder(r.Body).Decode(&lr); err != nil {
		h.log.Error("%s: %w", op, err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	resp, err := h.auth.Login(ctx, lr.Email, lr.Password)
	if err != nil {
		h.log.Error("%s: %w", op, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	b, _ := protojson.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")

	w.Write(b)
}
