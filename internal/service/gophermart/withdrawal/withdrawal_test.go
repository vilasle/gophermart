package withdrawal

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vilasle/gophermart/internal/repository/gophermart"
	"github.com/vilasle/gophermart/internal/service"
)

func TestWithdrawalService_Withdraw(t *testing.T) {
	type args struct {
		ctx context.Context
		dto service.WithdrawalRequest
	}

	type mockSetting struct {
		dtoIn  gophermart.WithdrawalRequest
		errOut error
		setup  func(*MockWithdrawalRepository, context.Context, gophermart.WithdrawalRequest, error)
	}

	repErr := errors.New("repository error")

	tests := []struct {
		name        string
		args        args
		mockSetting mockSetting
		err         error
	}{
		{
			name: "invalid format",
			args: args{
				ctx: context.Background(),
				dto: service.WithdrawalRequest{
					UserID:      "",
					OrderNumber: "",
					Sum:         0,
				},
			},
			mockSetting: mockSetting{
				dtoIn:  gophermart.WithdrawalRequest{},
				errOut: nil,
				setup:  func(m *MockWithdrawalRepository, ctx context.Context, dto gophermart.WithdrawalRequest, err error) {},
			},
			err: service.ErrInvalidFormat,
		},
		{
			name: "not enough points",
			args: args{
				ctx: context.Background(),
				dto: service.WithdrawalRequest{
					UserID:      "123456",
					OrderNumber: "954323",
					Sum:         100,
				},
			},
			mockSetting: mockSetting{
				dtoIn: gophermart.WithdrawalRequest{
					UserID:      "123456",
					OrderNumber: "954323",
					Sum:         100,
				},
				errOut: gophermart.ErrNotEnoughPoints,
				setup: func(m *MockWithdrawalRepository, ctx context.Context, dto gophermart.WithdrawalRequest, err error) {
					m.EXPECT().Expense(ctx, dto).Return(err)
				},
			},
			err: service.ErrNotEnoughPoints,
		},
		{
			name: "unknown repository error",
			args: args{
				ctx: context.Background(),
				dto: service.WithdrawalRequest{
					UserID:      "123456",
					OrderNumber: "954323",
					Sum:         100,
				},
			},
			mockSetting: mockSetting{
				dtoIn: gophermart.WithdrawalRequest{
					UserID:      "123456",
					OrderNumber: "954323",
					Sum:         100,
				},
				errOut: repErr,
				setup: func(m *MockWithdrawalRepository, ctx context.Context, dto gophermart.WithdrawalRequest, err error) {
					m.EXPECT().Expense(ctx, dto).Return(err)
				},
			},
			err: repErr,
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				dto: service.WithdrawalRequest{
					UserID:      "123456",
					OrderNumber: "954323",
					Sum:         100,
				},
			},
			mockSetting: mockSetting{
				dtoIn: gophermart.WithdrawalRequest{
					UserID:      "123456",
					OrderNumber: "954323",
					Sum:         100,
				},
				errOut: nil,
				setup: func(m *MockWithdrawalRepository, ctx context.Context, dto gophermart.WithdrawalRequest, err error) {
					m.EXPECT().Expense(ctx, dto).Return(err)
				},
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := NewMockWithdrawalRepository(ctrl)
			tt.mockSetting.setup(mock, tt.args.ctx, tt.mockSetting.dtoIn, tt.mockSetting.errOut)

			s := NewWithdrawalService(mock)

			err := s.Withdraw(tt.args.ctx, tt.args.dto)

			if tt.err != nil {
				assert.Error(t, err, tt.err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestWithdrawalService_List(t *testing.T) {
	type args struct {
		ctx context.Context
		dto service.WithdrawalListRequest
	}

	type mockSetting struct {
		dtoIn  gophermart.TransactionRequest
		dtoOut []gophermart.Transaction
		errOut error
		setup  func(*MockWithdrawalRepository, context.Context, gophermart.TransactionRequest, []gophermart.Transaction, error)
	}

	type want struct {
		dto []service.WithdrawalInfo
		err error
	}

	repErr := errors.New("repository error")

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
				dto: service.WithdrawalListRequest{
					UserID: "",
				},
			},
			mockSetting: mockSetting{
				dtoIn:  gophermart.TransactionRequest{},
				errOut: nil,
				setup: func(m *MockWithdrawalRepository, ctx context.Context, dtoIn gophermart.TransactionRequest, dtoOut []gophermart.Transaction, err error) {
				},
			},
			want: want{
				err: service.ErrInvalidFormat,
			},
		},
		{
			name: "unknown repository error",
			args: args{
				ctx: context.Background(),
				dto: service.WithdrawalListRequest{
					UserID: "123456",
				},
			},
			mockSetting: mockSetting{
				dtoIn: gophermart.TransactionRequest{
					UserID: "123456",
				},
				errOut: repErr,
				setup: func(m *MockWithdrawalRepository, ctx context.Context, dtoIn gophermart.TransactionRequest, dtoOut []gophermart.Transaction, err error) {
					m.EXPECT().Transactions(ctx, dtoIn).Return(dtoOut, err)
				},
			},
			want: want{
				err: repErr,
			},
		},
		{
			name: "there are only income transactions",
			args: args{
				ctx: context.Background(),
				dto: service.WithdrawalListRequest{
					UserID: "123456",
				},
			},
			mockSetting: mockSetting{
				dtoIn: gophermart.TransactionRequest{
					UserID: "123456",
				},
				dtoOut: []gophermart.Transaction{
					{
						Income:      true,
						UserID:      "123456",
						OrderNumber: "954323",
						Sum:         100,
					},
					{
						Income:      true,
						UserID:      "123456",
						OrderNumber: "4323",
						Sum:         100,
					},
					{
						Income:      true,
						UserID:      "123456",
						OrderNumber: "34534523",
						Sum:         100,
					},
				},
				errOut: nil,
				setup: func(m *MockWithdrawalRepository, ctx context.Context, dtoIn gophermart.TransactionRequest, dtoOut []gophermart.Transaction, err error) {
					m.EXPECT().Transactions(ctx, dtoIn).Return(dtoOut, err)
				},
			},
			want: want{
				dto: []service.WithdrawalInfo{},
				err: nil,
			},
		},
		{
			name: "there are income and outcome transactions",
			args: args{
				ctx: context.Background(),
				dto: service.WithdrawalListRequest{
					UserID: "123456",
				},
			},
			mockSetting: mockSetting{
				dtoIn: gophermart.TransactionRequest{
					UserID: "123456",
				},
				dtoOut: []gophermart.Transaction{
					{
						Income:      true,
						UserID:      "123456",
						OrderNumber: "954323",
						Sum:         100,
					},
					{
						Income:      true,
						UserID:      "123456",
						OrderNumber: "4323",
						Sum:         100,
					},
					{
						Income:      true,
						UserID:      "123456",
						OrderNumber: "34534523",
						Sum:         100,
					},
					{
						Income:      false,
						UserID:      "123456",
						OrderNumber: "954323",
						Sum:         100,
					},
					{
						Income:      false,
						UserID:      "123456",
						OrderNumber: "4323",
						Sum:         100,
					},
					{
						Income:      false,
						UserID:      "123456",
						OrderNumber: "34534523",
						Sum:         100,
					},
				},
				errOut: nil,
				setup: func(m *MockWithdrawalRepository, ctx context.Context, dtoIn gophermart.TransactionRequest, dtoOut []gophermart.Transaction, err error) {
					m.EXPECT().Transactions(ctx, dtoIn).Return(dtoOut, err)
				},
			},
			want: want{
				dto: []service.WithdrawalInfo{
					{
						OrderNumber: "954323",
						Sum:         100,
					},
					{
						OrderNumber: "4323",
						Sum:         100,
					},
					{
						OrderNumber: "34534523",
						Sum:         100,
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
			mock := NewMockWithdrawalRepository(ctrl)
			tt.mockSetting.setup(mock, tt.args.ctx, tt.mockSetting.dtoIn, tt.mockSetting.dtoOut, tt.mockSetting.errOut)

			s := NewWithdrawalService(mock)

			got, err := s.List(tt.args.ctx, tt.args.dto)

			if tt.want.err != nil {
				assert.Error(t, err, tt.want.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.dto, got)
			}
		})
	}
}

func TestWithdrawalService_Balance(t *testing.T) {
	type args struct {
		ctx context.Context
		dto service.UserBalanceRequest
	}
	type mockSetting struct {
		dtoIn  gophermart.TransactionRequest
		dtoOut []gophermart.Transaction
		errOut error
		setup  func(m *MockWithdrawalRepository, ctx context.Context, dtoIn gophermart.TransactionRequest, dtoOut []gophermart.Transaction, err error)
	}
	type want struct {
		dto service.UserBalance
		err     error
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
				dto: service.UserBalanceRequest{
					UserID: "",
				},
			},
			mockSetting: mockSetting{
				dtoIn:  gophermart.TransactionRequest{},
				dtoOut: []gophermart.Transaction{},
				errOut: nil,
				setup: func(m *MockWithdrawalRepository, ctx context.Context, dtoIn gophermart.TransactionRequest, dtoOut []gophermart.Transaction, err error) {
				},
			},
			want: want{
				dto: service.UserBalance{},
				err:     service.ErrInvalidFormat,
			},
		},
		{
			name: "unknown repository error",
			args: args{
				ctx: context.Background(),
				dto: service.UserBalanceRequest{
					UserID: "12345",
				},
			},
			mockSetting: mockSetting{
				dtoIn: gophermart.TransactionRequest{
					UserID: "12345",
				},
				dtoOut: []gophermart.Transaction{},
				errOut: repError,
				setup: func(m *MockWithdrawalRepository, ctx context.Context, dtoIn gophermart.TransactionRequest, dtoOut []gophermart.Transaction, err error) {
					m.EXPECT().Transactions(ctx, dtoIn).Return(dtoOut, err)
				},
			},
			want: want{
				dto: service.UserBalance{},
				err:     repError,
			},
		},
		{
			name: "success got balance",
			args: args{
				ctx: context.Background(),
				dto: service.UserBalanceRequest{
					UserID: "12345",
				},
			},
			mockSetting: mockSetting{
				dtoIn: gophermart.TransactionRequest{
					UserID: "12345",
				},
				dtoOut: []gophermart.Transaction{
					{
						Income:      false,
						UserID:      "12345",
						OrderNumber: "12345",
						Sum:         100,
					},
					{
						Income:      true,
						UserID:      "12345",
						OrderNumber: "12345",
						Sum:         2000,
					},
					{
						Income:      false,
						UserID:      "12345",
						OrderNumber: "12345",
						Sum:         1000,
					},
					{
						Income:      true,
						UserID:      "12345",
						OrderNumber: "12345",
						Sum:         100,
					},
				},
				errOut: nil,
				setup: func(m *MockWithdrawalRepository, ctx context.Context, dtoIn gophermart.TransactionRequest, dtoOut []gophermart.Transaction, err error) {
					m.EXPECT().Transactions(ctx, dtoIn).Return(dtoOut, err)
				},
			},
			want: want{
				dto: service.UserBalance{
					Withdrawn: 1100,
					Current: 1000,
				},
				err:     nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := NewMockWithdrawalRepository(ctrl)

			tt.mockSetting.setup(mock, tt.args.ctx, tt.mockSetting.dtoIn, tt.mockSetting.dtoOut, tt.mockSetting.errOut)

			s := NewWithdrawalService(mock)

			got, err := s.Balance(tt.args.ctx, tt.args.dto)

			if tt.want.err != nil {
				assert.Error(t, err, tt.want.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.dto, got)
			}

		})
	}
}
