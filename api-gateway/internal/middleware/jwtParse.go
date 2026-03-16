package middleware

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
)

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "unauthorixation", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		parts := strings.Split(tokenStr, ".")
		if len(parts) != 3 {
			http.Error(w, "invlaid token", http.StatusUnauthorized)
			return
		}

		payload, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		var claims map[string]any
		if err := json.Unmarshal(payload, &claims); err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		isAdmin, _ := claims["is_admin"].(bool)
		if !isAdmin {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
