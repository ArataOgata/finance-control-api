package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	resp "go-api/internal/dto/base"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(secretKey string) func(http.Handler) http.Handler {
	jwtSecret := []byte(secretKey)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				resp.SendError(w, http.StatusUnauthorized, "Требуется авторизация")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				resp.SendError(w, http.StatusUnauthorized, "Неверный формат заголовка Authorization")
				return
			}

			tokenStr := parts[1]

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return jwtSecret, nil
			})

			if err != nil || !token.Valid {
				resp.SendError(w, http.StatusUnauthorized, "Недействительный или истёкший токен")
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || claims == nil {
				resp.SendError(w, http.StatusUnauthorized, "Неверный формат claims")
				return
			}

			userID, ok := claims["sub"].(string)
			if !ok || userID == "" {
				resp.SendError(w, http.StatusUnauthorized, "В токене отсутствует или некорректный sub")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromRequest(r *http.Request) (string, bool) {
	uid, ok := r.Context().Value(UserIDKey).(string)
	return uid, ok
}
