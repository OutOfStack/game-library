package game

import (
	"context"
	"fmt"

	"cloud.google.com/go/civil"
	"github.com/OutOfStack/game-library/pkg/types"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// AddSale records information about game being on sale
func AddSale(ctx context.Context, db *sqlx.DB, ns NewSale, gameID int64) (*Sale, error) {
	beginDate, err := civil.ParseDate(ns.BeginDate)
	if err != nil {
		return nil, fmt.Errorf("parsing beginDate: %w", err)
	}
	endDate, err := civil.ParseDate(ns.EndDate)
	if err != nil {
		return nil, fmt.Errorf("parsing endDate: %w", err)
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

	s := Sale{
		ID:              lastInsertID,
		Name:            ns.Name,
		GameID:          gameID,
		BeginDate:       types.Date(beginDate),
		EndDate:         types.Date(endDate),
		DiscountPercent: ns.DiscountPercent,
	}

	return &s, nil
}

// ListSales returns all sales for specified game
func ListSales(ctx context.Context, db *sqlx.DB, gameID int64) ([]Sale, error) {
	sales := []Sale{}

	const q = `select id, name, game_id, begin_date, end_date, discount_percent
	from sales
	where game_id = $1`
	if err := db.SelectContext(ctx, &sales, q, gameID); err != nil {
		return nil, errors.Wrap(err, "selecting sales")
	}

	return sales, nil
}
