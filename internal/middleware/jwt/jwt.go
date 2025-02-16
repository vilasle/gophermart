package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"time"
)

/*
// ///////////////////////////// MODELS //////////////////////////////////////////////////////
// Claims — структура утверждений, которая включает стандартные утверждения и
// одно пользовательское UserID
type JWTClaims struct {
	jwt.RegisteredClaims
	UserID string
}

type JWTRefreshClaims struct {
	UserID string `json:"id"`
	jwt.RegisteredClaims
}

// continue without gzip
//h.ServeHTTP(res, or)

func JWTMiddleware(h http.Handler) http.Handler {
	jwtFunc := func(res http.ResponseWriter, req *http.Request) {
		// retrieve userID from the cookie
		// get token string from the cookies
		tokenString, err := req.Cookie("token")
		if err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		if tokenString.Value == "" {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		// TODO: Согласовать с backend проверку токена на тухлость
		//Check expiration date
		if tokenString.Expires. tokenString.Expires
		fmt.Println(cookie.Value)
		fmt.Println(cookie.Expires)


		// создаём экземпляр структуры с утверждениями
		claims := &JWTClaims{}
		// парсим из строки токена tokenString в структуру claims
		token, err := jwt.ParseWithClaims(tokenString.Value, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				res.WriteHeader(http.StatusUnauthorized)
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"]) // TODO: should I print it?
			} // anti-hacker check
			return []byte(config.SecretKey), nil // TODO: implement it from config (env)
		})
		if err != nil {
			http.Error(res, "Failed to parse token", http.StatusInternalServerError) // TODO: mb set another header?
			return
		}

		if !token.Valid { // TODO: проверяет ли это на
			http.Error(res, "Token is not valid", http.StatusUnauthorized)
			return
		}

		// TODO: ОТПРАВЛЯЮ claims.UserID на проверку в service или repository
		// В ответ получаю ОК или не ОК ЕСЛИ НЕ ОК ВЕРНУ Unathorized



	}
	return http.HandlerFunc(jwtFunc)
}





//////////////////////////////////////////////////////////

// genJWTTokenString create JWT token and return it in string type
func (c *Controller) genJWTTokenString() (string, string, error) { // TODO [MENTOR]: mb I should replace this func ???
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	//usId := string(time.Now().Unix())
	usID := service.GetRandString(time.Now().UTC().String())
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// set expiration time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.TokenExp)), //TODO [MENTOR] is it a good way to store it?
		},
		// set my own statement
		UserID: usID, // TODO [MENTOR]: how should I implement it better??
		// int(b[0] + b[1])
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(config.SecretKey)) // TODO [MENTOR]: how to store it better? how people store it in real projects? In env?
	// TODO: ok if env .. I set the env value secretKey on my PC e.g. and then start the app?
	if err != nil {
		return "", "", err
	}

	// возвращаем строку токена
	return tokenString, usID, nil

}




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
