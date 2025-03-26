package authorization

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"errors"

	"github.com/vilasle/gophermart/internal/repository/gophermart"
	"github.com/vilasle/gophermart/internal/service"
)

type AuthorizationService struct {
	rep gophermart.AuthorizationRepository
}

func NewAuthorizationService(rep gophermart.AuthorizationRepository) AuthorizationService {
	return AuthorizationService{rep: rep}
}

func (svc AuthorizationService) Register(ctx context.Context, dto service.RegisterRequest) (service.UserInfo, error) {
	if !checkFillingLoginPassword(dto.Login, dto.Password) {
		return service.UserInfo{}, service.ErrInvalidFormat
	}
	hash := svc.createHash(dto.Password)
	user := gophermart.AuthData{
		Login:        dto.Login,
		PasswordHash: hash[:],
	}

	result, err := svc.rep.AddUser(ctx, user)
	if err != nil {
		if errors.Is(err, gophermart.ErrDuplicate) {
			return service.UserInfo{}, service.ErrDuplicate
		}
		return service.UserInfo{}, err
	}

	return service.UserInfo{ID: result.ID}, nil
}

func (svc AuthorizationService) Authorize(ctx context.Context, dto service.AuthorizeRequest) (service.UserInfo, error) {
	if !checkFillingLoginPassword(dto.Login, dto.Password) {
		return service.UserInfo{}, service.ErrInvalidFormat
	}
	
	hash := svc.createHash(dto.Password)
	user := gophermart.AuthData{
		Login: dto.Login,
	}

	result, err := svc.rep.CheckUser(ctx, user)
	if err != nil {
		if errors.Is(err, gophermart.ErrEmptyResult) {
			return service.UserInfo{}, service.ErrWrongNameOrPassword
		}
		return service.UserInfo{}, err
	}

	passwordHash := []byte(result.PasswordHash)

	if ok := hmac.Equal(passwordHash, hash[:]); !ok {
		return service.UserInfo{}, service.ErrWrongNameOrPassword
	}

	return service.UserInfo{
		ID: result.ID,
	}, nil
}

func checkFillingLoginPassword(login, password string) bool {
	if login == "" || password == "" {
		return false
	}
	return true
}

func (svc AuthorizationService) CheckByUserID(ctx context.Context, id string) error {
	if id == "" {
		return service.ErrInvalidFormat
	}

	_, err := svc.rep.CheckUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, gophermart.ErrEmptyResult) {
			return service.ErrEntityDoesNotExists
		}
		return err
	}
	return nil
}

func (svc AuthorizationService) createHash(password string) [32]byte {
	hash := sha256.Sum256([]byte(password))
	return hash
}
