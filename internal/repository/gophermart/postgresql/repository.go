package postgresql

import (
	"context"
	"database/sql"

	mart "github.com/vilasle/gophermart/internal/repository/gophermart"
)

type PostgresqlGophermartRepository struct {
	db *sql.DB
}

func NewPostgresqlGophermartRepository(db *sql.DB) (*PostgresqlGophermartRepository, error) {
	r := &PostgresqlGophermartRepository{db: db}
	return r, r.createSchema()
}

func (r PostgresqlGophermartRepository) AddUser(ctx context.Context, dto mart.AuthData) (mart.UserInfo, error) {
	//TODO implement it
	panic("not implemented")
}

func (r PostgresqlGophermartRepository) CheckUser(ctx context.Context, dto mart.AuthData) (mart.UserInfo, error) {
	// TODO implement it
	panic("not implemented")
}

func (r PostgresqlGophermartRepository) Expense(ctx context.Context, dto mart.WithdrawalRequest) error {
	// TODO implement it
	panic("not implemented")
}

func (r PostgresqlGophermartRepository) Transactions(ctx context.Context, dto mart.TransactionRequest) ([]mart.Transaction, error) {
	// TODO implement it
	panic("not implemented")
}

func (r PostgresqlGophermartRepository) Create(ctx context.Context, dto mart.OrderCreateRequest) error {
	// TODO implement it
	panic("not implemented")
}

func (r PostgresqlGophermartRepository) Update(ctx context.Context, dto mart.OrderUpdateRequest) error {
	// TODO implement it
	panic("not implemented")
}

func (r PostgresqlGophermartRepository) List(ctx context.Context, dto mart.OrderListRequest) ([]mart.OrderInfo, error) {
	// TODO implement it
	panic("not implemented")
}
