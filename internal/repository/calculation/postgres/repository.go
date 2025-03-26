package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5/pgconn"
	decl "github.com/vilasle/gophermart/internal/repository/calculation"
)

const (
	codeDuplicateKey = "23505"
	//	tables
	calculationQueueTable = "calculation_queue"
	calculationTable      = "calculation"
	ruleTable             = "rules"
)

type CalculationRepository struct {
	db *sql.DB
}

func NewCalculationRepository(conn *sql.DB) (CalculationRepository, error) {
	r := CalculationRepository{db: conn}

	if err := r.createSchemeIfNotExists(); err != nil {
		return CalculationRepository{}, err
	}
	return r, nil
}

func (r CalculationRepository) AddCalculationToQueue(ctx context.Context, dto ...decl.AddingCalculation) error {
	sb := sqlbuilder.InsertInto(calculationQueueTable).
		Cols("order_number", "product_name", "price")

	for _, v := range dto {
		sb.Values(v.OrderNumber, v.ProductName, v.Price)
	}

	txt, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	_, err := r.db.ExecContext(ctx, txt, args...)

	return getRepositoryError(err)
}

func (r CalculationRepository) ClearCalculationsQueue(ctx context.Context, dto decl.ClearingCalculationQueue) error {
	sb := sqlbuilder.DeleteFrom(calculationQueueTable)
	sb.Where(sb.Equal("order_number", dto.OrderNumber))

	txt, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	_, err := r.db.ExecContext(ctx, txt, args...)

	return getRepositoryError(err)
}

func (r CalculationRepository) AddCalculationResult(ctx context.Context, dto decl.AddCalculationResult) error {
	txt, args := sqlbuilder.InsertInto(calculationTable).
		Cols("order_number", "points", "status").
		Values(dto.OrderNumber, dto.Value, dto.Status).
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	_, err := r.db.ExecContext(ctx, txt, args...)

	return getRepositoryError(err)
}

func (r CalculationRepository) GetCalculationsQueue(ctx context.Context) ([]decl.CalculationQueueInfo, error) {
	sb := sqlbuilder.Select("order_number", "product_name", "price").From(calculationQueueTable)

	txt, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)

	rows, err := r.db.QueryContext(ctx, txt, args...)
	if err != nil {
		return nil, getRepositoryError(err)
	}
	defer rows.Close()

	result, err := prepareCalculationQueue(rows)

	return result, getRepositoryError(err)
}

func prepareCalculationQueue(rows *sql.Rows) ([]decl.CalculationQueueInfo, error) {
	result := make([]decl.CalculationQueueInfo, 0)
	for rows.Next() {
		var v decl.CalculationQueueInfo
		err := rows.Scan(&v.OrderNumber, &v.ProductName, &v.Price)
		if err != nil {
			return nil, getRepositoryError(err)
		}
		result = append(result, v)
	}
	return result, rows.Err()
}

func (r CalculationRepository) UpdateCalculationResult(ctx context.Context, dto decl.AddCalculationResult) error {
	sb := sqlbuilder.Update(calculationTable)
	sb.Set(
		sb.Equal("status", dto.Status),
		sb.Equal("points", dto.Value),
	)
	sb.Where(sb.Equal("order_number", dto.OrderNumber))

	txt, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)

	_, err := r.db.ExecContext(ctx, txt, args...)

	return getRepositoryError(err)
}

func (r CalculationRepository) Calculations(ctx context.Context, dto decl.CalculationFilter) ([]decl.CalculationInfo, error) {
	sb := sqlbuilder.Select("order_number", "points", "status").From(calculationTable)

	if dto.OrderNumber != "" {
		sb.Where(sb.Equal("order_number", dto.OrderNumber))
	}

	txt, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)

	rows, err := r.db.QueryContext(ctx, txt, args...)
	if err != nil {
		return nil, getRepositoryError(err)
	}
	defer rows.Close()

	result, err := prepareCalculations(rows)
	return result, getRepositoryError(err)
}

func prepareCalculations(rows *sql.Rows) ([]decl.CalculationInfo, error) {
	result := make([]decl.CalculationInfo, 0)
	for rows.Next() {
		var c decl.CalculationInfo
		err := rows.Scan(&c.OrderNumber, &c.Value, &c.Status)
		if err != nil {
			return nil, getRepositoryError(err)
		}
		result = append(result, c)
	}
	return result, rows.Err()
}

func (r CalculationRepository) AddRules(ctx context.Context, dto ...decl.AddingRule) (id int16, err error) {
	sp := sqlbuilder.InsertInto("rules").Cols("match", "point", "way").Returning("id")
	for _, v := range dto {
		sp.Values(v.Match, v.Point, v.CalculationType)
	}
	txt, args := sp.BuildWithFlavor(sqlbuilder.PostgreSQL)
	row := r.db.QueryRowContext(ctx, txt, args...)

	err = row.Scan(&id)
	return id, getRepositoryError(err)
}

func (r CalculationRepository) Rules(ctx context.Context, dto decl.RuleFilter) ([]decl.RuleInfo, error) {
	sp := sqlbuilder.Select("id", "match", "point", "way").From(ruleTable)

	if dto.ID > 0 {
		sp.Where(sp.Equal("id", dto.ID))
	}

	txt, args := sp.BuildWithFlavor(sqlbuilder.PostgreSQL)
	rows, err := r.db.QueryContext(ctx, txt, args...)
	if err != nil {
		return nil, getRepositoryError(err)
	}

	defer rows.Close()

	rules, err := prepareRulesInfo(rows)

	return rules, getRepositoryError(err)
}

func prepareRulesInfo(rows *sql.Rows) ([]decl.RuleInfo, error) {
	rules := make([]decl.RuleInfo, 0)
	for rows.Next() {
		var rule decl.RuleInfo
		err := rows.Scan(&rule.ID, &rule.Match, &rule.Point, &rule.CalculationType)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, rows.Err()
}

func getRepositoryError(err error) error {
	if err != nil {
		return err
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case codeDuplicateKey:
			return decl.ErrDuplicate
		}
	}
	return err
}
