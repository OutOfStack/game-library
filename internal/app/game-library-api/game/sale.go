package game

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// AddSale records information about game being on sale
func AddSale(ctx context.Context, db *sqlx.DB, cs CreateSale) (int64, error) {
	const q = `insert into sales 
	(name, begin_date, end_date)
	values ($1, $2, $3)
	returning id`

	var lastInsertID int64
	err := db.QueryRowContext(ctx, q, cs.Name, cs.BeginDate, cs.EndDate).Scan(&lastInsertID)
	if err != nil {
		return 0, fmt.Errorf("inserting sale %v: %w", cs, err)
	}

	return lastInsertID, nil
}

// ListSales returns all sales
func ListSales(ctx context.Context, db *sqlx.DB) ([]Sale, error) {
	sales := []Sale{}

	const q = `select id, name, begin_date, end_date
	from sales`
	if err := db.SelectContext(ctx, &sales, q); err != nil {
		return nil, errors.Wrap(err, "selecting sales")
	}

	return sales, nil
}

// RetrieveSale returns sale by id
func RetrieveSale(ctx context.Context, db *sqlx.DB, saleID int64) (*Sale, error) {
	var sale Sale

	const q = `select id, name, begin_date, end_date
	from sales
	where id = $1
	limit 1`
	if err := db.GetContext(ctx, &sale, q, saleID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound{"sale", saleID}
		}
		return nil, err
	}

	return &sale, nil
}
