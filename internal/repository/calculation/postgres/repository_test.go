package postgres

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	repository "github.com/vilasle/gophermart/internal/repository/calculation"
)

func TestCreateTables(t *testing.T) {

	u, err := url.Parse("postgres://postgres:142543@172.17.0.2:5432/accrual")

	require.NoError(t, err)

	r, err := NewPostgresCalculationRepository(u)

	require.NoError(t, err)
	defer r.conn.Close()

	ctx := context.Background()
	//add rule
	id, err := r.AddRules(ctx, repository.AddingRule{
		Match:           "test",
		Point:           1,
		CalculationType: 1,
	})

	require.NoError(t, err)

	t.Log(id)

	//get rules
	rules, err := r.Rules(ctx, repository.RuleFilter{
		ID: id,
	})

	require.NoError(t, err)

	require.Equal(t, id, rules[0].ID)

	//add calculation to queue
	err = r.AddCalculationToQueue(ctx, repository.AddingCalculation{
		OrderNumber: "123456",
		ProductName: "test",
		Price:       1332.32,
	})

	require.NoError(t, err)

	err = r.SaveCalculationResult(ctx, repository.AddCalculationResult{
		OrderNumber: "123456",
		Status:      1,
		Value:       123,
	})

	require.NoError(t, err)

	calcs, err := r.Calculations(ctx, repository.CalculationFilter{
		OrderNumber: "123456",
	})

	require.NoError(t, err)

	require.Equal(t, 1, len(calcs))

	//save calculation result

	//get calculations

}
