package postgres

import "errors"

func (r CalculationRepository) createSchemeIfNotExists() error {
	errs := make([]error, 0, 3)
	errs = append(errs, r.createRuleScheme())
	errs = append(errs, r.createCalculationQueueScheme())
	errs = append(errs, r.createCalculationScheme())

	return errors.Join(errs...)
}

func (r CalculationRepository) createRuleScheme() error {
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS rules (
			id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			match VARCHAR(255) UNIQUE NOT NULL,
			point REAL NOT NULL,
			way SMALLINT NOT NULL
		);
	`)
	return err
}

func (r CalculationRepository) createCalculationQueueScheme() error {
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS calculation_queue (
			order_number VARCHAR(255) NOT NULL,
			product_name VARCHAR(255) NOT NULL,
			price REAL NOT NULL
		);
		CREATE INDEX IF NOT EXISTS calculation_queue_order_number_idx ON calculation_queue (order_number);
		CREATE INDEX IF NOT EXISTS calculation_queue_product_name_idx ON calculation_queue (product_name);
	`)

	return err
}

func (r CalculationRepository) createCalculationScheme() error {
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS calculation (
			order_number VARCHAR(255) UNIQUE NOT NULL,
			points REAL NOT NULL,
			status SMALLINT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS calculation_order_number_idx ON calculation (order_number);
	`)
	return err
}
