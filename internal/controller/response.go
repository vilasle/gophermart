package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vilasle/gophermart/internal/service"
)

type ControllerHandler func(*http.Request) Response

func (h ControllerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := h(r)
	resp.Write(w) // it calls func (t textResponse) Write(w http.ResponseWriter)  e.g.
}

type Response interface {
	Write(http.ResponseWriter)
}

// /////////////////////////////////////////////////////////////////////////////////////
type baseResponse struct {
	data    []byte
	cookies []http.Cookie
	header  map[string]string
	err     error
}

func (r baseResponse) Write(w http.ResponseWriter) {
	for _, cookie := range r.cookies {
		http.SetCookie(w, &cookie)
	}

	for k, v := range r.header {
		w.Header().Set(k, v)
	}

	w.WriteHeader(getErrorCode(r.err))
	w.Write(r.data)
}

type ResponseType = int

const (
	TypeText ResponseType = iota + 1
	TypeJson
)

type textResponse struct {
	data    any
	cookies []http.Cookie
	err     error
}

type jsonResponse struct {
	data    any
	cookies []http.Cookie
	err     error
}

func NewResponse(err error, data any, token string, kind ResponseType, cookies ...http.Cookie) Response {
	if cookies == nil {
		cookies = []http.Cookie{}
	}

	switch kind {
	case TypeJson:
		return jsonResponse{data: data, cookies: cookies, err: err} // TODO: mb MUST use data.([]controller.OrderInf) ???
	default:
		return textResponse{data: token, cookies: cookies, err: err}
	}
}

func (r textResponse) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")

	base := baseResponse{
		data:    []byte(r.data.(string)),
		cookies: r.cookies,
		err:     r.err,
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
		data:    body,
		cookies: r.cookies,
		err:     r.err,
		header: map[string]string{
			"Content-Type": "application/json",
		},
	}

	base.Write(w)

}

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
