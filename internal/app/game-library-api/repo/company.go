package repo

import (
	"context"
	"fmt"
)

// CreateCompany creates new company
func (s *Storage) CreateCompany(ctx context.Context, c Company) (int32, error) {
	ctx, span := tracer.Start(ctx, "db.company.create")
	defer span.End()

	var id int32
	const q = `
	INSERT INTO companies (name, igdb_id)
	VALUES ($1, $2)
	ON CONFLICT (igdb_id) DO NOTHING
	RETURNING id`

	if err := s.db.QueryRowContext(ctx, q, c.Name, c.IGDBID).Scan(&id); err != nil {
		return 0, fmt.Errorf("create company with name %s and igdb id %d: %v", c.Name, c.IGDBID, err)
	}

	return id, nil
}

// GetCompanies returns companies
func (s *Storage) GetCompanies(ctx context.Context) ([]Company, error) {
	ctx, span := tracer.Start(ctx, "db.company.get")
	defer span.End()

	companies := make([]Company, 0)
	const q = `
	SELECT id, name, igdb_id
	FROM companies`

	if err := s.db.SelectContext(ctx, &companies, q); err != nil {
		return nil, err
	}

	return companies, nil
}
