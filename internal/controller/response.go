package controller

/*
import (
	json2 "encoding/json"
	"errors"
	"github.com/vilasle/gophermart/internal/service"
	"net/http"
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

///////////////////////////////////////////////////////////////////////////////////////
type ResponseType = int

const (
	TEXT = iota +1
	JSON
)

type textResponse struct{
	text []byte
	err error
}

type jsonResponse struct{
	data []byte
	err error
}

func NewResponse(data any, err error, kind ResponseType) Response {
	switch kind {
	case TEXT:
		return textResponse{text: []byte(data.(string)), err: err}
	case JSON:
		return jsonResponse{data: []byte(data.(string)), err: err}
	default:
		return nil
	}
}

// create Write method to implement Response interface here
func (r textResponse) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(getErrorCode(r.err))
	w.Write(r.text)
}
// create Write method to implement Response interface here
func (r jsonResponse) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	json, err := json.Marshal(r.data)
	if err != nil {
		r.err = err
	}else {
		r.data = json

	}
	w.WriteHeader(getErrorCode(r.err))
	w.Write(r.data)

}



func getErrorCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if isBadRequest(err) {
		return http.StatusBadRequest // 400
	} else { // any possible undefined error
		return http.StatusInternalServerError // 500
	}
}

func isBadRequest (err error) bool {
	if errors.Is(err, service.ErrInvalidFormat) { // TODO: 401 for login handler (неверная пара логин и пароль)
		return true
	}
		if errors.Is(err, service.ErrDuplicate) { //
		return true
	}
		if errors.Is(err, service.ErrOrderUploadAnotherUser) {
		return true
	}
		if errors.Is(err, service.ErrWrongNumberOfOrder) {
		return true
	}
		if errors.Is(err, service.ErrEntityDoesNotExists) {
		return true
	}
		if errors.Is(err, service.ErrLimit) {
		return true
	}

}

}












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
