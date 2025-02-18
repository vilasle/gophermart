package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type connection struct {
	*pgxpool.Pool
}

func newConnection(url string) (*connection, error) {
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, err
	}

	baseCtx := context.Background()
	ctx, cancel := context.WithTimeout(baseCtx, 5*time.Second)
	defer cancel()

	return &connection{pool}, pool.Ping(ctx)
}

func (c *connection) ExecReturnOnlyError(sql string, args ...interface{}) error {
	baseCtx := context.Background()
	ctx, cancel := context.WithTimeout(baseCtx, 5*time.Second)
	defer cancel()

	_, err := c.Pool.Exec(ctx, sql, args...)
	return err
}

func (c *connection) Close() {
	c.Pool.Close()
}


