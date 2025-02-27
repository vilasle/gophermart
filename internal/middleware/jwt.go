package middleware

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/vilasle/gophermart/internal/service"
	"net/http"
	"time"
)

// TODO: implement it in config (env)
const tokenExp = time.Hour * 1
const secretKey = "supersecretkey"

type contextKey string

const UserIDKey contextKey = "userID"
const CookieKey string = "token"

// ///////////////////////////// MODELS //////////////////////////////////////////////////////
// Claims — структура утверждений, которая включает стандартные утверждения и
// одно пользовательское UserID
type JWTClaims struct {
	jwt.RegisteredClaims
	UserID string
}

// TODO: the logic
// 1. Проверка на наличие userID вообще, на валидность
// 2. Проверка expiration date, если она истекла, то прошу сервис сгенерить новый userID, генерю token, отдаю в куку
// 3. Проверка userID в сервисе. Если ОК => передаю упр хэндлеру, НЕ ОК => unathorized
// 4.

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

			// ИЗВЛЕКАЮ userID
			// создаём экземпляр структуры с утверждениями
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

			// TODO: ОТПРАВЛЯЮ claims.UserID на проверку в service
			// В ответ получаю ОК или не ОК ЕСЛИ НЕ ОК ВЕРНУ Unathorized
			// как в jwt лучше пробросить сервис для вызова, например в данном случае
			//
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
				Expires:  time.Now().Add(tokenExp), // TODO: есть внутри токена и в куке тжс?
			})

			// add userID to context to use it in controller
			ctx := context.WithValue(req.Context(), UserIDKey, claims.UserID)
			// expiration date - OK - continue
			next.ServeHTTP(res, req.WithContext(ctx)) // continue

		})

	}
}

/*
func JWTMiddleware(h http.Handler) http.Handler {
	jwtFunc := func(res http.ResponseWriter, req *http.Request) {
		// get token string from the cookies
		tokenString, err := req.Cookie("token")
		// check it
		if err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		if tokenString.Value == "" {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		// ИЗВЛЕКАЮ userID
		// создаём экземпляр структуры с утверждениями
		claims := &JWTClaims{}
		// парсим из строки токена tokenString в структуру claims
		token, err := jwt.ParseWithClaims(tokenString.Value, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				res.WriteHeader(http.StatusUnauthorized) //TODO [MENTOR] mb use StatusInternalServerError?
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

		// TODO: ОТПРАВЛЯЮ claims.UserID на проверку в service
		// В ответ получаю ОК или не ОК ЕСЛИ НЕ ОК ВЕРНУ Unathorized
		// как в jwt лучше пробросить сервис для вызова, например в данном случае
		//
		err = service.AuthorizationService.Authorize(req.Context(), ..)
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
			Expires:  time.Now().Add(tokenExp), // TODO: есть внутри токена и в куке тжс?
		}


		// add userID to context to use it in controller
		ctx := context.WithValue(req.Context(), "userID", claims.UserID)
		// expiration date - OK - continue
		h.ServeHTTP(res, req.WithContext(ctx)) // continue


		// TODO: вопрос: НУЖНО проверять тухлость. Опция token.Valid (1), также тухлость можно проверить вот так(2):
		//Check expiration date
		//if tokenString.Expires.After(time.Now()) {
		// }
		// первый (1) случай проверки можно сделать только ПОСЛЕ парсинга токена НО там помимио проверки на тухлость
		// висят и другие проверки .
		// второй (2) случай - ДО парсинга
		// Получается в 1 случае ненужные доп проверки и тогда сервису как действовать? Может вообще не проверять на Valid
		// тогда и проверить только тухлость до парсинг, либо научить сервис реагировать на эти доп проверки Valid
		// да и вообще получается DRY нарушение, я проверяю тухлость до парсинга, а потом Valid проверит его после парсинга


	}
	return http.HandlerFunc(jwtFunc)
}

*/

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

/*

func (c *Controller) setAuthToken(w http.ResponseWriter, tokenStr string) {

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenStr,
		Secure:   false,
		HttpOnly: true,
		Expires:  time.Now().Add(config.TokenExp),
	})

}

func (c *Controller) retrieveUserIDFromTokenString(r *http.Request) (string, error) {
	// get token string from the cookies
	tokenString, err := r.Cookie("token")

	if err != nil {
		c.logger.Info("No token!", zap.Error(err))
		return "", errors.New("no token")
	}
	// TODO: [MENTOR] SHOULD I CHECK
	if tokenString.Value == "" {
		c.logger.Info("Empty token!", zap.Error(err))
		return "", errors.New("empty token")
	}
	// создаём экземпляр структуры с утверждениями
	claims := &models.Claims{}
	// парсим из строки токена tokenString в структуру claims
	token, err := jwt.ParseWithClaims(tokenString.Value, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		} // anti-hacker check
		return []byte(config.SecretKey), nil
	})
	if err != nil {
		c.logger.Info("Can't parse token!", zap.Error(err))
		return "", errors.New("can't parse token")
	}

	if !token.Valid {
		c.logger.Info("Invalid token!", zap.Error(err))
		return "", errors.New("invalid token")
	}

	c.logger.Info("Successfully retrieved token!", zap.String("token", tokenString.Value))
	// возвращаем ID пользователя в читаемом виде
	return claims.UserID, nil

}





// ///////////////////////////////////////////////////////////////////////////////////////////
// the jwt middleware function
func RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Список эндпоинтов, для которых не требуется авторизация
		notAuth := []string{"/api/user/new", "/api/user/login"}
		requestPath := r.URL.Path //текущий путь запроса

		//проверяем, не требует ли запрос аутентификации
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		if !isAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		// Аутентификация прошла успешно, направляем запрос следующему обработчику
		next.ServeHTTP(w, r)
	})
}

// middleware 2
func JwtAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		t := strings.Split(authHeader, " ")
		if len(t) == 2 {
			authToken := t[1]
			authorized, err := tokenutil.IsAuthorized(authToken, secret)
			if authorized {
				userID, err := tokenutil.ExtractIDFromToken(authToken, secret)
				if err != nil {
					c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
					c.Abort()
					return
				}
				c.Set("x-user-id", userID)
				c.Next()
				return
			}
			c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
			c.Abort()
			return
		}
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Not authorized"})
		c.Abort()
	}


func CreateAccessToken(user *domain.User, secret string, expiry int) (accessToken string, err error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry)).Unix()
	claims := &domain.JwtCustomClaims{
		Name: user.Name,
		ID:   user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return t, err
}

func CreateRefreshToken(user *domain.User, secret string, expiry int) (refreshToken string, err error) {
	claimsRefresh := &domain.JwtCustomRefreshClaims{
		ID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(expiry)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	rt, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return rt, err
}

func IsAuthorized(requestToken string, secret string) (bool, error) {
	_, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func ExtractIDFromToken(requestToken string, secret string) (string, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok && !token.Valid {
		return "", fmt.Errorf("Invalid Token")
	}

	return claims["id"].(string), nil
}


*/
