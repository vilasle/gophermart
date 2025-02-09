package controller

import "net/http"

type ControllerHandler func(*http.Request) Response

func (h ControllerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := h(r)
	resp.Write(w)
}

type Response interface {
	Write(http.ResponseWriter)
}
