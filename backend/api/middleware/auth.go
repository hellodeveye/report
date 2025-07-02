package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/hellodeveye/report/pkg/auth"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// 检查Bearer token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// 将用户信息添加到上下文中
		ctx := context.WithValue(r.Context(), auth.UserOpenIDKey, claims.OpenID)
		ctx = context.WithValue(ctx, auth.UserNameKey, claims.Name)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuthMiddleware 可选认证中间件（不强制要求登录）
func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString := parts[1]
				if claims, err := auth.ValidateToken(tokenString); err == nil {
					// 将用户信息添加到上下文中
					ctx := context.WithValue(r.Context(), auth.UserOpenIDKey, claims.OpenID)
					ctx = context.WithValue(ctx, auth.UserNameKey, claims.Name)
					r = r.WithContext(ctx)
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
