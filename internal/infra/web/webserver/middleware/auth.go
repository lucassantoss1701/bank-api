package middleware

import (
	"context"
	"errors"
	"lucassantoss1701/bank/configs"
	"lucassantoss1701/bank/internal/entity"
	"lucassantoss1701/bank/internal/infra/web"
	"lucassantoss1701/bank/internal/infra/web/responses"
	"net/http"

	"github.com/golang-jwt/jwt"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("authorization")

		claims, err := validateToken(token)
		if err != nil {
			message := err.Error()
			err := entity.NewErrorHandler(entity.UNAUTHORIZED_ERROR)
			err.Add(message)
			responses.Err(w, err)
			return
		}

		accountID, err := getAccountID(claims)
		if err != nil {
			message := err.Error()
			err := entity.NewErrorHandler(entity.UNAUTHORIZED_ERROR)
			err.Add(message)
			responses.Err(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), web.AccountIDKey, accountID)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func getAccountID(claims jwt.MapClaims) (string, error) {
	if accountID, ok := claims["account_id"].(string); ok {
		return accountID, nil
	}
	return "", errors.New("token claims without account id")
}

func validateToken(token string) (jwt.MapClaims, error) {
	if token == "" {
		return nil, errors.New("token must not be empty")
	}

	claims, err := parseToken(token)
	if err != nil {
		return nil, errors.New("token is invalid")
	}
	return claims, nil
}

func parseToken(token string) (jwt.MapClaims, error) {
	var secret = configs.Get().Security.Secret
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}
