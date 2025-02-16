package controller

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/vilasle/gophermart/internal/controller"
	"github.com/vilasle/gophermart/internal/service"
)


type registerRequest struct {
	login    string `json:"login"`
	password string `json:"password"`
}

type authorizeRequest struct {
	login    string
	password string
}

type Controller struct {
	authSvc service.AuthorizationService
}
// TODO: JWT   middleware "хэндлер доступен только аутентифицированным пользователям"
// POST /api/user/register
func (c Controller) UserRegister(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		//check the body
		if r.Body == http.NoBody { // http.NoBody - not nil, len =0
			return controller.NewResponse([]byte{}, controller.ErrInvalidFormat, controller.TEXT)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 { // TODO: это лишняя проверка?
			return controller.NewResponse([]byte{}, controller.ErrInvalidFormat, controller.TEXT)
		}
		// Unmarshal login and password
		// proxy struct to unmarshal
		regReq := struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}

		err = json.Unmarshal(body, &regReq)
		if err != nil {
			return controller.NewResponse([]byte{}, controller.InternalError, controller.TEXT)
		}
		// fill the acceptable struct for response
		user := service.RegisterRequest {
			Login: regReq.Login,
			Password: regReq.Password,
		}

		// Передаю логин и пароль в сервис на проверку
		err, userID := c.authSvc.Register(r.Context(),user)   //



		if err != nil {
			return controller.NewResponse([]byte{}, err, controller.TEXT)

		}
		// TODO: Где расположить эту функцию генерации токена genJWTTokenString (кинул в jwt)? Она по идее стоит особняком от jwt middleware, т.к.
		// TODO: располагаем в controller и дёргаем в eregister и login
		// jwt middleware не распространяет своё действие на register и login?
		token, err :=
		return c.authSvc.Register(r.Context(), user) //
	}
}

// POST /api/user/login   - даёт userID

// TODO: it checks if user is registered if so => return Cookies for him
func (c Controller) UserLogin(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		//check the body
		if r.Body == http.NoBody { // http.NoBody - not nil, len =0
			return controller.NewResponse([]byte{}, controller.ErrInvalidFormat, controller.TEXT)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 { // TODO: это лишняя проверка?
			return controller.NewResponse([]byte{}, controller.ErrInvalidFormat, controller.TEXT)
		}
		// Unmarshal login and password
		var user authorizeRequest  // implement my own struct to unmarshall
		err = json.Unmarshal(body, &user)
		if err != nil {
			return controller.NewResponse([]byte{}, controller.InternalError, controller.TEXT)
		}
		// Передаю логин и пароль в сервис на проверку

		// TODO: как тут преедать третий параметр? Тут требуется конкретная struct, но в service она без struct tag
		err, userID := c.authSvc.Authorize(r.Context(),user)   // TODO: ПЕРЕДАЮ login & password в сервис при этом создал свою authorizeRequest с struct tags
		if err != nil {
			// TODO: Как мне тут распознать и вернуть ошибку сервиса: через отдельную функцию типа checkErr в controller?
			/*Возможные коды ответа:
			  200 — пользователь успешно аутентифицирован;
			  400 — неверный формат запроса;
			  401 — неверная пара логин/пароль;
			  500 — внутренняя ошибка сервера.
			*/

			return controller.NewResponse([]byte{}, controller.InternalError, controller.TEXT)
		}

	}
}

// POST /api/user/orders
// mw extract userID => to service  // if not userID - ask service to generate userID
// upload order number, check it
func (c Controller) RelateOrderWithUser(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		//check the body

		// move the string(body) into the func in service to check order number (LUNA) and save it
	}
}

// GET /api/user/orders
// Хендлер доступен только авторизованному пользователю. Номера заказа в выдаче должны быть отсортированы по времени
//загрузки от самых старых к самым новым. Формат даты — RFC3339.
func (c Controller) ListOrdersRelatedWithUser(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		//TODO: НУЖНО ЛИ ПРОВЕРЯТЬ CONTENT-LENGTH == 0?
		// TODO: retrieve userID from cookie and match it with the database?
	}
}

// GET /api/user/balance
// ендлер доступен только авторизованному пользователю. В ответе должны содержаться данные о текущей сумме баллов
//лояльности, а также сумме использованных за весь период регистрации баллов.
func (c Controller) BalanceStateByUser(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		////TODO: НУЖНО ЛИ ПРОВЕРЯТЬ CONTENT-LENGTH == 0?

	}
}

// POST /api/user/balance/withdraw
// Хендлер доступен только авторизованному пользователю. Номер заказа представляет собой гипотетический номер
//нового заказа пользователя, в счёт оплаты которого списываются баллы.
// Примечание: для успешного списания достаточно успешной регистрации запроса,
func (c Controller) Withdraw(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {

		panic("not implemented")
	}
}

// GET /api/user/withdrawals (AUTH only)
func (c Controller) ListOfWithdrawals(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		//////TODO: НУЖНО ЛИ ПРОВЕРЯТЬ CONTENT-LENGTH == 0?
		panic("not implemented")
	}
}

// TODO: Для взаимодействия с системой доступен один хендлер:
//
//    GET /api/orders/{number} — получение информации о расчёте начислений баллов лояльности.




func checkAndGetBody (body io.ReadCloser) []byte, error {
	body, err := io.ReadAll(body)
	if err != nil || len(body) == 0 {
		return nil, getErrorCode(err)
}
}

*/