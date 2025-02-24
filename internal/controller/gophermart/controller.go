package controller

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"net/http"
	"time"

	"github.com/vilasle/gophermart/internal/controller"
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
type OrderInf struct {
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
type WithdrawalInf struct {
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
	authSvc     service.AuthorizationService
	orderUp     service.OrderService
	withdrawSvc service.WithdrawalService
	accrualSvc  service.AccrualService
}

// POST /api/user/register
func (c Controller) UserRegister(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		//check the body
		if r.Body == http.NoBody { // http.NoBody - not nil, len =0
			return controller.NewResponse(service.ErrInvalidFormat, nil, "", controller.ERROR)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 { // TODO: это лишняя проверка?
			return controller.NewResponse(service.ErrInvalidFormat, nil, "", controller.ERROR)
		}

		// proxy struct to unmarshal
		regReq := registerReq{}
		// Unmarshal login and password
		err = json.Unmarshal(body, &regReq)
		if err != nil {
			return controller.NewResponse(err, nil, "", controller.ERROR)
		}

		// fill the acceptable struct for response
		user := service.RegisterRequest{
			Login:    regReq.Login,
			Password: regReq.Password,
		}

		// Передаю логин и пароль в сервис на проверку, получаем userID
		userID, err := c.authSvc.Register(r.Context(), user) //
		if err != nil {
			return controller.NewResponse(err, nil, "", controller.ERROR)

		}
		// Если всё ок, то производим генерацию токена и его запись в куки
		tokenStr, err := genJWTTokenString(userID.ID)
		if err != nil {
			return controller.NewResponse(err, nil, "", controller.ERROR)
		}
		// generate response (set cookie) and response
		return controller.NewResponse(nil, nil, tokenStr, controller.TEXT)
	}
}

// POST /api/user/login
// TODO: it checks if user is registered if so => return Cookies for him
func (c Controller) UserLogin(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		//check the body
		if r.Body == http.NoBody { // http.NoBody - not nil, len =0
			return controller.NewResponse(service.ErrInvalidFormat, nil, "", controller.ERRORJSON)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 { // TODO: это лишняя проверка?
			return controller.NewResponse(service.ErrInvalidFormat, nil, "", controller.ERRORJSON)
		}

		// proxy struct to unmarshal
		regReq := registerReq{}
		// Unmarshal login and password
		err = json.Unmarshal(body, &regReq)
		if err != nil {
			return controller.NewResponse(err, nil, "", controller.ERROR)
		}
		// fill the acceptable struct for response
		user := service.AuthorizeRequest{
			Login:    regReq.Login,
			Password: regReq.Password,
		}

		// Передаю логин и пароль в сервис на проверку
		userInfo, err := c.authSvc.Authorize(r.Context(), user) //
		if err != nil {
			return controller.NewResponse(err, nil, "", controller.ERROR)
		}
		// Если всё ок, то производим генерацию токена
		tokenStr, err := genJWTTokenString(userInfo.ID)
		if err != nil {
			return controller.NewResponse(err, nil, "", controller.ERROR)
		}
		// set cookie to mold the response
		return controller.NewResponse(nil, nil, tokenStr, controller.TEXT)
	}
}

// POST /api/user/orders
// mw extract userID => to service  // if not userID - ask service to generate userID
// upload order number, check it
func (c Controller) RelateOrderWithUser(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		//check the body
		if r.Body == http.NoBody { // http.NoBody - not nil, len =0
			return controller.NewResponse(service.ErrInvalidFormat, nil, "", controller.ERROR)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 { // TODO: это лишняя проверка?
			return controller.NewResponse(service.ErrInvalidFormat, nil, "", controller.ERROR)
		}

		// move the string(body) into the func in service to check order number (LUNA) and save it
		err = c.orderUp.Register(r.Context(), service.RegisterOrderRequest{Number: string(body)})
		return controller.NewResponse(err, nil, "", controller.ERROR)

	}

}

// GET /api/user/orders
// Хендлер доступен только авторизованному пользователю. Номера заказа в выдаче должны быть отсортированы по времени
// загрузки от самых старых к самым новым. Формат даты — RFC3339.
func (c Controller) ListOrdersRelatedWithUser(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		if r.ContentLength != 0 { //TODO: НУЖНО ЛИ ПРОВЕРЯТЬ CONTENT-LENGTH == 0? - ДА, ибо может быть GET с телом?
			return controller.NewResponse(controller.ErrInvalidFormat, nil, "", controller.ERROR)
		}
		// get userID from jwt context (by the key) to get order list related with a specific user
		userID := r.Context().Value("userID")
		orderInfo, err := c.orderUp.List(r.Context(), service.ListOrderRequest{UserID: userID.(string)})
		if err != nil {
			return controller.NewResponse(err, nil, "", controller.ERROR)
		}
		//fill the proxy slice of structs (with struct tags) to marshal the response
		orInfo := make([]OrderInf, 0, len(orderInfo))
		for i := range orderInfo {
			orInfo = append(orInfo, OrderInf{Number: orderInfo[i].Number, Status: orderInfo[i].Status, Accrual: orderInfo[i].Accrual, CreatedAt: orderInfo[i].CreatedAt})
		}
		return controller.NewResponse(nil, orInfo, "", controller.JSON)

	}
}

// GET /api/user/balance
// хендлер доступен только авторизованному пользователю. В ответе должны содержаться данные о текущей сумме баллов
// лояльности, а также сумме использованных за весь период регистрации баллов.
func (c Controller) BalanceStateByUser(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		if r.ContentLength != 0 { //TODO: НУЖНО ЛИ ПРОВЕРЯТЬ CONTENT-LENGTH == 0?
			return controller.NewResponse(service.ErrInvalidFormat, nil, "", controller.ERROR)
		}
		// get userID from jwt context (by the key) to get order list related with a specific user
		userID := r.Context().Value("userID")
		balanceInfo, err := c.withdrawSvc.Balance(r.Context(), service.UserBalanceRequest{UserID: userID.(string)})
		if err != nil {
			return controller.NewResponse(err, nil, "", controller.ERROR)
		}
		// fill proxy struct to marshal response
		balInfo := UserBal{Current: balanceInfo.Current, Withdrawn: balanceInfo.Withdrawn}
		// mold the response
		return controller.NewResponse(nil, balInfo, "", controller.JSON)

	}
}

// POST /api/user/balance/withdraw
// Хендлер доступен только авторизованному пользователю. Номер заказа представляет собой гипотетический номер
// нового заказа пользователя, в счёт оплаты которого списываются баллы.
// Примечание: для успешного списания достаточно успешной регистрации запроса,
func (c Controller) Withdraw(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		//check the body
		if r.Body == http.NoBody { // http.NoBody - not nil, len =0
			return controller.NewResponse(service.ErrInvalidFormat, nil, "", controller.ERROR)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 { // TODO: это лишняя проверка?
			return controller.NewResponse(service.ErrInvalidFormat, nil, "", controller.ERROR)
		}
		// proxy struct to unmarshal OrderNumber & Sum
		type ProductRow struct {
			order string  `json:"order"`
			sum   float64 `json:"sum"`
		}
		withdrawalReq := ProductRow{}
		//unmarshalling
		err = json.Unmarshal(body, &withdrawalReq)
		if err != nil {
			return controller.NewResponse(err, nil, "", controller.ERROR)
		}

		// get userID from jwt context (by the key) to get order list related with a specific user
		userID := r.Context().Value("userID")

		err = c.withdrawSvc.Withdraw(r.Context(), service.WithdrawalRequest{UserID: userID.(string), OrderNumber: withdrawalReq.order, Sum: withdrawalReq.sum})
		return controller.NewResponse(err, nil, "", controller.ERROR)
	}
}

// GET /api/user/withdrawals (AUTH only)
func (c Controller) ListOfWithdrawals(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		if r.ContentLength != 0 { //////TODO: НУЖНО ЛИ ПРОВЕРЯТЬ CONTENT-LENGTH == 0?
			return controller.NewResponse(service.ErrInvalidFormat, nil, "", controller.ERROR)
		}
		// get userID from jwt context (by the key) to get order list related with a specific user
		userID := r.Context().Value("userID")

		withdrawalInfo, err := c.withdrawSvc.List(r.Context(), service.WithdrawalListRequest{UserID: userID.(string)})
		if err != nil {
			return controller.NewResponse(err, nil, "", controller.ERROR)
		}
		// create&fill the proxy struct to marshal data in response
		withdrawList := make([]WithdrawalInf, 0, len(withdrawalInfo))
		for i := range withdrawalInfo {
			withdrawList = append(withdrawList, WithdrawalInf{OrderNumber: withdrawalInfo[i].OrderNumber, Sum: withdrawalInfo[i].Sum, Status: withdrawalInfo[i].Status})
		}
		return controller.NewResponse(nil, withdrawList, "", controller.JSON)

	}
}

// TODO: Для взаимодействия с системой доступен один хендлер
// TODO: гофермарт - как сервер, тут этот хэндлер для чего? ОБычно же хэндлер обрабатывает запрос к определённому
// эндроинту, тут эндпоинт это  GET /api/orders/{number}. Кто будет к нему обращаться в гофермарте?? Или же этот эндпоинт
// посылает запрос на эндпоинт с тем же названием, но в accruel
// Вот в accruel
// gophermart, для получения бонусов.
//GET /api/orders/{number} —  этот хэндлер для того, чтобы к нему обращался гофермарт (получение информации о расчёте начислений баллов лояльности)
//POST /api/orders — сюда посылает запрос условный клиент (регистрация нового совершённого заказа);
// POST /api/goods —  сюда посылает запрос условный АДМИН акруэла (регистрация информации о новой механике вознаграждения за товар)

// GET /api/orders/{number} — получение информации о расчёте начислений баллов лояльности.
func (c Controller) GetCalculationInfo(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		// TODO: нужно ли проверять body да и вообще что-то рповерять, ведь тут доверенный сервис?
		//check the body
		if r.Body == http.NoBody { // http.NoBody - not nil, len =0
			return controller.NewResponse(service.ErrInvalidFormat, nil, "", controller.ERROR)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 { // TODO: это лишняя проверка?
			return controller.NewResponse(service.ErrInvalidFormat, nil, "", controller.ERROR)
		}
		// get order number
		orderNum := r.PathValue("number")

		// get order processing info
		accrualInf, err := c.accrualSvc.Accruals(r.Context(), service.AccrualsFilterRequest{Number: orderNum})
		if err != nil {
			return controller.NewResponse(err, nil, "", controller.ERROR)
		}
		//fill proxy-struct to mold response
		accInf := AccrualsInf{OrderNumber: accrualInf.OrderNumber, Status: accrualInf.Status, Accrual: accrualInf.Accrual}
		return controller.NewResponse(nil, accInf, "", controller.JSON)
	}
}
func genJWTTokenString(userID string) (string, error) { // создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	type JWTClaims struct {
		jwt.RegisteredClaims
		userID string
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			// set expiration time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		// set my own statement
		userID: userID,
	})

	// создаём строку токена с подписью
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	// возвращаем строку токена
	return tokenString, nil

}
