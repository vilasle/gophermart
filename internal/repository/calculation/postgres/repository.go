package postgres

import (
	"context"
	"net/url"

	"github.com/huandu/go-sqlbuilder"
	decl "github.com/vilasle/gophermart/internal/repository/calculation"
)

type CalculationRepository struct {
	*connection
}

func NewPostgresCalculationRepository(url url.URL) (CalculationRepository, error) {
	conn, err := newConnection(url.String())
	if err != nil {
		//TODO could not create connection
		return CalculationRepository{}, err
	}

	r := CalculationRepository{connection: conn}

	err = r.createSchemeIfNotExists()
	if err != nil {
		//TODO could not create scheme
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
	_, err := r.Exec(ctx, txt, args...)
	//TODO check and wrap error
	return err
}

func (r CalculationRepository) SaveCalculationResult(ctx context.Context, dto decl.AddCalculationResult) error {
	txt, args := sqlbuilder.InsertInto("calculation").
		Cols("order_number", "points", "status").
		Values(dto.OrderNumber, dto.Value, dto.Status).
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	_, err := r.Exec(ctx, txt, args...)
	//TODO check and wrap error
	return err
}

func (r CalculationRepository) Calculations(ctx context.Context, dto decl.CalculationFilter) ([]decl.CalculationInfo, error) {
	//OrderNumber string

	sb := sqlbuilder.Select("order_number", "points", "status").From("calculation")

	if dto.OrderNumber != "" {
		sb.Where(sb.Equal("order_number", dto.OrderNumber))
	}

	txt, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)

	rows, err := r.Query(ctx, txt, args...)
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

func (r CalculationRepository) AddRules(ctx context.Context, dto ...decl.AddingRule) (id string, err error) {
	sp := sqlbuilder.InsertInto("rules").Cols("match", "point", "way")
	for _, v := range dto {
		sp.Values(v.Match, v.Point, v.CalculationType)
	}
	txt, args := sp.BuildWithFlavor(sqlbuilder.PostgreSQL)
	_, err = r.Exec(ctx, txt, args...)
	//TODO check and wrap error
	return id, err
}

func (r CalculationRepository) Rules(ctx context.Context, dto decl.RuleFilter) ([]decl.RuleInfo, error) {
	sp := sqlbuilder.Select("id", "match", "point", "way").From("rules")

	if dto.ID > 0 {
		sp.Where(sp.Equal("id", dto.ID))
	}

	txt, args := sp.BuildWithFlavor(sqlbuilder.PostgreSQL)
	rows, err := r.Query(ctx, txt, args...)
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
	}
	return rules, nil
}
