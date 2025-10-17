package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/TeseySTD/GoHospitalApi/auth"
)

type contextKey string

const (
	UserContextKey contextKey = "user"
	RoleContextKey contextKey = "role"
)

func JWTAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		
		if authHeader == "" {
			respondError(w, http.StatusUnauthorized, "Authorization header missing")
			return
		}
		
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			respondError(w, http.StatusUnauthorized, "Invalid authorization format. Use: Bearer <token>")
			return
		}
		
		tokenString := parts[1]
		
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}
		
		ctx := context.WithValue(r.Context(), UserContextKey, claims.Username)
		ctx = context.WithValue(ctx, RoleContextKey, claims.Role)
		
		next(w, r.WithContext(ctx))
	}
}

func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(RoleContextKey).(string)
		if !ok {
			respondError(w, http.StatusUnauthorized, "Role not found in context")
			return
		}
		
		if !auth.IsAdmin(role) {
			respondError(w, http.StatusForbidden, "Admin access required")
			return
		}
		
		next(w, r)
	}
}

func RoleBasedAccess(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(RoleContextKey).(string)
		if !ok {
			respondError(w, http.StatusUnauthorized, "Role not found in context")
			return
		}
		
		if auth.IsReader(role) && r.Method != http.MethodGet {
			respondError(w, http.StatusForbidden, "Read-only access. Only GET requests allowed")
			return
		}
		
		next(w, r)
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}