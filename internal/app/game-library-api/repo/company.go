package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
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
		return 0, fmt.Errorf("create company with name %s and igdb id %d: %v", c.Name, c.IGDBID.Int64, err)
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

// GetCompanyIDByName returns company id by name
// If company does not exist returns ErrNotFound
func (s *Storage) GetCompanyIDByName(ctx context.Context, name string) (int32, error) {
	ctx, span := tracer.Start(ctx, "db.game.getcompanyidbyname")
	defer span.End()

	var id int32
	const q = `SELECT id
	FROM companies
	WHERE lower(name) = $1`

	if err := s.db.GetContext(ctx, &id, q, strings.ToLower(name)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNotFound[string]{Entity: "company", ID: name}
		}
		return 0, err
	}

	return id, nil
}
