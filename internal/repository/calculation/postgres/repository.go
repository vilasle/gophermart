package postgres

import (
	"context"
	"database/sql"

	"github.com/huandu/go-sqlbuilder"
	decl "github.com/vilasle/gophermart/internal/repository/calculation"
)

type CalculationRepository struct {
	db *sql.DB
}

func NewPostgresCalculationRepository(conn *sql.DB) (CalculationRepository, error) {
	r := CalculationRepository{db: conn}

	if err := r.createSchemeIfNotExists(); err != nil {
		return CalculationRepository{}, err
	}
	return r, nil
}

func (r CalculationRepository) AddCalculationToQueue(ctx context.Context, dto ...decl.AddingCalculation) error {
	sb := sqlbuilder.InsertInto("calculation_queue").Cols("order_number", "product_name", "price")
	for _, v := range dto {
		sb.Values(v.OrderNumber, v.ProductName, v.Price)
	}
	txt, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	_, err := r.db.ExecContext(ctx, txt, args...)
	//TODO check and wrap error
	return err
}

func (r CalculationRepository) SaveCalculationResult(ctx context.Context, dto decl.AddCalculationResult) error {
	txt, args := sqlbuilder.InsertInto("calculation").
		Cols("order_number", "points", "status").
		Values(dto.OrderNumber, dto.Value, dto.Status).
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	_, err := r.db.ExecContext(ctx, txt, args...)
	//TODO check and wrap error
	return err
}

func (r CalculationRepository) Calculations(ctx context.Context, dto decl.CalculationFilter) ([]decl.CalculationInfo, error) {
	sb := sqlbuilder.Select("order_number", "points", "status").From("calculation")

	if dto.OrderNumber != "" {
		sb.Where(sb.Equal("order_number", dto.OrderNumber))
	}

	txt, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)

	rows, err := r.db.QueryContext(ctx, txt, args...)
	if err != nil {
		//TODO check and wrap error
		return nil, err
	}

	result := make([]decl.CalculationInfo, 0)
	for rows.Next() {
		var c decl.CalculationInfo
		err := rows.Scan(&c.OrderNumber, &c.Value, &c.Status)
		if err != nil {
			//TODO check and wrap error
			return nil, err
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
	//TODO check and wrap error
	return id, err
}

func (r CalculationRepository) Rules(ctx context.Context, dto decl.RuleFilter) ([]decl.RuleInfo, error) {
	sp := sqlbuilder.Select("id", "match", "point", "way").From("rules")

	if dto.ID > 0 {
		sp.Where(sp.Equal("id", dto.ID))
	}

	txt, args := sp.BuildWithFlavor(sqlbuilder.PostgreSQL)
	rows, err := r.db.QueryContext(ctx, txt, args...)
	if err != nil {
		return nil, err
	}
	rules := make([]decl.RuleInfo, 0)
	for rows.Next() {
		var rule decl.RuleInfo
		err := rows.Scan(&rule.ID, &rule.Match, &rule.Point, &rule.CalculationType)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}
