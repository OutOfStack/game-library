package game

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// AddSale records information about game being on sale
func AddSale(ctx context.Context, db *sqlx.DB, ns NewSale, gameID int64) (*GetSale, error) {
	_, err := Retrieve(ctx, db, gameID)
	if err != nil {
		return nil, err
	}

	const q = `insert into sales 
	(name, game_id, begin_date, end_date, discount_percent)
	values ($1, $2, $3, $4, $5)
	returning id`

	var lastInsertID int64
	err = db.QueryRowContext(ctx, q, ns.Name, gameID, ns.BeginDate, ns.EndDate, ns.DiscountPercent).Scan(&lastInsertID)
	if err != nil {
		return nil, fmt.Errorf("inserting sale %v: %w", ns, err)
	}

	getSale := ns.mapToGetSale(lastInsertID, gameID)

	return getSale, nil
}

// ListSales returns all sales for specified game
func ListSales(ctx context.Context, db *sqlx.DB, gameID int64) ([]GetSale, error) {
	sales := []Sale{}

	const q = `select id, name, game_id, begin_date, end_date, discount_percent
	from sales
	where game_id = $1`
	if err := db.SelectContext(ctx, &sales, q, gameID); err != nil {
		return nil, errors.Wrap(err, "selecting sales")
	}

	getSales := []GetSale{}
	for _, s := range sales {
		getSales = append(getSales, *s.mapToGetSale())
	}

	return getSales, nil
}
