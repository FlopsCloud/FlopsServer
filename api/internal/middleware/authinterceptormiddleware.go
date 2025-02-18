package middleware

import (
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthInterceptorMiddleware struct {
	Secret string
}

func NewAuthInterceptorMiddleware(secret string) *AuthInterceptorMiddleware {
	return &AuthInterceptorMiddleware{
		Secret: secret,
	}
}

func (m *AuthInterceptorMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		logx.Info(tokenString)
		// claims, err := token.ParseToken(tokenString, []byte(m.Secret))
		// if err != nil {
		// 	http.Error(w, "Invalid token", http.StatusUnauthorized)
		// 	return
		// }

		// // Extract user information from claims
		// userId, ok := claims["userId"].(string)
		// if !ok {
		// 	http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		// 	return
		// }

		// // Add user information to the request context
		// ctx := context.WithValue(r.Context(), "userId", userId)
		// r = r.WithContext(ctx)

		next(w, r)
	}
}
