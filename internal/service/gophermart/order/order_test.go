package order

import (
	"context"
	"errors"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vilasle/gophermart/internal/repository/gophermart"
	"github.com/vilasle/gophermart/internal/service"
)

func TestOrderService_Register(t *testing.T) {
	type args struct {
		ctx context.Context
		dto service.RegisterOrderRequest
	}

	type mockSetting struct {
		//order rep List()
		dtoListIn  gophermart.OrderListRequest
		dtoListOut []gophermart.OrderInfo
		errListOut error
		setupList  func(m *MockOrderRepository, ctx context.Context, dtoIn gophermart.OrderListRequest, dtoOut []gophermart.OrderInfo, err error)

		//order rep Create()
		dtoCreateIn  gophermart.OrderCreateRequest
		errCreateOut error
		setupCreate  func(m *MockOrderRepository, ctx context.Context, dtoIn gophermart.OrderCreateRequest, err error)

		//order rep Update()
		dtoUpdateIn  gophermart.OrderUpdateRequest
		errUpdateOut error
		setupUpdate  func(m *MockOrderRepository, ctx context.Context, dtoIn gophermart.OrderUpdateRequest, err error)

		//accrual Accruals()
		dtoAccrualIn  service.AccrualsFilterRequest
		dtoAccrualOut service.AccrualsInfo
		errAccrualOut error
		setupAccrual  func(m *MockAccrualService, ctx context.Context, dtoIn service.AccrualsFilterRequest, dtoOut service.AccrualsInfo, err error)

		//withdrawal Income()
		dtoIncomeIn  gophermart.WithdrawalRequest
		errIncomeOut error
		setupIncome  func(m *MockWithdrawalRepository, ctx context.Context, dtoIn gophermart.WithdrawalRequest, err error)
	}

	type want struct {
		dto []service.OrderInfo
		err error
	}

	repError := errors.New("repository error")

	tests := []struct {
		name        string
		args        args
		mockSetting mockSetting
		want        want
		wait        time.Duration
	}{
		{
			name: "invalid format",
			args: args{
				ctx: context.Background(),
				dto: service.RegisterOrderRequest{},
			},
			mockSetting: mockSetting{
				setupList: func(m *MockOrderRepository, ctx context.Context, dtoIn gophermart.OrderListRequest, dtoOut []gophermart.OrderInfo, err error) {
				},
				setupCreate: func(m *MockOrderRepository, ctx context.Context, dtoIn gophermart.OrderCreateRequest, err error) {},
				setupUpdate: func(m *MockOrderRepository, ctx context.Context, dtoIn gophermart.OrderUpdateRequest, err error) {},
				setupAccrual: func(m *MockAccrualService, ctx context.Context, dtoIn service.AccrualsFilterRequest, dtoOut service.AccrualsInfo, err error) {
				},
				setupIncome: func(m *MockWithdrawalRepository, ctx context.Context, dtoIn gophermart.WithdrawalRequest, err error) {
				},
			},
			want: want{
				dto: []service.OrderInfo{},
				err: service.ErrInvalidFormat,
			},
		},
		{
			name: "repository List() error",
			args: args{
				ctx: context.Background(),
				dto: service.RegisterOrderRequest{
					Number: "1234567890",
					UserID: "1234567890",
				},
			},
			mockSetting: mockSetting{
				dtoListIn: gophermart.OrderListRequest{
					UserID:      "1234567890",
					OrderNumber: "1234567890",
				},
				dtoListOut: []gophermart.OrderInfo{},
				errListOut: repError,
				setupList: func(m *MockOrderRepository, ctx context.Context, dtoIn gophermart.OrderListRequest, dtoOut []gophermart.OrderInfo, err error) {
					m.EXPECT().List(ctx, dtoIn).Return(dtoOut, err)
				},
				setupCreate: func(m *MockOrderRepository, ctx context.Context, dtoIn gophermart.OrderCreateRequest, err error) {},
				setupUpdate: func(m *MockOrderRepository, ctx context.Context, dtoIn gophermart.OrderUpdateRequest, err error) {},
				setupAccrual: func(m *MockAccrualService, ctx context.Context, dtoIn service.AccrualsFilterRequest, dtoOut service.AccrualsInfo, err error) {
				},
				setupIncome: func(m *MockWithdrawalRepository, ctx context.Context, dtoIn gophermart.WithdrawalRequest, err error) {
				},
			},
			want: want{
				dto: []service.OrderInfo{},
				err: repError,
			},
		},
		{
			name: "repository List() there registered orders",
			args: args{
				ctx: context.Background(),
				dto: service.RegisterOrderRequest{
					Number: "1234567890",
					UserID: "1234567890",
				},
			},
			mockSetting: mockSetting{
				dtoListIn: gophermart.OrderListRequest{
					UserID:      "1234567890",
					OrderNumber: "1234567890",
				},
				dtoListOut: []gophermart.OrderInfo{
					{
						Number:  "1234567890",
						Status:  gophermart.StatusInvalid,
						Accrual: 0,
					},
				},
				errListOut: nil,
				setupList: func(m *MockOrderRepository, ctx context.Context, dtoIn gophermart.OrderListRequest, dtoOut []gophermart.OrderInfo, err error) {
					m.EXPECT().List(ctx, dtoIn).Return(dtoOut, err)
				},
				setupCreate: func(m *MockOrderRepository, ctx context.Context, dtoIn gophermart.OrderCreateRequest, err error) {},
				setupUpdate: func(m *MockOrderRepository, ctx context.Context, dtoIn gophermart.OrderUpdateRequest, err error) {},
				setupAccrual: func(m *MockAccrualService, ctx context.Context, dtoIn service.AccrualsFilterRequest, dtoOut service.AccrualsInfo, err error) {
				},
				setupIncome: func(m *MockWithdrawalRepository, ctx context.Context, dtoIn gophermart.WithdrawalRequest, err error) {
				},
			},
			want: want{
				dto: []service.OrderInfo{},
				err: service.ErrDuplicate,
			},
			wait: time.Millisecond * 500,
		},
		{
			name: "repository Create() error",
			args: args{
				ctx: context.Background(),
				dto: service.RegisterOrderRequest{
					Number: "1234567890",
					UserID: "1234567890",
				},
			},
			mockSetting: mockSetting{
				dtoListIn: gophermart.OrderListRequest{
					UserID:      "1234567890",
					OrderNumber: "1234567890",
				},
				dtoListOut: []gophermart.OrderInfo{},
				errListOut: nil,
				setupList: func(m *MockOrderRepository, ctx context.Context, dtoIn gophermart.OrderListRequest, dtoOut []gophermart.OrderInfo, err error) {
					m.EXPECT().List(ctx, dtoIn).Return(dtoOut, err)
				},

				dtoCreateIn: gophermart.OrderCreateRequest{
					UserID: "1234567890",
					Number: "1234567890",
				},
				errCreateOut: repError,
				setupCreate: func(m *MockOrderRepository, ctx context.Context, dtoIn gophermart.OrderCreateRequest, err error) {
					m.EXPECT().Create(ctx, dtoIn).Return(err)
				},
				setupUpdate: func(m *MockOrderRepository, ctx context.Context, dtoIn gophermart.OrderUpdateRequest, err error) {},
				setupAccrual: func(m *MockAccrualService, ctx context.Context, dtoIn service.AccrualsFilterRequest, dtoOut service.AccrualsInfo, err error) {
				},
				setupIncome: func(m *MockWithdrawalRepository, ctx context.Context, dtoIn gophermart.WithdrawalRequest, err error) {
				},
			},
			want: want{
				dto: []service.OrderInfo{},
				err: service.ErrDuplicate,
			},
			wait: time.Millisecond * 500,
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				dto: service.RegisterOrderRequest{
					Number: "1234567890",
					UserID: "1234567890",
				},
			},
			mockSetting: mockSetting{
				dtoListIn: gophermart.OrderListRequest{
					UserID:      "1234567890",
					OrderNumber: "1234567890",
				},
				dtoListOut: []gophermart.OrderInfo{},
				errListOut: nil,
				setupList: func(m *MockOrderRepository, ctx context.Context, dtoIn gophermart.OrderListRequest, dtoOut []gophermart.OrderInfo, err error) {
					m.EXPECT().List(ctx, dtoIn).Return(dtoOut, err)
				},

				dtoCreateIn: gophermart.OrderCreateRequest{
					UserID: "1234567890",
					Number: "1234567890",
				},
				errCreateOut: nil,
				setupCreate: func(m *MockOrderRepository, ctx context.Context, dtoIn gophermart.OrderCreateRequest, err error) {
					m.EXPECT().Create(ctx, dtoIn).Return(err)
				},
				dtoAccrualIn: service.AccrualsFilterRequest{
					Number: "1234567890",
				},
				dtoAccrualOut: service.AccrualsInfo{
					OrderNumber: "1234567890",
					Status:      "PROCESSED",
					Accrual:     100,
				},
				errAccrualOut: nil,
				setupAccrual: func(m *MockAccrualService, ctx context.Context, dtoIn service.AccrualsFilterRequest, dtoOut service.AccrualsInfo, err error) {
					m.EXPECT().Accruals(ctx, dtoIn).Return(dtoOut, err)
				},
				dtoUpdateIn: gophermart.OrderUpdateRequest{
					Number: "1234567890",
					Status: gophermart.StatusProcessed,
					UserID: "1234567890",
				},
				errUpdateOut: nil,
				setupUpdate: func(m *MockOrderRepository, ctx context.Context, dtoIn gophermart.OrderUpdateRequest, err error) {
					m.EXPECT().Update(gomock.Any(), dtoIn).Return(err)
				},
				dtoIncomeIn: gophermart.WithdrawalRequest{
					UserID:      "1234567890",
					OrderNumber: "1234567890",
					Sum:         100,
				},
				setupIncome: func(m *MockWithdrawalRepository, ctx context.Context, dtoIn gophermart.WithdrawalRequest, err error) {
					m.EXPECT().Income(gomock.Any(), dtoIn).Return(err)
				},
			},
			want: want{
				dto: []service.OrderInfo{},
				err: nil,
			},
			wait: time.Second * 5,
		},
		{
			name: "success with the second attempt",
			args: args{
				ctx: context.Background(),
				dto: service.RegisterOrderRequest{
					Number: "1234567890",
					UserID: "1234567890",
				},
			},
			mockSetting: mockSetting{
				dtoListIn: gophermart.OrderListRequest{
					UserID:      "1234567890",
					OrderNumber: "1234567890",
				},
				dtoListOut: []gophermart.OrderInfo{},
				errListOut: nil,
				setupList: func(m *MockOrderRepository, ctx context.Context, dtoIn gophermart.OrderListRequest, dtoOut []gophermart.OrderInfo, err error) {
					m.EXPECT().List(ctx, dtoIn).Return(dtoOut, err)
				},

				dtoCreateIn: gophermart.OrderCreateRequest{
					UserID: "1234567890",
					Number: "1234567890",
				},
				errCreateOut: nil,
				setupCreate: func(m *MockOrderRepository, ctx context.Context, dtoIn gophermart.OrderCreateRequest, err error) {
					m.EXPECT().Create(ctx, dtoIn).Return(err)
				},
				dtoAccrualIn: service.AccrualsFilterRequest{
					Number: "1234567890",
				},
				dtoAccrualOut: service.AccrualsInfo{
					OrderNumber: "1234567890",
					Status:      "PROCESSED",
					Accrual:     100,
				},
				errAccrualOut: nil,
				setupAccrual: func(m *MockAccrualService, ctx context.Context, dtoIn service.AccrualsFilterRequest, dtoOut service.AccrualsInfo, err error) {
					m.EXPECT().Accruals(ctx, dtoIn).Return(dtoOut, err)
				},
				dtoUpdateIn: gophermart.OrderUpdateRequest{
					Number: "1234567890",
					Status: gophermart.StatusProcessed,
					UserID: "1234567890",
				},
				errUpdateOut: nil,
				setupUpdate: func(m *MockOrderRepository, ctx context.Context, dtoIn gophermart.OrderUpdateRequest, err error) {
					m.EXPECT().Update(gomock.Any(), dtoIn).Return(err)
				},
				dtoIncomeIn: gophermart.WithdrawalRequest{
					UserID:      "1234567890",
					OrderNumber: "1234567890",
					Sum:         100,
				},
				setupIncome: func(m *MockWithdrawalRepository, ctx context.Context, dtoIn gophermart.WithdrawalRequest, err error) {
					m.EXPECT().Income(gomock.Any(), dtoIn).Return(err)
				},
			},
			want: want{
				dto: []service.OrderInfo{},
				err: nil,
			},
			wait: time.Second * 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repOrder := NewMockOrderRepository(ctrl)
			accSvc := NewMockAccrualService(ctrl)
			repTx := NewMockWithdrawalRepository(ctrl)

			tt.mockSetting.setupList(repOrder, tt.args.ctx, tt.mockSetting.dtoListIn, tt.mockSetting.dtoListOut, tt.mockSetting.errListOut)
			tt.mockSetting.setupCreate(repOrder, tt.args.ctx, tt.mockSetting.dtoCreateIn, tt.mockSetting.errCreateOut)
			tt.mockSetting.setupUpdate(repOrder, tt.args.ctx, tt.mockSetting.dtoUpdateIn, tt.mockSetting.errUpdateOut)
			tt.mockSetting.setupAccrual(accSvc, tt.args.ctx, tt.mockSetting.dtoAccrualIn, tt.mockSetting.dtoAccrualOut, tt.mockSetting.errAccrualOut)
			tt.mockSetting.setupIncome(repTx, tt.args.ctx, tt.mockSetting.dtoIncomeIn, tt.mockSetting.errIncomeOut)

			svc := NewOrderService(repOrder, accSvc, repTx)

			err := svc.Register(tt.args.ctx, tt.args.dto)

			if tt.want.err != nil {
				assert.Error(t, err, tt.want.err)
			} else {
				assert.NoError(t, err)
			}

			svc.Close()

			time.Sleep(tt.wait)
		})
	}
}

func TestOrderService_List(t *testing.T) {
	type args struct {
		ctx context.Context
		dto service.ListOrderRequest
	}

	type mockSetting struct {
		dtoIn  gophermart.OrderListRequest
		dtoOut []gophermart.OrderInfo
		errOut error
		setup  func(m *MockOrderRepository, ctx context.Context, dtoIn gophermart.OrderListRequest, dtoOut []gophermart.OrderInfo, err error)
	}

	type want struct {
		dto []service.OrderInfo
		err error
	}

	repError := errors.New("repository error")

	tests := []struct {
		name        string
		args        args
		mockSetting mockSetting
		want        want
	}{
		{
			name: "invalid format",
			args: args{
				ctx: context.Background(),
				dto: service.ListOrderRequest{
					UserID: "",
				},
			},
			mockSetting: mockSetting{
				setup: func(m *MockOrderRepository, ctx context.Context, dtoIn gophermart.OrderListRequest, dtoOut []gophermart.OrderInfo, err error) {
				},
			},
			want: want{
				dto: []service.OrderInfo{},
				err: service.ErrInvalidFormat,
			},
		},
		{
			name: "unknown repository error",
			args: args{
				ctx: context.Background(),
				dto: service.ListOrderRequest{
					UserID: "1234567",
				},
			},
			mockSetting: mockSetting{
				dtoIn: gophermart.OrderListRequest{
					UserID: "1234567",
				},
				dtoOut: []gophermart.OrderInfo{},
				errOut: repError,
				setup: func(m *MockOrderRepository, ctx context.Context, dtoIn gophermart.OrderListRequest, dtoOut []gophermart.OrderInfo, errOut error) {
					m.EXPECT().List(ctx, dtoIn).Return(dtoOut, errOut)
				},
			},
			want: want{
				dto: []service.OrderInfo{},
				err: repError,
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				dto: service.ListOrderRequest{
					UserID: "1234567",
				},
			},
			mockSetting: mockSetting{
				dtoIn: gophermart.OrderListRequest{
					UserID: "1234567",
				},
				dtoOut: []gophermart.OrderInfo{
					{
						Number:  "123456",
						Status:  gophermart.StatusInvalid,
						Accrual: 0,
					},
					{
						Number:  "65432",
						Status:  gophermart.StatusProcessed,
						Accrual: 100,
					},
				},
				errOut: nil,
				setup: func(m *MockOrderRepository, ctx context.Context, dtoIn gophermart.OrderListRequest, dtoOut []gophermart.OrderInfo, errOut error) {
					m.EXPECT().List(ctx, dtoIn).Return(dtoOut, errOut)
				},
			},
			want: want{
				dto: []service.OrderInfo{
					{
						Number:  "123456",
						Status:  "INVALID",
						Accrual: 0,
					},
					{
						Number:  "65432",
						Status:  "PROCESSED",
						Accrual: 100,
					},
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repOrder := NewMockOrderRepository(ctrl)
			accSvc := NewMockAccrualService(ctrl)
			repTx := NewMockWithdrawalRepository(ctrl)

			tt.mockSetting.setup(repOrder, tt.args.ctx, tt.mockSetting.dtoIn, tt.mockSetting.dtoOut, tt.mockSetting.errOut)

			svc := NewOrderService(repOrder, accSvc, repTx)

			got, err := svc.List(tt.args.ctx, tt.args.dto)

			if tt.want.err != nil {
				assert.Error(t, err, tt.want.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.dto, got)
			}

			svc.Close()

			time.Sleep(time.Microsecond * 500)
		})
	}
}
