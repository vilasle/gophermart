package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/vilasle/gophermart/internal/service"
)

const tokenExp = time.Hour * 1
const secretKey = "supersecretkey"

type contextKey string

const UserIDKey contextKey = "userID"
const CookieKey string = "token"

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID string
}

func JWTMiddleware(some ...service.AuthorizationService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			fmt.Println("JWT MW Is available!")
			// get token string from the cookies
			tokenString, err := req.Cookie(CookieKey)
			// check it
			if err != nil {
				res.WriteHeader(http.StatusUnauthorized)
				return
			}
			if tokenString.Value == "" {
				res.WriteHeader(http.StatusUnauthorized)
				return
			}

			claims := &JWTClaims{}

			// парсим из строки токена tokenString в структуру claims
			token, err := jwt.ParseWithClaims(tokenString.Value, claims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					res.WriteHeader(http.StatusUnauthorized)                                 //TODO [MENTOR] mb use StatusInternalServerError?
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"]) // TODO: should I print it?
				} // anti-hacker check
				return []byte(secretKey), nil
			})
			if err != nil {
				// TODO: по логике норм ли? jwt по идее не обязан вызывать методы сервиса
				// (НАДО ЛИ ЧТО-ТО ПИСАТЬ ТИПА КАК СНИЗУ)
				// просим авторизоваться, т.к. с токеном что-то не так (например просрочился)
				http.Error(res, "Bad token, please authorize again", http.StatusUnauthorized) // TODO: mb set another header?
				return
			}
			// check является ли токен ненулевым, имеет ли он AccessToken и не истёк ли срок его действия.
			if !token.Valid {
				http.Error(res, "Token is not valid, please authorize again", http.StatusUnauthorized)
				return
			}

			err = some[0].CheckByUserID(req.Context(), claims.UserID)
			if err != nil {
				http.Error(res, "Failed to validate token", http.StatusUnauthorized)
				return // TODO: ВЫЙДУ ВООБЩЕ ИЗ ОБРАБОТКИ ЗАПРОСА И УПР НЕ УЙДЕТ В ХЭНДЛЕР ДАЛЬШЕ?
			}

			// ГЕНЕРЮ ТОКЕН, вставляю его в куку, передаю управление хэндлеру
			newToken, err := genJWTTokenString(claims.UserID)
			if err != nil {
				http.Error(res, "Failed to generate token", http.StatusInternalServerError)
			}
			http.SetCookie(res, &http.Cookie{
				Name:     "token",
				Value:    newToken,
				Secure:   false,
				HttpOnly: true,
				Expires:  time.Now().Add(tokenExp),
			})

			// add userID to context to use it in controller
			ctx := context.WithValue(req.Context(), UserIDKey, claims.UserID)
			// expiration date - OK - continue
			next.ServeHTTP(res, req.WithContext(ctx)) // continue

		})

	}
}

//////////////////////////////////////////////////////////

// genJWTTokenString create JWT token and return it in string type
func genJWTTokenString(userID string) (string, error) { // создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			// set expiration time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		// set my own statement
		UserID: userID,
	})

	// создаём строку токена с подписью
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	// возвращаем строку токена
	return tokenString, nil
}
