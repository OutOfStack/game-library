package game

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// AddSale records information about game being on sale
func AddSale(ctx context.Context, db *sqlx.DB, ns NewSale) (*GetSale, error) {
	const q = `insert into sales 
	(name, begin_date, end_date)
	values ($1, $2, $3)
	returning id`

	var lastInsertID int64
	err := db.QueryRowContext(ctx, q, ns.Name, ns.BeginDate, ns.EndDate).Scan(&lastInsertID)
	if err != nil {
		return nil, fmt.Errorf("inserting sale %v: %w", ns, err)
	}

	getSale := ns.mapToGetSale(lastInsertID)

	return getSale, nil
}

// ListSales returns all sales
func ListSales(ctx context.Context, db *sqlx.DB) ([]GetSale, error) {
	sales := []Sale{}

	const q = `select id, name, begin_date, end_date
	from sales`
	if err := db.SelectContext(ctx, &sales, q); err != nil {
		return nil, errors.Wrap(err, "selecting sales")
	}

	getSales := []GetSale{}
	for _, s := range sales {
		getSales = append(getSales, *s.mapToGetSale())
	}

	return getSales, nil
}

// RetrieveSale returns sale by id
func RetrieveSale(ctx context.Context, db *sqlx.DB, saleID int64) (*GetSale, error) {
	var sale Sale

	const q = `select id, name, begin_date, end_date
	from sales
	limit 1`
	if err := db.GetContext(ctx, &sale, q); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	getSale := sale.mapToGetSale()

	return getSale, nil
}
