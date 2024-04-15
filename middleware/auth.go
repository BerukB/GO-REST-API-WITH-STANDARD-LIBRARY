package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"github.com/BerukB/GO-REST-API-WITH-STANDARD-LIBRARY/common"
	"github.com/BerukB/GO-REST-API-WITH-STANDARD-LIBRARY/handlers"
)

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			handlers.HandleError(w, r, common.UNAUTHORIZED, "No Authorization", http.StatusUnauthorized)
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}

			return []byte("12345678"), nil
		})

		if err != nil {
			handlers.HandleError(w, r, common.UNAUTHORIZED, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			handlers.HandleError(w, r, common.UNAUTHORIZED, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Settin userID in the context

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			userID, ok := claims["userID"].(string)

			if !ok {

				customErr := &common.CustomError{
					Type:       common.UNAUTHORIZED,
					Message:    common.UNAUTHORIZED,
					StatusCode: http.StatusUnauthorized,
				}

				ctxWithError := context.WithValue(r.Context(), common.ErrorKey, customErr)
				r = r.WithContext(ctxWithError)
				handlers.HandleRequestError(w, r)

				return
			}

			ctxWithUserID := context.WithValue(r.Context(), common.UserIDKey, userID)
			r = r.WithContext(ctxWithUserID)
		}

		next.ServeHTTP(w, r)
	})
}
