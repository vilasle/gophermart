package gophermart

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v4"

	"github.com/vilasle/gophermart/internal/controller"
	"github.com/vilasle/gophermart/internal/logger"

	_mdw "github.com/vilasle/gophermart/internal/middleware"
	"github.com/vilasle/gophermart/internal/service"
)

// TODO [MENTOR]
// 1) If error => should I switch response to only the header with error (no content-type or text/plain) anyway??
// 2) Is it needed to set h.Set("X-Content-Type-Options", "nosniff") everywhere to evade malware js activation?
// 3) How to handle w.Write(dataMarsh) the best?
// 4) Лучше хранить токен в куках или в хэдэре авторизэйшн (почему для хэдеров не нужна БД)?
// TODO: implement it in config (env)
const TokenExp = time.Hour * 1
const secretKey = "supersecretkey"

////////////////proxy-structs to convert data to structs with struct tags /////////////////////////////////////////////

// OrderInf is used to marshal data in GET /api/user/orders
type OrderInfo struct {
	Number    string    `json:"number"`
	Status    string    `json:"status"`
	Accrual   float64   `json:"accrual,omitempty"` // there may be no any reward
	CreatedAt time.Time `json:"uploaded_at"`
}

// regReq is used to unmarshal data in POST /api/user/register & POST /api/user/login
type registerReq struct { // TODO: лучше тут хранить или в хэндлере с т.з. памяти?
	Login    string `json:"login"` // TODO: if here => replace ProductRow
	Password string `json:"password"`
}

// UserBal is used to marshal response body in GET /api/user/balance
type UserBal struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

// WithdrawalInf is used as a proxy struct to marshal response body in GET /api/user/withdrawals
type WithdrawalInfo struct {
	OrderNumber string  `json:"order"`
	Sum         float64 `json:"sum"`
	Status      string  `json:"processed_at"`
}

// AccrualsInf is used as a proxy struct to unmarshal response body in GET /api/orders/{number}
type AccrualsInf struct {
	OrderNumber string  `json:"order"`
	Status      string  `json:"status"`
	Accrual     float64 `json:"accrual,omitempty"` // TODO: omitempty is it ok?
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Controller struct {
	AuthSvc     service.AuthorizationService
	OrderSvc    service.OrderService
	WithdrawSvc service.WithdrawalService
}

// POST /api/user/register
func (c Controller) UserRegister() controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		reqID := r.Context().Value(middleware.RequestIDKey)
		if reqID == nil {
			reqID = ""
		}

		log := logger.With("id", reqID.(string))

		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 {
			return controller.NewResponse(service.ErrInvalidFormat, nil, controller.TypeText, 0)
		}

		log.Info("registering user", "body", string(body))

		// proxy struct to unmarshal
		regReq := registerReq{}
		// Unmarshal login and password
		err = json.Unmarshal(body, &regReq)
		if err != nil {
			return controller.NewResponse(err, nil, controller.TypeText, 0)
		}

		// fill the acceptable struct for response
		user := service.RegisterRequest{
			Login:    regReq.Login,
			Password: regReq.Password,
		}

		// Передаю логин и пароль в сервис на проверку, получаем userID
		userID, err := c.AuthSvc.Register(r.Context(), user) //
		if err != nil {
			return controller.NewResponse(err, nil, controller.TypeText, 0)

		}
		// Если всё ок, то производим генерацию токена и его запись в куки
		tokenStr, err := genJWTTokenString(userID.ID)
		if err != nil {
			return controller.NewResponse(err, nil, controller.TypeText, 0)
		}

		// generate response (set cookie) and response
		return controller.NewResponse(nil, nil, controller.TypeText, http.StatusOK, http.Cookie{
			Name:     _mdw.CookieKey,
			Value:    tokenStr,
			Secure:   false,
			HttpOnly: true,
			Expires:  time.Now().Add(TokenExp),
		})
	}
}

// POST /api/user/login
func (c Controller) UserLogin() controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		reqID := r.Context().Value(middleware.RequestIDKey)
		if reqID == nil {
			reqID = ""
		}

		log := logger.With("id", reqID.(string))

		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 { // TODO: это лишняя проверка?
			return controller.NewResponse(service.ErrInvalidFormat, nil, controller.TypeJSON, 0)
		}
		log.Info("logging in user", "body", string(body))

		// proxy struct to unmarshal
		regReq := registerReq{}
		// Unmarshal login and password
		err = json.Unmarshal(body, &regReq)
		if err != nil {
			return controller.NewResponse(err, nil, controller.TypeText, 0)
		}
		// fill the acceptable struct for response
		user := service.AuthorizeRequest{
			Login:    regReq.Login,
			Password: regReq.Password,
		}

		// Передаю логин и пароль в сервис на проверку
		userInfo, err := c.AuthSvc.Authorize(r.Context(), user) //
		if err != nil {
			return controller.NewResponse(err, nil, controller.TypeText, 0)
		}
		// Если всё ок, то производим генерацию токена
		tokenStr, err := genJWTTokenString(userInfo.ID)
		if err != nil {
			return controller.NewResponse(err, nil, controller.TypeText, 0)
		}
		// set cookie to mold the response
		return controller.NewResponse(nil, nil, controller.TypeText, 0, http.Cookie{
			Name:     _mdw.CookieKey,
			Value:    tokenStr,
			Secure:   false,
			HttpOnly: true,
			Expires:  time.Now().Add(TokenExp),
		})
	}
}

// POST /api/user/orders
func (c Controller) RelateOrderWithUser() controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		reqID := r.Context().Value(middleware.RequestIDKey)
		if reqID == nil {
			reqID = ""
		}

		log := logger.With("id", reqID.(string))

		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 {
			return controller.NewResponse(service.ErrInvalidFormat, nil, controller.TypeText, 0)
		}
		log.Info("relate order with user", "body", string(body))

		userID, ok := r.Context().Value(_mdw.UserIDKey).(string)
		if !ok {
			return controller.NewResponse(service.ErrWrongNameOrPassword, nil, controller.TypeText, 0)
		}

		// move the string(body) into the func in service to check order number (LUNA) and save it
		err = c.OrderSvc.Register(r.Context(), service.RegisterOrderRequest{
			Number: string(body),
			UserID: userID,
		})

		return controller.NewResponse(err, nil, controller.TypeText, http.StatusAccepted)
	}

}

// GET /api/user/orders
func (c Controller) ListOrdersRelatedWithUser() controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		// get userID from jwt context (by the key) to get order list related with a specific user
		userID, ok := r.Context().Value(_mdw.UserIDKey).(string)
		if !ok {
			return controller.NewResponse(service.ErrWrongNameOrPassword, nil, controller.TypeText, 0)
		}

		orderInfo, err := c.OrderSvc.List(r.Context(), service.ListOrderRequest{UserID: userID})
		if err != nil {
			return controller.NewResponse(err, nil, controller.TypeText, 0)
		}
		//fill the proxy slice of structs (with struct tags) to marshal the response
		orInfo := fillListOfOrders(orderInfo)
		return controller.NewResponse(nil, orInfo, controller.TypeJSON, 0)
	}
}

// GET /api/user/balance
func (c Controller) BalanceStateByUser() controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		reqID := r.Context().Value(middleware.RequestIDKey)
		if reqID == nil {
			reqID = ""
		}

		log := logger.With("id", reqID.(string))

		// get userID from jwt context (by the key) to get order list related with a specific user
		userID, ok := r.Context().Value(_mdw.UserIDKey).(string)
		if !ok {
			return controller.NewResponse(service.ErrWrongNameOrPassword, nil, controller.TypeText, 0)
		}
		log.Info("getting balance", "userID", userID)

		balanceInfo, err := c.WithdrawSvc.Balance(r.Context(), service.UserBalanceRequest{UserID: userID})
		if err != nil {
			return controller.NewResponse(err, nil, controller.TypeText, 0)
		}
		// fill proxy struct to marshal response
		balInfo := UserBal{Current: balanceInfo.Current, Withdrawn: balanceInfo.Withdrawn}
		// mold the response
		return controller.NewResponse(nil, balInfo, controller.TypeJSON, 0)

	}
}

// POST /api/user/balance/withdraw
func (c Controller) Withdraw() controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		reqID := r.Context().Value(middleware.RequestIDKey)
		if reqID == nil {
			reqID = ""
		}

		log := logger.With("id", reqID.(string))

		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 { // TODO: это лишняя проверка?
			return controller.NewResponse(service.ErrInvalidFormat, nil, controller.TypeText, 0)
		}

		log.Info("withdraw", "body", string(body))

		userID, ok := r.Context().Value(_mdw.UserIDKey).(string)
		if !ok {
			return controller.NewResponse(service.ErrWrongNameOrPassword, nil, controller.TypeText, 0)
		}

		log.Info("withdraw", "userID", userID)

		// proxy struct to unmarshal OrderNumber & Sum
		inputBody := struct {
			Order string  `json:"order"`
			Sum   float64 `json:"sum"`
		}{}

		err = json.Unmarshal(body, &inputBody)
		if err != nil {
			return controller.NewResponse(err, nil, controller.TypeText, 0)
		}

		err = c.WithdrawSvc.Withdraw(
			r.Context(), service.WithdrawalRequest{
				UserID:      userID,
				OrderNumber: inputBody.Order,
				Sum:         inputBody.Sum,
			})

		return controller.NewResponse(err, nil, controller.TypeText, 0)
	}
}

// GET /api/user/withdrawals (AUTH only)
func (c Controller) ListOfWithdrawals() controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		reqID := r.Context().Value(middleware.RequestIDKey)
		if reqID == nil {
			reqID = ""
		}

		log := logger.With("id", reqID.(string))

		// get userID from jwt context (by the key) to get order list related with a specific user
		userID, ok := r.Context().Value(_mdw.UserIDKey).(string)
		if !ok {
			return controller.NewResponse(service.ErrWrongNameOrPassword, nil, controller.TypeText, 0)
		}
		log.Info("getting withdrawals", "userID", userID)

		withdrawalInfo, err := c.WithdrawSvc.List(r.Context(), service.WithdrawalListRequest{UserID: userID})
		if err != nil {
			return controller.NewResponse(err, nil, controller.TypeText, 0)
		}

		withdrawList := fillListOfWithdrawals(withdrawalInfo)

		return controller.NewResponse(nil, withdrawList, controller.TypeJSON, 0)
	}
}

func genJWTTokenString(userID string) (string, error) {
	type JWTClaims struct {
		jwt.RegisteredClaims
		UserID string
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			// set expiration time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},

		UserID: userID,
	})
	return token.SignedString([]byte(secretKey))
}

func fillListOfOrders(orderInfo []service.OrderInfo) []OrderInfo {
	orders := make([]OrderInfo, 0, len(orderInfo))
	for _, v := range orderInfo {
		order := OrderInfo{
			Number:    v.Number,
			Status:    v.Status,
			Accrual:   v.Accrual,
			CreatedAt: v.CreatedAt,
		}
		orders = append(orders, order)
	}
	return orders
}

func fillListOfWithdrawals(withdrawalInfo []service.WithdrawalInfo) []WithdrawalInfo {
	withdrawList := make([]WithdrawalInfo, 0, len(withdrawalInfo))
	for _, v := range withdrawalInfo {
		s := v.Sum
		if s < 0 {
			s = -s
		}
		withdrawList = append(withdrawList, WithdrawalInfo{
			OrderNumber: v.OrderNumber,
			Sum:         s,
			Status:      v.CreatedAt.Format(time.RFC3339),
		})
	}
	return withdrawList
}
