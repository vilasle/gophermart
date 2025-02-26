package postgresql

import "errors"

func (r PostgresqlGophermartRepository) createSchema() error {
	errs := make([]error, 0, 3)

	errs = append(errs, r.createUserTable())
	errs = append(errs, r.createOrderTable())
	errs = append(errs, r.createTransactionTable())

	return errors.Join(errs...)
}

func (r PostgresqlGophermartRepository) createUserTable() error {
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS "user" (
			id UUID PRIMARY KEY, 
			login VARCHAR(255) 
			UNIQUE NOT NULL, 
			password BYTEA
		);

		CREATE INDEX IF NOT EXISTS "user_login_idx" ON "user" (login);
	`)
	return err
}

func (r PostgresqlGophermartRepository) createOrderTable() error {
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS "order" (
			id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
			number VARCHAR(255) NOT NULL UNIQUE,
			user_id UUID NOT NULL,
			create_at TIMESTAMP NOT NULL,
			status SMALLINT NOT NULL,
			sum REAL NOT NULL,
			FOREIGN KEY (user_id) REFERENCES "user" (id)
		);
		CREATE INDEX IF NOT EXISTS "order_number_idx" ON "order" (number);
		CREATE INDEX IF NOT EXISTS "order_user_id_idx" ON "order" (user_id);
	`)
	return err
}

func (r PostgresqlGophermartRepository) createTransactionTable() error {
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS "transaction" (
			id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
			order_number VARCHAR(255) NOT NULL,
			user_id UUID NOT NULL,
			income BOOLEAN NOT NULL,
			sum REAL NOT NULL,
			create_at TIMESTAMP NOT NULL,
			FOREIGN KEY (user_id) REFERENCES "user" (id)
		);
		CREATE INDEX IF NOT EXISTS "transaction_user_id_idx" ON "transaction" (user_id);
	`)
	return err
}
