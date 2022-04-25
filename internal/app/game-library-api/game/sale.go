package game

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/api/trace"
)

// AddSale records information about game being on sale
func AddSale(ctx context.Context, db *sqlx.DB, cs CreateSaleReq) (int64, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "sql.sale.addsale")
	defer span.End()

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

// GetSales returns all sales
func GetSales(ctx context.Context, db *sqlx.DB) ([]Sale, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "sql.sale.getsales")
	defer span.End()

	sales := []Sale{}

	const q = `select id, name, begin_date, end_date
	from sales`
	if err := db.SelectContext(ctx, &sales, q); err != nil {
		return nil, errors.Wrap(err, "selecting sales")
	}

	return sales, nil
}

// RetrieveSale returns sale by id
// If such entity does not exist returns error ErrNotFound{}
func RetrieveSale(ctx context.Context, db *sqlx.DB, saleID int64) (*Sale, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "sql.sale.retrievesale")
	defer span.End()

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
