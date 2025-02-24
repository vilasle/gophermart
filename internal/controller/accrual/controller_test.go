package controller

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vilasle/gophermart/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

// https://kovardin.ru/articles/go/testirovanie-http-hendlerov-v-go/
// https://www.youtube.com/watch?v=Mvw5fbHGJFw
// Смысл мока у меня: например, мокаю интерфейс сервиса, выставляю то, что примет мок и что он должен вернуть.
func TestHandler_orderInfo(t *testing.T) {
	type behaviour func(srv *MockAccrualService, ctx context.Context, dto service.AccrualsFilterRequest)
	// type behaviour func(srv *MockAccrualService, ctx context.Context, dto service.AccrualsFilterRequest) (service.AccrualsInfo, error)
	type args struct { // the input parameters of the testing method and a mock
		ctx  context.Context
		dto1 service.AccrualsFilterRequest
		dto2 service.AccrualsInfo
		behaviour
	}

	testTable := []struct {
		name string
		args args // input method data + mock
		//wantData service.AccrualsInfo
		wantData string // what I want from the handler to response
		wantErr  error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				dto1: service.AccrualsFilterRequest{
					Number: "12345678",
				},
				dto2: service.AccrualsInfo{
					OrderNumber: "12345678",
					Status:      "PROCESSED",
					Accrual:     500, //////////////// TODO праввильно ли что dto2 добавил как тогда REturn делать
				},
				behaviour: func(srv *MockAccrualService, ctx context.Context, dto1 service.AccrualsFilterRequest, dto2 service.AccrualsInfo) {
					srv.EXPECT().Accruals(ctx, dto1.Number).Return(service.AccrualsInfo{OrderNumber: "12345678", Status: "PROCESSED", Accrual: 500}, nil)

				},
			},
			wantData: `{"order": "12345678", "status": "PROCESSED", "accrual": 500}`,
			wantErr:  nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t) // контроллер мока сервиса
			defer ctrl.Finish()

			storageMock := NewMockAccrualService(ctrl) // create the service mock
			testCase.args.behaviour(storageMock, testCase.args.ctx, testCase.args.dto1)
			//handler :=  c
			//services := &repository.URLStorage{storageMock} //
			// init handler
			//handler := http.HandlerFunc(Controller.OrderInfo)
			// Test server
			//r: = chi.NewRouter()
			//r.POST("/api/orders/{number}"), handler.Accruals)

			w := httptest.NewRecorder() // fake response
			//req := httptest.NewRequest(http.MethodGet, "/api/orders/{number}", nil)
			//request := httptest.NewRequest(http.MethodPost, "http://127.0.0.1:8080/", bytes.NewReader([]byte("https://practicum.yandex.ru/")))
			res := w.Result()

			// Perform the request
			//handler := http.HandlerFunc(Controller.OrderInfo)
			//handler.ServeHTTP(w,req)// test request
			// Assert
			assert.Equal(t, http.StatusOK, w.Code)
			//assert.Equal(t, testCase.wantData, w.Body.String())
			assert.Equal(t, testCase.wantData, res.Body)
			/*testTable := []struct {
				name         string
				ctx          context.Context
				expectedCode int
				behaviour    // Mock behaviour
			}{
				{
					name: "success",
					ctx:  context.Background(),
					behaviour: func(rep *mock_repository.MockURLStorage, ctx context.Context) error {
						rep.EXPECT().Ping(ctx).Return(nil)
						return nil
					},
					expectedCode: http.StatusOK,
				},
			}

			for _, testCase := range testTable {
				t.Run(testCase.name, func(t *testing.T) {
					c := gomock.NewController(t)
					defer c.Finish()

					storageMock := mock_repository.NewMockURLStorage(c) // create the service mock
					testCase.behaviour(storageMock, testCase.ctx)

					services := &repository.URLStorage{storageMock} //
					handler := NewHand

				})
			}
			*/
		})
	}
}

/*func TestGet(t *testing.T) {
	// создаём контроллер
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// создаём объект-заглушку
	m := NewMockAccrualService(ctrl)

	// гарантируем, что заглушка
	// при вызове с аргументом "Key" вернёт "Value"
	in1 := service.AccrualsFilterRequest{} // TODO: на вход что передаю
	res1 := service.AccrualsInfo{}         // TODO: что в результате
	value := []byte("Value")
	m.EXPECT().Accruals(context.Background(), in1).Return(res1, nil)

	// какую функцию тестирую передав в неё объект-заглушку // TODO:

	val, err := с.OrderInfo(m, "Key")
	// и проверяем возвращаемые значения
	require.NoError(t, err)
	require.Equal(t, val, value)
}
*/
