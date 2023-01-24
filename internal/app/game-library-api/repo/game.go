package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/OutOfStack/game-library/pkg/types"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.opentelemetry.io/otel"
)

// Storage provides required dependencies for repository
type Storage struct {
	db *sqlx.DB
}

// New creates new Storage
func New(db *sqlx.DB) *Storage {
	return &Storage{
		db: db,
	}
}

var tracer = otel.Tracer("")

// GetGames returns list of games limited by pageSize and starting Id
func (s *Storage) GetGames(ctx context.Context, pageSize int, lastID int32) (list []Game, err error) {
	ctx, span := tracer.Start(ctx, "db.game.get")
	defer span.End()

	const q = `
	SELECT id, name, developer, publisher, release_date, genre, logo_url, rating
	FROM games
	WHERE id > $1
	ORDER BY id
	FETCH FIRST $2 ROWS ONLY`

	if err = s.db.SelectContext(ctx, &list, q, lastID, pageSize); err != nil {
		return nil, err
	}

	return list, nil
}

// SearchGames returns list of games by search query
func (s *Storage) SearchGames(ctx context.Context, search string) (list []Game, err error) {
	ctx, span := tracer.Start(ctx, "db.game.search")
	defer span.End()

	const q = `
	SELECT id, name, developer, publisher, release_date, genre, logo_url, rating
	FROM games
	WHERE LOWER(name) LIKE $1`

	if err = s.db.SelectContext(ctx, &list, q, strings.ToLower(search)+"%"); err != nil {
		return nil, err
	}

	return list, nil
}

// GetGameByID returns game by id.
// If game does not exist returns ErrNotFound
func (s *Storage) GetGameByID(ctx context.Context, id int32) (g Game, err error) {
	ctx, span := tracer.Start(ctx, "db.game.getbyid")
	defer span.End()

	const q = `SELECT id, name, developer, publisher, release_date, genre, logo_url, rating
	FROM games
	WHERE id = $1`

	if err = s.db.GetContext(ctx, &g, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Game{}, ErrNotFound[int32]{Entity: "game", ID: id}
		}
		return Game{}, err
	}

	return g, nil
}

// CreateGame creates new game
func (s *Storage) CreateGame(ctx context.Context, cg CreateGame) (id int32, err error) {
	ctx, span := tracer.Start(ctx, "db.game.create")
	defer span.End()

	const q = `INSERT INTO games
	(name, developer, publisher, release_date, genre, logo_url)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id`

	if err = s.db.QueryRowContext(ctx, q, cg.Name, cg.Developer, cg.Publisher, cg.ReleaseDate, pq.StringArray(cg.Genre), cg.LogoURL).Scan(&id); err != nil {
		return 0, fmt.Errorf("inserting game %s: %w", cg.Name, err)
	}

	return id, nil
}

// UpdateGame updates game
// If game does not exist returns ErrNotFound
func (s *Storage) UpdateGame(ctx context.Context, id int32, ug UpdateGame) error {
	ctx, span := tracer.Start(ctx, "db.game.update")
	defer span.End()

	const q = `UPDATE games 
	SET name = $2, developer = $3, publisher = $4, release_date = $5, genre = $6, logo_url = $7
	WHERE id = $1`

	releaseDate, err := types.ParseDate(ug.ReleaseDate)
	if err != nil {
		return fmt.Errorf("invalid date %s: %v", releaseDate.String(), err)
	}
	res, err := s.db.ExecContext(ctx, q, id, ug.Name, ug.Developer, ug.Publisher, releaseDate.String(), pq.StringArray(ug.Genre), ug.LogoURL)
	if err != nil {
		return fmt.Errorf("updating game %d: %v", id, err)
	}

	return checkRowsAffected(res, "game", id)
}

// UpdateGameRating updates game rating
// If game does not exist returns ErrNotFound
func (s *Storage) UpdateGameRating(ctx context.Context, id int32) error {
	ctx, span := tracer.Start(ctx, "db.game.updaterating")
	defer span.End()

	const q = `UPDATE games 
	SET rating = (
		SELECT SUM(rating)::numeric / COUNT(rating) 
		FROM ratings 
		WHERE game_id = $1)
	WHERE id = $1`

	res, err := s.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("updating game %d rating: %v", id, err)
	}

	return checkRowsAffected(res, "game", id)
}

// DeleteGame deletes game by id.
// If game does not exist returns ErrNotFound
func (s *Storage) DeleteGame(ctx context.Context, id int32) error {
	ctx, span := tracer.Start(ctx, "db.game.delete")
	defer span.End()

	const q = `DELETE FROM games 
	WHERE id = $1`
	res, err := s.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("deleting game %d: %v", id, err)
	}
	return checkRowsAffected(res, "game", id)
}
