package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// CreateCompany creates new company
func (s *Storage) CreateCompany(ctx context.Context, c Company) (id int32, err error) {
	ctx, span := tracer.Start(ctx, "db.createCompany")
	defer span.End()

	const q = `
	INSERT INTO companies (name, igdb_id, created_at)
	VALUES ($1, $2, $3)
	ON CONFLICT (igdb_id) DO NOTHING
	RETURNING id`

	if err = s.db.QueryRowContext(ctx, q, c.Name, c.IGDBID, time.Now()).Scan(&id); err != nil {
		return 0, fmt.Errorf("create company with name %s and igdb id %d: %v", c.Name, c.IGDBID.Int64, err)
	}

	return id, nil
}

// GetCompanies returns companies
func (s *Storage) GetCompanies(ctx context.Context) (companies []Company, err error) {
	ctx, span := tracer.Start(ctx, "db.getCompanies")
	defer span.End()

	const q = `
	SELECT id, name, igdb_id
	FROM companies`

	if err = s.db.SelectContext(ctx, &companies, q); err != nil {
		return nil, err
	}

	return companies, nil
}

// GetCompanyIDByName returns company id by name
// If company does not exist returns ErrNotFound
func (s *Storage) GetCompanyIDByName(ctx context.Context, name string) (id int32, err error) {
	ctx, span := tracer.Start(ctx, "db.getCompanyIDByName")
	defer span.End()

	const q = `SELECT id
	FROM companies
	WHERE lower(name) = $1`

	if err = s.db.GetContext(ctx, &id, q, strings.ToLower(name)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNotFound[string]{Entity: "company", ID: name}
		}
		return 0, err
	}

	return id, nil
}

// GetTopDevelopers returns top developers by amount of games
func (s *Storage) GetTopDevelopers(ctx context.Context, limit int64) (companies []Company, err error) {
	ctx, span := tracer.Start(ctx, "db.getTopDevelopers")
	defer span.End()

	const q = `
	SELECT c.id, c.name, c.igdb_id
	FROM companies c
	JOIN (
		SELECT unnest(developers) AS company_id FROM games
	) AS g ON c.id = g.company_id
	GROUP BY c.id, c.name, c.igdb_id
	ORDER BY COUNT(*) DESC
	LIMIT $1`

	if err = s.db.SelectContext(ctx, &companies, q, limit); err != nil {
		return nil, err
	}

	return companies, nil
}

// GetTopPublishers returns top publishers by amount of games
func (s *Storage) GetTopPublishers(ctx context.Context, limit int64) (companies []Company, err error) {
	ctx, span := tracer.Start(ctx, "db.getTopPublishers")
	defer span.End()

	const q = `
	SELECT c.id, c.name, c.igdb_id
	FROM companies c
	JOIN (
		SELECT unnest(publishers) AS company_id FROM games
	) AS g ON c.id = g.company_id
	GROUP BY c.id, c.name, c.igdb_id
	ORDER BY COUNT(*) DESC
	LIMIT $1`

	if err = s.db.SelectContext(ctx, &companies, q, limit); err != nil {
		return nil, err
	}

	return companies, nil
}
