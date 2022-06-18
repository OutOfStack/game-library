package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

// AddSale records information about game being on sale
func (s *Storage) AddSale(ctx context.Context, cs CreateSale) (int64, error) {
	ctx, span := tracer.Start(ctx, "db.sale.addsale")
	defer span.End()

	const q = `insert into sales 
	(name, begin_date, end_date)
	values ($1, $2, $3)
	returning id`

	var lastInsertID int64
	err := s.DB.QueryRowContext(ctx, q, cs.Name, cs.BeginDate, cs.EndDate).Scan(&lastInsertID)
	if err != nil {
		return 0, fmt.Errorf("inserting sale %v: %w", cs, err)
	}

	return lastInsertID, nil
}

// GetSales returns all sales
func (s *Storage) GetSales(ctx context.Context) ([]Sale, error) {
	ctx, span := tracer.Start(ctx, "db.sale.getsales")
	defer span.End()

	sales := []Sale{}

	const q = `select id, name, begin_date, end_date
	from sales`
	if err := s.DB.SelectContext(ctx, &sales, q); err != nil {
		return nil, errors.Wrap(err, "selecting sales")
	}

	return sales, nil
}

// RetrieveSale returns sale by id
// If such entity does not exist returns error ErrNotFound{}
func (s *Storage) RetrieveSale(ctx context.Context, saleID int64) (*Sale, error) {
	ctx, span := tracer.Start(ctx, "db.sale.retrievesale")
	defer span.End()

	var sale Sale

	const q = `select id, name, begin_date, end_date
	from sales
	where id = $1
	limit 1`
	if err := s.DB.GetContext(ctx, &sale, q, saleID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound{"sale", saleID}
		}
		return nil, err
	}

	return &sale, nil
}
