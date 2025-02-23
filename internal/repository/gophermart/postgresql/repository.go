package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/huandu/go-sqlbuilder"
	mart "github.com/vilasle/gophermart/internal/repository/gophermart"
)

type PostgresqlGophermartRepository struct {
	db *sql.DB
}

func NewPostgresqlGophermartRepository(db *sql.DB) (*PostgresqlGophermartRepository, error) {
	r := &PostgresqlGophermartRepository{db: db}
	return r, r.createSchema()
}

// AuthorizationRepository
func (r PostgresqlGophermartRepository) AddUser(ctx context.Context, dto mart.AuthData) (mart.UserInfo, error) {
	sb := sqlbuilder.InsertInto("users").
		Cols("id", "login", "password").
		Values("gen_random_uuid()", dto.Login, dto.PasswordHash).
		Returning("id")

	txt, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	row := r.db.QueryRowContext(ctx, txt, args...)

	var id string
	err := row.Scan(&id)
	return mart.UserInfo{ID: id}, err
}

func (r PostgresqlGophermartRepository) CheckUser(ctx context.Context, dto mart.AuthData) (mart.UserInfo, error) {
	sb := sqlbuilder.Select("id", "login", "password").From("users")
	sb.Where(sb.Equal("login", dto.Login))

	txt, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	row := r.db.QueryRowContext(ctx, txt, args...)

	var id, login, password string
	err := row.Scan(&id, &login, &password)

	return mart.UserInfo{ID: id, Login: login, PasswordHash: password}, err
}

// WithdrawalRepository
func (r PostgresqlGophermartRepository) Expense(ctx context.Context, dto mart.WithdrawalRequest) error {
	v := dto.Sum
	if v > 0 {
		v = -v
	}

	sbCh := sqlbuilder.Select("SUM(sum)").From("transaction").GroupBy("user_id")
	sbCh.Where(sbCh.Equal("user_id", dto.UserID))

	txt1, args1 := sbCh.BuildWithFlavor(sqlbuilder.PostgreSQL)

	sbAdd := sqlbuilder.InsertInto("transactions").
		Cols("order_number", "user_id", "income", "sum", "created_at").
		Values(dto.OrderNumber, dto.UserID, false, v, "now()")

	txt2, args2 := sbAdd.BuildWithFlavor(sqlbuilder.PostgreSQL)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, txt1, args1...)

	var sum float64
	if err := row.Scan(&sum); err == nil {
		if sum < dto.Sum {
			return mart.ErrDoesNotEnoughPoints
		}
	}

	if _, err := tx.ExecContext(ctx, txt2, args2...); err != nil {
		return err
	}

	return tx.Commit()
}

func (r PostgresqlGophermartRepository) Income(ctx context.Context, dto mart.WithdrawalRequest) error {
	v := dto.Sum
	if v < 0 {
		v = -v
	}
	sbAdd := sqlbuilder.InsertInto("transactions").
		Cols("order_number", "user_id", "income", "sum", "created_at").
		Values(dto.OrderNumber, dto.UserID, true, v, "now()")

	txt, args := sbAdd.BuildWithFlavor(sqlbuilder.PostgreSQL)

	_, err := r.db.ExecContext(ctx, txt, args...)
	if err != nil {
		return err
	}
	return err

}

func (r PostgresqlGophermartRepository) Transactions(ctx context.Context, dto mart.TransactionRequest) ([]mart.Transaction, error) {
	sb := sqlbuilder.Select("order_number", "user_id", "income", "sum", "created_at").From("transactions")
	sb.Where(sb.Equal("user_id", dto.UserID))

	txt, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	rows, err := r.db.QueryContext(ctx, txt, args...)
	if err != nil {
		return nil, err
	}

	transactions := make([]mart.Transaction, 0)
	for rows.Next() {
		transaction := mart.Transaction{}
		err := rows.Scan(&transaction.OrderNumber, &transaction.UserID, &transaction.Income, &transaction.Sum, &transaction.CreatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

// OrderRepository
func (r PostgresqlGophermartRepository) Create(ctx context.Context, dto mart.OrderCreateRequest) error {
	sb := sqlbuilder.InsertInto("order").
		Cols("number", "user_id", "create_at", "status", "sum").
		Values(dto.Number, dto.UserID, "now()", mart.StatusNew, 0)

	txt, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	_, err := r.db.ExecContext(ctx, txt, args...)
	return err
}

func (r PostgresqlGophermartRepository) Update(ctx context.Context, dto mart.OrderUpdateRequest) error {
	sb := sqlbuilder.Update("order")
	sb.Set(fmt.Sprintf("status = %d", dto.Status))

	sb.Where(sb.Equal("number", dto.Number))

	txt, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	_, err := r.db.ExecContext(ctx, txt, args...)
	return err
}

func (r PostgresqlGophermartRepository) List(ctx context.Context, dto mart.OrderListRequest) ([]mart.OrderInfo, error) {
	sp := sqlbuilder.
		Select("number", "create_at", "status", "sum").
		From("order")
	sp.Where(sp.Equal("user_id", dto.UserID))

	txt, args := sp.BuildWithFlavor(sqlbuilder.PostgreSQL)
	rows, err := r.db.QueryContext(ctx, txt, args...)
	if err != nil {
		return nil, err
	}

	orders := make([]mart.OrderInfo, 0)
	for rows.Next() {
		order := mart.OrderInfo{}
		err := rows.Scan(&order.Number, &order.CreatedAt, &order.Status, &order.Accrual)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}
