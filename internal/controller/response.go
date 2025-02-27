package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vilasle/gophermart/internal/service"
)

type ControllerHandler func(r *http.Request) Response

func (h ControllerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := h(r)
	resp.Write(w) // it calls func (t textResponse) Write(w http.ResponseWriter)  e.g.
}

type Response interface {
	Write(http.ResponseWriter)
}

// /////////////////////////////////////////////////////////////////////////////////////
type baseResponse struct {
	data        []byte
	successCode int
	cookies     []http.Cookie
	header      map[string]string
	err         error
}

func (r baseResponse) Write(w http.ResponseWriter) {
	for _, cookie := range r.cookies {
		http.SetCookie(w, &cookie)
	}

	for k, v := range r.header {
		w.Header().Set(k, v)
	}

	w.WriteHeader(getErrorCode(r.err, r.successCode))
	w.Write(r.data)
}

type ResponseType = int

const (
	TypeText ResponseType = iota + 1
	TypeJSON
)

type textResponse struct {
	data        any
	successCode int
	cookies     []http.Cookie
	err         error
}

type jsonResponse struct {
	data        any
	successCode int
	cookies     []http.Cookie
	err         error
}

func NewResponse(err error, data any, kind ResponseType, httpCodeIfSuccess int, cookies ...http.Cookie) Response {
	if cookies == nil {
		cookies = []http.Cookie{}
	}

	if httpCodeIfSuccess == 0 {
		httpCodeIfSuccess = http.StatusOK
	}

	switch kind {
	case TypeJSON:
		return jsonResponse{data: data, successCode: httpCodeIfSuccess, cookies: cookies, err: err} // TODO: mb MUST use data.([]controller.OrderInf) ???
	default:
		return textResponse{data: data, successCode: httpCodeIfSuccess, cookies: cookies, err: err}
	}
}

func (r textResponse) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")

	var body []byte

	if v, ok := r.data.(string); ok {
		body = []byte(v)
	}

	base := baseResponse{
		data:        body,
		successCode: r.successCode,
		cookies:     r.cookies,
		err:         r.err,
		header: map[string]string{
			"Content-Type": "text/plain",
		},
	}

	base.Write(w)
}

func (r jsonResponse) Write(w http.ResponseWriter) {
	var (
		body []byte
		err  = r.err
	)

	body, err = json.Marshal(r.data)
	if err != nil { //
		body = []byte(`{"error": "StatusInternalServerError"}`)
	}

	base := baseResponse{
		data:        body,
		successCode: r.successCode,
		cookies:     r.cookies,
		err:         r.err,
		header: map[string]string{
			"Content-Type": "application/json",
		},
	}

	base.Write(w)

}

func getErrorCode(err error, successCode int) int {
	if err == nil {
		return successCode
	}

	if errors.Is(err, service.ErrInvalidFormat) {
		return http.StatusBadRequest
	}

	if errors.Is(err, service.ErrNotEnoughPoints) { // for POST /api/orders in accrual (status 202)
		return http.StatusPaymentRequired
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
