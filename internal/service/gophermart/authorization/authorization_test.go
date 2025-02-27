package authorization

import (
	"context"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vilasle/gophermart/internal/repository/gophermart"
	"github.com/vilasle/gophermart/internal/service"
)

func TestAuthorizationService_Register(t *testing.T) {
	type args struct {
		ctx context.Context
		dto service.RegisterRequest
	}

	type mockSetting struct {
		dtoIn  gophermart.AuthData
		dtoOut gophermart.UserInfo
		errOut error
		setup  func(*MockAuthorizationRepository, context.Context, gophermart.AuthData, gophermart.UserInfo, error)
	}

	passwordHash := []byte{
		94, 136, 72, 152, 218, 40, 4, 113, 81, 208, 229, 111, 141, 198, 41, 39, 115, 96, 61, 13, 106, 171, 189,
		214, 42, 17, 239, 114, 29, 21, 66, 216,
	}

	repErr := errors.New("repository error")

	_, _ = repErr, passwordHash

	tests := []struct {
		name        string
		args        args
		mockSetting mockSetting
		want        service.UserInfo
		err         error
	}{
		{
			name: "invalid dto",
			args: args{
				ctx: context.Background(),
				dto: service.RegisterRequest{
					Login:    "",
					Password: "",
				},
			},
			mockSetting: mockSetting{
				dtoIn:  gophermart.AuthData{},
				dtoOut: gophermart.UserInfo{},
				errOut: nil,
				setup: func(m *MockAuthorizationRepository, ctx context.Context, dto gophermart.AuthData, result gophermart.UserInfo, err error) {
				},
			},
			want: service.UserInfo{},
			err:  service.ErrInvalidFormat,
		},
		{
			name: "duplicate error",
			args: args{
				ctx: context.Background(),
				dto: service.RegisterRequest{
					Login:    "login",
					Password: "password",
				},
			},
			mockSetting: mockSetting{
				dtoIn: gophermart.AuthData{
					Login:        "login",
					PasswordHash: passwordHash,
				},
				dtoOut: gophermart.UserInfo{},
				errOut: gophermart.ErrDuplicate,
				setup: func(m *MockAuthorizationRepository, ctx context.Context, dto gophermart.AuthData, result gophermart.UserInfo, err error) {
					m.EXPECT().AddUser(ctx, dto).Return(result, err)
				},
			},
			want: service.UserInfo{},
			err:  service.ErrDuplicate,
		},
		{
			name: "another repository error",
			args: args{
				ctx: context.Background(),
				dto: service.RegisterRequest{
					Login:    "login",
					Password: "password",
				},
			},
			mockSetting: mockSetting{
				dtoIn: gophermart.AuthData{
					Login:        "login",
					PasswordHash: passwordHash,
				},
				dtoOut: gophermart.UserInfo{},
				errOut: repErr,
				setup: func(m *MockAuthorizationRepository, ctx context.Context, dto gophermart.AuthData, result gophermart.UserInfo, err error) {
					m.EXPECT().AddUser(ctx, dto).Return(result, err)
				},
			},
			want: service.UserInfo{},
			err:  repErr,
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				dto: service.RegisterRequest{
					Login:    "login",
					Password: "password",
				},
			},
			mockSetting: mockSetting{
				dtoIn: gophermart.AuthData{
					Login:        "login",
					PasswordHash: passwordHash,
				},
				dtoOut: gophermart.UserInfo{
					ID: "1234567890",
				},
				errOut: nil,
				setup: func(m *MockAuthorizationRepository, ctx context.Context, dto gophermart.AuthData, result gophermart.UserInfo, err error) {
					m.EXPECT().AddUser(ctx, dto).Return(result, err)
				},
			},
			want: service.UserInfo{
				ID: "1234567890",
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := NewMockAuthorizationRepository(ctrl)

			tt.mockSetting.setup(mock, tt.args.ctx, tt.mockSetting.dtoIn, tt.mockSetting.dtoOut, tt.mockSetting.errOut)

			svc := NewAuthorizationService(mock)

			got, err := svc.Register(tt.args.ctx, tt.args.dto)

			if tt.err != nil && errors.Is(err, tt.err) {
				assert.ErrorIs(t, err, tt.err)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAuthorizationService_Authorize(t *testing.T) {
	type args struct {
		ctx context.Context
		dto service.AuthorizeRequest
	}

	type mockSetting struct {
		dtoIn  gophermart.AuthData
		dtoOut gophermart.UserInfo
		errOut error
		setup  func(*MockAuthorizationRepository, context.Context, gophermart.AuthData, gophermart.UserInfo, error)
	}

	passwordHash := []byte{
		94, 136, 72, 152, 218, 40, 4, 113, 81, 208, 229, 111, 141, 198, 41, 39, 115, 96, 61, 13, 106, 171, 189,
		214, 42, 17, 239, 114, 29, 21, 66, 216,
	}

	repErr := errors.New("repository error")

	_, _ = repErr, passwordHash

	tests := []struct {
		name        string
		args        args
		mockSetting mockSetting
		want        service.UserInfo
		err         error
	}{
		{
			name: "invalid dto",
			args: args{
				ctx: context.Background(),
				dto: service.AuthorizeRequest{
					Login:    "",
					Password: "",
				},
			},
			mockSetting: mockSetting{
				dtoIn:  gophermart.AuthData{},
				dtoOut: gophermart.UserInfo{},
				errOut: nil,
				setup: func(m *MockAuthorizationRepository, ctx context.Context, dto gophermart.AuthData, result gophermart.UserInfo, err error) {
				},
			},
			want: service.UserInfo{},
			err:  service.ErrInvalidFormat,
		},
		{
			name: "empty result error",
			args: args{
				ctx: context.Background(),
				dto: service.AuthorizeRequest{
					Login:    "login",
					Password: "password",
				},
			},
			mockSetting: mockSetting{
				dtoIn: gophermart.AuthData{
					Login:        "login",
				},
				dtoOut: gophermart.UserInfo{},
				errOut: gophermart.ErrEmptyResult,
				setup: func(m *MockAuthorizationRepository, ctx context.Context, dto gophermart.AuthData, result gophermart.UserInfo, err error) {
					m.EXPECT().CheckUser(ctx, dto).Return(result, err)
				},
			},
			want: service.UserInfo{},
			err:  service.ErrWrongNameOrPassword,
		},
		{
			name: "another repository error",
			args: args{
				ctx: context.Background(),
				dto: service.AuthorizeRequest{
					Login:    "login",
					Password: "password",
				},
			},
			mockSetting: mockSetting{
				dtoIn: gophermart.AuthData{
					Login:        "login",
				},
				dtoOut: gophermart.UserInfo{},
				errOut: repErr,
				setup: func(m *MockAuthorizationRepository, ctx context.Context, dto gophermart.AuthData, result gophermart.UserInfo, err error) {
					m.EXPECT().CheckUser(ctx, dto).Return(result, err)
				},
			},
			want: service.UserInfo{},
			err:  repErr,
		},
		{
			name: "wrong password",
			args: args{
				ctx: context.Background(),
				dto: service.AuthorizeRequest{
					Login:    "login",
					Password: "password",
				},
			},
			mockSetting: mockSetting{
				dtoIn: gophermart.AuthData{
					Login:        "login",
				},
				dtoOut: gophermart.UserInfo{
					ID:           "1234567890",
					PasswordHash: passwordHash[15:],
				},
				errOut: nil,
				setup: func(m *MockAuthorizationRepository, ctx context.Context, dto gophermart.AuthData, result gophermart.UserInfo, err error) {
					m.EXPECT().CheckUser(ctx, dto).Return(result, err)
				},
			},
			want: service.UserInfo{
				ID: "1234567890",
			},
			err: service.ErrWrongNameOrPassword,
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				dto: service.AuthorizeRequest{
					Login:    "login",
					Password: "password",
				},
			},
			mockSetting: mockSetting{
				dtoIn: gophermart.AuthData{
					Login:        "login",
				},
				dtoOut: gophermart.UserInfo{
					ID:           "1234567890",
					PasswordHash: passwordHash,
				},
				errOut: nil,
				setup: func(m *MockAuthorizationRepository, ctx context.Context, dto gophermart.AuthData, result gophermart.UserInfo, err error) {
					m.EXPECT().CheckUser(ctx, dto).Return(result, err)
				},
			},
			want: service.UserInfo{
				ID: "1234567890",
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := NewMockAuthorizationRepository(ctrl)

			tt.mockSetting.setup(mock, tt.args.ctx, tt.mockSetting.dtoIn, tt.mockSetting.dtoOut, tt.mockSetting.errOut)

			svc := NewAuthorizationService(mock)

			got, err := svc.Authorize(tt.args.ctx, tt.args.dto)

			if tt.err != nil && errors.Is(err, tt.err) {
				assert.ErrorIs(t, err, tt.err)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
