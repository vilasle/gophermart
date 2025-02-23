package accrual

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/vilasle/gophermart/internal/repository/gophermart"
	"github.com/vilasle/gophermart/internal/service"
)

func TestAccrualServiceHTTP_Accruals(t *testing.T) {
	type args struct {
		ctx context.Context
		dto service.AccrualsFilterRequest
	}

	type mockSetting struct {
		dtoIn  gophermart.AccrualRequest
		dtoOut gophermart.AccrualInfo
		errOut error
		setup  func(*MockAccrualRepository, context.Context, gophermart.AccrualRequest, gophermart.AccrualInfo, error)
	}

	tests := []struct {
		name        string
		args        args
		mockSetting mockSetting
		want        service.AccrualsInfo
		wantErr     bool
	}{
		{
			name: "empty dto",
			args: args{
				ctx: context.Background(),
				dto: service.AccrualsFilterRequest{},
			},
			want:    service.AccrualsInfo{},
			wantErr: true,
			mockSetting: mockSetting{
				dtoIn:  gophermart.AccrualRequest{},
				dtoOut: gophermart.AccrualInfo{},
				errOut: nil,
				setup: func(mar *MockAccrualRepository, ctx context.Context, ar gophermart.AccrualRequest, ai gophermart.AccrualInfo, err error) {
				},
			},
		},
		{
			name: "repository raise error",
			args: args{
				ctx: context.Background(),
				dto: service.AccrualsFilterRequest{Number: "1234567890"},
			},
			mockSetting: mockSetting{
				dtoIn: gophermart.AccrualRequest{
					OrderNumber: "1234567890",
				},
				dtoOut: gophermart.AccrualInfo{},
				errOut: errors.New("repository error"),
				setup: func(mar *MockAccrualRepository, ctx context.Context, ar gophermart.AccrualRequest, ai gophermart.AccrualInfo, err error) {
					mar.EXPECT().AccrualByOrder(ctx, ar).Return(ai, err)
				},
			},
			want:    service.AccrualsInfo{},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				dto: service.AccrualsFilterRequest{Number: "1234567890"},
			},
			want: service.AccrualsInfo{
				OrderNumber: "1234567890",
				Status:      "PROCESSED",
				Accrual:     100,
			},
			mockSetting: mockSetting{
				dtoIn: gophermart.AccrualRequest{
					OrderNumber: "1234567890",
				},
				dtoOut: gophermart.AccrualInfo{
					Number:  "1234567890",
					Status:  "PROCESSED",
					Accrual: 100,
				},
				errOut: nil,
				setup: func(mar *MockAccrualRepository, ctx context.Context, ar gophermart.AccrualRequest, ai gophermart.AccrualInfo, err error) {
					mar.EXPECT().AccrualByOrder(ctx, ar).Return(ai, err)
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rep := NewMockAccrualRepository(ctrl)

			tt.mockSetting.setup(rep, tt.args.ctx, tt.mockSetting.dtoIn, tt.mockSetting.dtoOut, tt.mockSetting.errOut)

			s := NewAccrualService(rep)

			got, err := s.Accruals(tt.args.ctx, tt.args.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("AccrualServiceHTTP.Accruals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AccrualServiceHTTP.Accruals() = %v, want %v", got, tt.want)
			}
		})
	}
}
