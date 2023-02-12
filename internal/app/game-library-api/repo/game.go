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

// GetGames returns list of games limited by pageSize and starting id
func (s *Storage) GetGames(ctx context.Context, pageSize int, lastID int32) (list []Game, err error) {
	ctx, span := tracer.Start(ctx, "db.game.get")
	defer span.End()

	const q = `
	SELECT id, name, developer, publisher, release_date, genre, logo_url, rating, summary, genres, platforms, screenshots, developers, 
	       publishers, websites, slug, igdb_rating, igdb_id
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
	SELECT id, name, developer, developers, publisher, publishers, release_date, genre, genres, logo_url, rating, summary, 
	       platforms, screenshots, websites, slug, igdb_rating, igdb_id
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

	const q = `SELECT id, name, developer, developers, publisher, publishers, release_date, genre, genres, logo_url, rating, 
       summary, platforms, screenshots, websites, slug, igdb_rating, igdb_id
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

// GetGameIDByIGDBID returns game id by igdb id.
// If game does not exist returns ErrNotFound
func (s *Storage) GetGameIDByIGDBID(ctx context.Context, igdbID int64) (id int32, err error) {
	ctx, span := tracer.Start(ctx, "db.game.getidbyigdbid")
	defer span.End()

	const q = `SELECT id
	FROM games
	WHERE igdb_id = $1`

	if err = s.db.GetContext(ctx, &id, q, igdbID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNotFound[int64]{Entity: "game", ID: igdbID}
		}
		return 0, err
	}

	return id, nil
}

// CreateGame creates new game
func (s *Storage) CreateGame(ctx context.Context, cg CreateGame) (id int32, err error) {
	ctx, span := tracer.Start(ctx, "db.game.create")
	defer span.End()

	const q = `INSERT INTO games
    (name, developer, publisher, developers, publishers, release_date, genre, genres, logo_url, summary, platforms, screenshots, 
     	websites, slug, igdb_rating, igdb_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14::varchar(50), $15, $16)
	RETURNING id`

	err = s.db.QueryRowContext(ctx, q, cg.Name, cg.Developer, cg.Publisher, pq.Int32Array(cg.Developers), pq.Int32Array(cg.Publishers),
		cg.ReleaseDate, pq.StringArray(cg.Genre), pq.Int32Array(cg.Genres), cg.LogoURL, cg.Summary, pq.Int32Array(cg.Platforms),
		pq.StringArray(cg.Screenshots), pq.StringArray(cg.Websites), cg.Slug, cg.IGDBRating, cg.IGDBID).Scan(&id)
	if err != nil {
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
	SET name = $2, developers = $3, publishers = $4, release_date = $5, genres = $6, logo_url = $7, summary = $8, platforms = $9,
	    screenshots = $10, websites = $11, slug = $12, igdb_rating = $13, igdb_id = $14
	WHERE id = $1`

	releaseDate, err := types.ParseDate(ug.ReleaseDate)
	if err != nil {
		return fmt.Errorf("invalid date %s: %v", releaseDate.String(), err)
	}
	res, err := s.db.ExecContext(ctx, q, id, ug.Name, pq.Int32Array(ug.Developers), pq.Int32Array(ug.Publishers), releaseDate.String(),
		pq.Int32Array(ug.Genres), ug.LogoURL, ug.Summary, pq.Int32Array(ug.Platforms), pq.StringArray(ug.Screenshots),
		pq.StringArray(ug.Websites), ug.Slug, ug.IGDBRating, ug.IGDBID)
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
