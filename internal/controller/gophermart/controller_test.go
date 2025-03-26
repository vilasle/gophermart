package gophermart

import (
	"testing"
)

// create  a mock///////////////////////////////////////////////////////////////////////////////////////
// type TestService struct {
// 	err error
// }

// // can return defined errors ErrInvalidFormat, ErrDuplicate and undefined error
// func (s *TestService) Register(context.Context, service.RegisterOrderRequest) (service.UserInfo, error) {
// 	return service.UserInfo{}, s.err
// }

// // can return defined errors ErrInvalidFormat, ErrWrongNameOrPassword and undefined error
// func (s *TestService) Authorize(context.Context, service.AuthorizeRequest) (service.UserInfo, error) {
// 	return service.UserInfo{}, s.err
// }

// // func to work with mock as a pointer
// func (s *TestService) SetError(err error) {
// 	s.err = err
// }

/////////////////////////////////////////////////////////////////////////////////////////////////////////

func TestController_UserRegister(t *testing.T) {
	// type want struct {
	// 	code int
	// }
	// tests := []struct { // the array of structures
	// 	name string
	// 	want want
	// }{
	// 	{
	// 		name: "send empty body #1",
	// 		want: want{
	// 			code: 400,
	// 		},
	// 	},
	// }

	// for _, test := range tests {
	// 	t.Run(test.name, func(t *testing.T) {

	// 		t1 := &TestService{}
	// 		t1.SetError(service.ErrWrongNameOrPassword)
	// 		t1.Register(context.TODO()).Return("", "")

	// 		c := &Controller{t1}
	// 		c.Register()

	// 		request := httptest.NewRequest(http.MethodPost, "http://127.0.0.1:8080/", bytes.NewReader([]byte("https://practicum.yandex.ru/")))
	// 		request.Header.Set("Content-Type", "text/plain")
	// 		// create a new Recorder
	// 		w := httptest.NewRecorder()
	// 		//CreateShortURLHandler(w, request)

	// 		res := w.Result()
	// 		// check response code
	// 		assert.Equal(t, test.want.code, res.StatusCode)
	// 		// get and check the body

	// 		_, err := io.ReadAll(res.Body)
	// 		defer res.Body.Close() // we must use defer after io.ReadAll to avoid issues
	// 		// TODO: mb I should handle error from res.Body.Close() ?
	// 		require.NoError(t, err)

	// 		// todo new request
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}

	// 	})
	// }
}
