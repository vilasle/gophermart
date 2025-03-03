package postgresql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5/pgconn"
	mart "github.com/vilasle/gophermart/internal/repository/gophermart"
)

const (
	codeDuplicateKey = "23505"
)

type PostgresqlGophermartRepository struct {
	db *sql.DB
}

func NewPostgresqlGophermartRepository(db *sql.DB) (PostgresqlGophermartRepository, error) {
	r := PostgresqlGophermartRepository{db: db}
	return r, r.createSchema()
}

// AuthorizationRepository
func (r PostgresqlGophermartRepository) AddUser(ctx context.Context, dto mart.AuthData) (mart.UserInfo, error) {
	sb := sqlbuilder.InsertInto(`"user"`).Cols("id", "login", "password")
	sb.Values(sqlbuilder.Raw("gen_random_uuid()"), dto.Login, dto.PasswordHash)
	sb.Returning("id")

	txt, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	row := r.db.QueryRowContext(ctx, txt, args...)

	var id string
	err := row.Scan(&id)
	return mart.UserInfo{ID: id}, getRepositoryError(err)
}

func (r PostgresqlGophermartRepository) CheckUser(ctx context.Context, dto mart.AuthData) (mart.UserInfo, error) {
	sb := sqlbuilder.Select("id", "login", "password").From(`"user"`)
	sb.Where(sb.Equal("login", dto.Login))

	txt, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	row := r.db.QueryRowContext(ctx, txt, args...)

	var id, login string
	var password []byte
	err := row.Scan(&id, &login, &password)

	return mart.UserInfo{ID: id, Login: login, PasswordHash: password}, getRepositoryError(err)
}

func (r PostgresqlGophermartRepository) CheckUserByID(ctx context.Context, reqID string) (mart.UserInfo, error) {
	sb := sqlbuilder.Select("id", "login", "password").From(`"user"`)
	sb.Where(sb.Equal("id", reqID))

	txt, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	row := r.db.QueryRowContext(ctx, txt, args...)

	var id, login string
	var password []byte

	err := row.Scan(&id, &login, &password)

	return mart.UserInfo{ID: id, Login: login, PasswordHash: password}, getRepositoryError(err)
}

// WithdrawalRepository
func (r PostgresqlGophermartRepository) Expense(ctx context.Context, dto mart.WithdrawalRequest) error {
	v := dto.Sum
	if v > 0 {
		v = -v
	}

	sbCh := sqlbuilder.Select("SUM(sum)").From(`"transaction"`).GroupBy("user_id")
	sbCh.Where(sbCh.Equal("user_id", dto.UserID))

	txt1, args1 := sbCh.BuildWithFlavor(sqlbuilder.PostgreSQL)

	sbAdd := sqlbuilder.InsertInto(`"transaction"`).
		Cols("order_number", "user_id", "income", "sum", "created_at").
		Values(dto.OrderNumber, dto.UserID, false, v, sqlbuilder.Raw("now()"))

	txt2, args2 := sbAdd.BuildWithFlavor(sqlbuilder.PostgreSQL)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, txt1, args1...)

	var sum float64
	if err := row.Scan(&sum); err == nil && sum < dto.Sum {
		return mart.ErrNotEnoughPoints
	} else if err == sql.ErrNoRows {
		return mart.ErrNotEnoughPoints
	}

	if _, err := tx.ExecContext(ctx, txt2, args2...); err != nil {
		return getRepositoryError(err)
	}

	return tx.Commit()
}

func (r PostgresqlGophermartRepository) Income(ctx context.Context, dto mart.WithdrawalRequest) error {
	v := dto.Sum
	if v < 0 {
		v = -v
	}
	sbAdd := sqlbuilder.InsertInto(`"transaction"`).
		Cols("order_number", "user_id", "income", "sum", "created_at").
		Values(dto.OrderNumber, dto.UserID, true, v, sqlbuilder.Raw("now()"))

	txt, args := sbAdd.BuildWithFlavor(sqlbuilder.PostgreSQL)

	_, err := r.db.ExecContext(ctx, txt, args...)
	return getRepositoryError(err)

}

func (r PostgresqlGophermartRepository) Transactions(ctx context.Context, dto mart.TransactionRequest) ([]mart.Transaction, error) {
	sb := sqlbuilder.Select("order_number", "user_id", "income", "sum", "created_at").
		From(`"transaction"`)
	sb.Where(sb.Equal("user_id", dto.UserID))

	txt, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	rows, err := r.db.QueryContext(ctx, txt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions, err := scanTransactions(rows)

	return transactions, getRepositoryError(err)
}

func scanTransactions(rows *sql.Rows) ([]mart.Transaction, error) {
	transactions := make([]mart.Transaction, 0)
	for rows.Next() {
		transaction := mart.Transaction{}
		if err := rows.Scan(&transaction.OrderNumber, &transaction.UserID, &transaction.Income,
			&transaction.Sum, &transaction.CreatedAt); err != nil {

			return nil, getRepositoryError(err)
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

// OrderRepository
func (r PostgresqlGophermartRepository) Create(ctx context.Context, dto mart.OrderCreateRequest) error {
	sb := sqlbuilder.InsertInto(`"order"`).
		Cols("number", "user_id", "created_at", "status", "sum").
		Values(dto.Number, dto.UserID, sqlbuilder.Raw("now()"), mart.StatusNew, 0)

	txt, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	_, err := r.db.ExecContext(ctx, txt, args...)
	return getRepositoryError(err)
}

func (r PostgresqlGophermartRepository) Update(ctx context.Context, dto mart.OrderUpdateRequest) error {
	sb := sqlbuilder.Update(`"order"`)
	sb.Set(
		sb.Equal("status", dto.Status),
		sb.Equal("sum", dto.Accrual),
	)

	sb.Where(sb.Equal("number", dto.Number))

	txt, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	_, err := r.db.ExecContext(ctx, txt, args...)
	return getRepositoryError(err)
}

func (r PostgresqlGophermartRepository) List(ctx context.Context, dto mart.OrderListRequest) ([]mart.OrderInfo, error) {
	sp := sqlbuilder.
		Select("user_id", "number", "created_at", "status", "sum").
		From(`"order"`)

	if len(dto.OrderNumber) > 0 {
		sp.Where(sp.Equal("number", dto.OrderNumber))
	}

	if len(dto.UserID) > 0 {
		sp.Where(sp.Equal("user_id", dto.UserID))
	}

	if dto.Status > 0 {
		sp.Where(sp.Equal("status", dto.Status))
	}

	txt, args := sp.BuildWithFlavor(sqlbuilder.PostgreSQL)
	rows, err := r.db.QueryContext(ctx, txt, args...)
	if err != nil {
		return nil, getRepositoryError(err)
	}

	if rows.Err() != nil {
		return nil, getRepositoryError(rows.Err())
	}

	defer rows.Close()

	orders, err := scanAsOrdersInfo(rows)

	return orders, getRepositoryError(err)
}

func scanAsOrdersInfo(rows *sql.Rows) ([]mart.OrderInfo, error) {
	orders := make([]mart.OrderInfo, 0)
	for rows.Next() {
		order := mart.OrderInfo{}
		err := rows.Scan(&order.UserID, &order.Number, &order.CreatedAt, &order.Status, &order.Accrual)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func getRepositoryError(err error) error {
	if err == nil {
		return err
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case codeDuplicateKey:
			return mart.ErrDuplicate
		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		return mart.ErrEmptyResult
	}

	return err
}
