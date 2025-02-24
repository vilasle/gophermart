package controller

import (
	"encoding/json"
	"errors"
	controller "github.com/vilasle/gophermart/internal/controller/gophermart"
	"github.com/vilasle/gophermart/internal/service"
	"net/http"
	"time"
)

// ControllerHandler implements HandlerFunc
type ControllerHandler func(*http.Request) Response

// define ServeHTTP to implement HandlerFunc type
func (h ControllerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := h(r)
	resp.Write(w) // it calls func (t textResponse) Write(w http.ResponseWriter)  e.g.
}

type Response interface {
	Write(http.ResponseWriter)
}

// /////////////////////////////////////////////////////////////////////////////////////
type ResponseType = int

const ( //  TODO: figure it our
	TEXT = iota + 1
	JSON
	ERRORJSON
	ERRORTEXT
	ERROR
)

type simpleTextResponseNoBody struct {
	//text string
	token string // TODO: mb put it into data in signature like an element of a struct?
	err   error
}

type jsonResponse struct {
	data any
	err  error
}

type responseWithJSError struct { // TODO: mb I can somehow unite it with JSON>?
	emptyData []byte
	err       error
}

type responseWithTEXTError struct { // TODO: mb I can somehow unite it with JSON>?
	emptyData []byte
	err       error
}

type responseWithError struct {
	err error
}

func NewResponse(err error, data any, token string, kind ResponseType) Response {
	switch kind {
	case TEXT:
		return simpleTextResponseNoBody{token: token, err: err}
	case JSON:
		return jsonResponse{data: data, err: err} // TODO: mb MUST use data.([]controller.OrderInf) ???
	case ERRORJSON:
		return responseWithJSError{emptyData: []byte(data.(string)), err: err}
	case ERRORTEXT:
		return responseWithTEXTError{emptyData: []byte(data.(string)), err: err}
	case ERROR:
		return responseWithError{err: err}

	default:
		return nil
	}
}
func (r responseWithError) Write(w http.ResponseWriter) {
	w.WriteHeader(getErrorCode(r.err))
	return // TODO: надо ли?
}

func (r responseWithJSError) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(getErrorCode(r.err))
	// w.Write(r.emptyData) TODO: убрал запись пустых байт
}

func (r responseWithTEXTError) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(getErrorCode(r.err))
	//w.Write(r.emptyData) TODO: убрал запись пустых байт
}

// create Write method to implement Response interface here
func (r simpleTextResponseNoBody) Write(w http.ResponseWriter) {
	// add token to the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    r.token,
		Secure:   false,
		HttpOnly: true,
		Expires:  time.Now().Add(controller.TokenExp), // coincides with token options
	})
	w.WriteHeader(getErrorCode(r.err))
}

// create Write method to implement Response interface here
func (r jsonResponse) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	dataMarsh, err := json.Marshal(r.data)
	if err != nil { //
		intErr := `{"error": "StatusInternalServerError"}`
		http.Error(w, intErr, http.StatusInternalServerError)
		return // TODO: надо ли?
	}
	w.WriteHeader(getErrorCode(r.err))
	w.Write(dataMarsh) // TODO: how to handle it the best? http.Error ?
}

// I MUST do it cause it is not be recognized by getErrorCode (if use controller preset errors)
func getErrorCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if errors.Is(err, service.ErrInvalidFormat) {
		return http.StatusBadRequest
	}
	if errors.Is(err, service.StatusOrderSuccessfullyAccepted) { // for POST /api/orders in accrual (status 202)
		return http.StatusAccepted
	}
	if errors.Is(err, service.ErrDuplicate) {
		return http.StatusConflict // 409 - логин уже занят
	}
	if errors.Is(err, service.ErrWrongNameOrPassword) {
		return http.StatusUnauthorized //401 — неверная пара логин/пароль;
	}
	if errors.Is(err, service.ErrOrderUploadAnotherUser) {
		return http.StatusConflict // 409 — номер заказа уже был загружен другим пользователем;
	}

	if errors.Is(err, service.ErrWrongNumberOfOrder) {
		return http.StatusUnprocessableEntity // 422 — неверный формат номера заказа;
	}
	if errors.Is(err, service.ErrEntityDoesNotExists) {
		return http.StatusNoContent // 204 — заказ не зарегистрирован в системе расчёта.
	}
	if errors.Is(err, service.ErrLimit) {
		return http.StatusTooManyRequests // 429 — превышено количество запросов к сервису
	}

	return http.StatusInternalServerError

}

/*
// responses description
type simpleResponse struct {
	data []byte
	err     error
}
// to implement Response interface
func (r simpleResponse) Write(w http.ResponseWriter) {
	w.WriteHeader(getStatusCode(r.err))

	w.Write(r.data)
}

type JSONResponse struct {
	sp simpleResponse
}

func NewJSONResponse(content []byte, err error) Response {
	return JSONResponse{sp: simpleResponse{data: content, err: err}}
}

func (r JSONResponse) Write(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	r.sp.Write(w)

*/
