package auth

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/yaricks657/final-project/internal/manager"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		manager.Mng.Log.LogInfo("авторизация идет ", r.RemoteAddr)

		pass := manager.Mng.Cnf.Password
		fmt.Println(pass)
		if len(pass) > 0 {
			cookie, err := r.Cookie("token")
			if err != nil {
				sendErrorResponse(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			tokenStr := cookie.Value
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return []byte(manager.Mng.Cnf.JWTSecret), nil
			})
			if err != nil || !token.Valid {
				sendErrorResponse(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				sendErrorResponse(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			storedPasswordHash := hashPassword(pass)
			if claims["passwordHash"] != storedPasswordHash {
				sendErrorResponse(w, "Authentication required", http.StatusUnauthorized)
				return
			}
		}
		manager.Mng.Log.LogInfo("авторизация успешна ", r.RemoteAddr)

		next(w, r)
	})
}
