package auth

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/yaricks657/final-project/internal/manager"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		pass := manager.Mng.Cnf.Password
		if len(pass) > 0 {
			fmt.Println(r.Cookie("token"))
			cookie, err := r.Cookie("token")
			if err != nil {
				sendErrorResponse(w, fmt.Sprintf("Authentication required: %s", err), http.StatusUnauthorized)
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
				sendErrorResponse(w, "Authentication required2", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				sendErrorResponse(w, "Authentication required3", http.StatusUnauthorized)
				return
			}

			storedPasswordHash, err := hashPassword(pass)
			if err != nil {
				manager.Mng.Log.LogError("Ошибка при хешировании: ", err)
				return
			}
			if claims["passwordHash"] != storedPasswordHash {
				sendErrorResponse(w, "Authentication required4", http.StatusUnauthorized)
				return
			}
		}

		next(w, r)
	})
}
