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

// ListGameSales returns all sales for specified game
func ListGameSales(ctx context.Context, db *sqlx.DB, gameID int64) ([]GameSale, error) {
	_, err := Retrieve(ctx, db, gameID)
	if err != nil {
		return nil, err
	}

	gameSales := []GameSale{}

	const q = `select sg.game_id, sg.sale_id, s.name as sale, s.begin_date, s.end_date, sg.discount_percent
	from sales_games sg
	left join sales s on s.id = sg.sale_id
	left join games g on g.id = sg.game_id
	where sg.game_id = $1`
	if err := db.SelectContext(ctx, &gameSales, q, gameID); err != nil {
		return nil, errors.Wrapf(err, "selecting sales of game with id %q", gameID)
	}

	return gameSales, nil
}
