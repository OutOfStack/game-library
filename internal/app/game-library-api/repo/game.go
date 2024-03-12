package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/OutOfStack/game-library/pkg/types"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.opentelemetry.io/otel"
)

var (
	tracer = otel.Tracer("")
)

var (
	psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

// SortOrder - type for query sort order
type SortOrder string

// SortOrder valeus
const (
	AscendingSortOrder  SortOrder = "ASC"
	DescendingSortOrder SortOrder = "DESC"
)

const (
	gameReleaseYearCoeff = 2.5
	gameRatingCoeff      = 2.0
)

// GamesOrderBy type of games ordering
type GamesOrderBy struct {
	Field string
	Order SortOrder
}

// OrderGamesBy options
var (
	OrderGamesByDefault = GamesOrderBy{
		Field: "weight",
		Order: DescendingSortOrder,
	}
	OrderGamesByReleaseDate = GamesOrderBy{
		Field: "release_date",
		Order: DescendingSortOrder,
	}
	OrderGamesByName = GamesOrderBy{
		Field: "name",
		Order: AscendingSortOrder,
	}
)

// GamesFilter games filter
type GamesFilter struct {
	Name        string
	DeveloperID int32
	PublisherID int32
	GenreID     int32
	OrderBy     GamesOrderBy
}

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

// GetGames returns list of games with specified pageSize at specified page
func (s *Storage) GetGames(ctx context.Context, pageSize, page int, filter GamesFilter) (list []Game, err error) {
	ctx, span := tracer.Start(ctx, "db.getGames")
	defer span.End()

	query := psql.Select("id", "name", "release_date", "logo_url", "rating", "summary", "genres", "platforms",
		"screenshots", "developers", "publishers", "websites", "slug", "igdb_rating", "igdb_id",
		fmt.Sprintf("(extract(year FROM release_date)/%f + igdb_rating + rating/%f) weight", gameReleaseYearCoeff, gameRatingCoeff)).
		From("games").
		Limit(uint64(pageSize)).
		Offset(uint64((page - 1) * pageSize)).
		OrderBy(fmt.Sprintf("%s %s", filter.OrderBy.Field, filter.OrderBy.Order))

	if filter.Name != "" {
		query = query.Where(sq.Like{"LOWER(name)": "%" + strings.ToLower(filter.Name) + "%"})
	}
	if filter.GenreID != 0 {
		query = query.Where(sq.Expr("? = ANY(genres)", filter.GenreID))
	}
	if filter.PublisherID != 0 {
		query = query.Where(sq.Expr("? = ANY(publishers)", filter.PublisherID))
	}
	if filter.DeveloperID != 0 {
		query = query.Where(sq.Expr("? = ANY(developers)", filter.DeveloperID))
	}

	q, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	if err = s.db.SelectContext(ctx, &list, q, args...); err != nil {
		return nil, err
	}

	return list, nil
}

// GetGamesCount returns games count
func (s *Storage) GetGamesCount(ctx context.Context, filter GamesFilter) (count uint64, err error) {
	ctx, span := tracer.Start(ctx, "db.getGamesCount")
	defer span.End()

	query := psql.Select("COUNT(id)").
		From("games")

	if filter.Name != "" {
		query = query.Where(sq.Like{"LOWER(name)": "%" + strings.ToLower(filter.Name) + "%"})
	}
	if filter.GenreID != 0 {
		query = query.Where(sq.Expr("? = ANY(genres)", filter.GenreID))
	}
	if filter.PublisherID != 0 {
		query = query.Where(sq.Expr("? = ANY(publishers)", filter.PublisherID))
	}
	if filter.DeveloperID != 0 {
		query = query.Where(sq.Expr("? = ANY(developers)", filter.DeveloperID))
	}

	q, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}

	if err = s.db.GetContext(ctx, &count, q, args...); err != nil {
		return 0, err
	}

	return count, nil
}

// GetGameByID returns game by id.
// If game does not exist returns ErrNotFound
func (s *Storage) GetGameByID(ctx context.Context, id int32) (game Game, err error) {
	ctx, span := tracer.Start(ctx, "db.getGameByID")
	defer span.End()

	const q = `SELECT id, name, developers, publishers, release_date, genres, logo_url, rating, summary, platforms,
       screenshots, websites, slug, igdb_rating, igdb_id
	FROM games
	WHERE id = $1`

	if err = s.db.GetContext(ctx, &game, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Game{}, ErrNotFound[int32]{Entity: "game", ID: id}
		}
		return Game{}, err
	}

	return game, nil
}

// GetGameIDByIGDBID returns game id by igdb id.
// If game does not exist returns ErrNotFound
func (s *Storage) GetGameIDByIGDBID(ctx context.Context, igdbID int64) (id int32, err error) {
	ctx, span := tracer.Start(ctx, "db.getGameIdByIgdbId")
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
	ctx, span := tracer.Start(ctx, "db.createGame")
	defer span.End()

	const q = `INSERT INTO games
    (name, developers, publishers, release_date, genres, logo_url, summary, platforms, screenshots,
     	websites, slug, igdb_rating, igdb_id, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11::varchar(50), $12, $13, $14)
	RETURNING id`

	err = s.db.QueryRowContext(ctx, q, cg.Name, pq.Int32Array(cg.Developers), pq.Int32Array(cg.Publishers),
		cg.ReleaseDate, pq.Int32Array(cg.Genres), cg.LogoURL, cg.Summary, pq.Int32Array(cg.Platforms),
		pq.StringArray(cg.Screenshots), pq.StringArray(cg.Websites), cg.Slug, cg.IGDBRating, cg.IGDBID, time.Now()).
		Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("inserting game %s: %w", cg.Name, err)
	}

	return id, nil
}

// UpdateGame updates game
// If game does not exist returns ErrNotFound
func (s *Storage) UpdateGame(ctx context.Context, id int32, ug UpdateGame) error {
	ctx, span := tracer.Start(ctx, "db.updateGame")
	defer span.End()

	const q = `UPDATE games
	SET name = $2, developers = $3, publishers = $4, release_date = $5, genres = $6, logo_url = $7, summary = $8, platforms = $9,
	    screenshots = $10, websites = $11, slug = $12, igdb_rating = $13, igdb_id = $14, updated_at = $15
	WHERE id = $1`

	releaseDate, err := types.ParseDate(ug.ReleaseDate)
	if err != nil {
		return fmt.Errorf("invalid date %s: %v", releaseDate.String(), err)
	}
	res, err := s.db.ExecContext(ctx, q, id, ug.Name, pq.Int32Array(ug.Developers), pq.Int32Array(ug.Publishers), releaseDate.String(),
		pq.Int32Array(ug.Genres), ug.LogoURL, ug.Summary, pq.Int32Array(ug.Platforms), pq.StringArray(ug.Screenshots),
		pq.StringArray(ug.Websites), ug.Slug, ug.IGDBRating, ug.IGDBID, time.Now())
	if err != nil {
		return fmt.Errorf("updating game %d: %v", id, err)
	}

	return checkRowsAffected(res, "game", id)
}

// UpdateGameRating updates game rating
// If game does not exist returns ErrNotFound
func (s *Storage) UpdateGameRating(ctx context.Context, id int32) error {
	ctx, span := tracer.Start(ctx, "db.updateGameRating")
	defer span.End()

	const q = `UPDATE games
	SET rating = (
		SELECT COALESCE(SUM(rating)::numeric / COUNT(rating), 0)
		FROM ratings
		WHERE game_id = $1),
	    updated_at = $2
	WHERE id = $1`

	res, err := s.db.ExecContext(ctx, q, id, time.Now())
	if err != nil {
		return fmt.Errorf("updating game %d rating: %v", id, err)
	}

	return checkRowsAffected(res, "game", id)
}

// DeleteGame deletes game by id.
// If game does not exist returns ErrNotFound
func (s *Storage) DeleteGame(ctx context.Context, id int32) error {
	ctx, span := tracer.Start(ctx, "db.deleteGame")
	defer span.End()

	const q = `DELETE FROM games
	WHERE id = $1`
	res, err := s.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("deleting game %d: %v", id, err)
	}
	return checkRowsAffected(res, "game", id)
}
