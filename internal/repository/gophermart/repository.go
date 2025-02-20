package gophermart

import "context"

type AuthorizationRepository interface {
	AddUser(context.Context, AuthData) (UserInfo, error)
	CheckUser(context.Context, AuthData) (UserInfo, error)
}
