package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/OutOfStack/game-library/pkg/types"
	"github.com/georgysavva/scany/v2/pgxscan"
)

const (
	gameReleaseYearCoeff = 2.5
	gameRatingCoeff      = 2.0
)

// OrderGamesBy options
var (
	OrderGamesByDefault = model.OrderBy{
		Field: "weight",
		Order: model.DescendingSortOrder,
	}
	OrderGamesByReleaseDate = model.OrderBy{
		Field: "release_date",
		Order: model.DescendingSortOrder,
	}
	OrderGamesByName = model.OrderBy{
		Field: "name",
		Order: model.AscendingSortOrder,
	}
)

// GetGames returns list of games with specified pageSize at specified page
func (s *Storage) GetGames(ctx context.Context, pageSize, page int, filter model.GamesFilter) (list []model.Game, err error) {
	ctx, span := tracer.Start(ctx, "db.getGames")
	defer span.End()

	query := psql.Select("id", "name", "release_date", "logo_url", "rating", "summary", "genres", "platforms",
		"screenshots", "developers", "publishers", "websites", "slug", "igdb_rating", "igdb_id",
		fmt.Sprintf("(extract(year FROM release_date)/%f + igdb_rating + rating/%f) weight", gameReleaseYearCoeff, gameRatingCoeff)).
		From("games").
		Limit(uint64(pageSize)).
		Offset(uint64((page - 1) * pageSize))

	if filter.OrderBy.Field != "" {
		query = query.OrderBy(fmt.Sprintf("%s %s", filter.OrderBy.Field, filter.OrderBy.Order))
	}

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

	if err = pgxscan.Select(ctx, s.db, &list, q, args...); err != nil {
		return nil, err
	}

	return list, nil
}

// GetGamesCount returns games count
func (s *Storage) GetGamesCount(ctx context.Context, filter model.GamesFilter) (count uint64, err error) {
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

	if err = pgxscan.Get(ctx, s.db, &count, q, args...); err != nil {
		return 0, err
	}

	return count, nil
}

// GetGameByID returns game by id.
// If game does not exist returns apperr.Error with NotFound status code
func (s *Storage) GetGameByID(ctx context.Context, id int32) (game model.Game, err error) {
	ctx, span := tracer.Start(ctx, "db.getGameByID")
	defer span.End()

	const q = `
		SELECT id, name, developers, publishers, release_date, genres, logo_url, rating, summary, platforms,
       		screenshots, websites, slug, igdb_rating, igdb_id
		FROM games
		WHERE id = $1`

	if err = pgxscan.Get(ctx, s.db, &game, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Game{}, apperr.NewNotFoundError("game", id)
		}
		return model.Game{}, err
	}

	return game, nil
}

// GetGameIDByIGDBID returns game id by igdb id.
// If game does not exist returns apperr.Error with NotFound status code
func (s *Storage) GetGameIDByIGDBID(ctx context.Context, igdbID int64) (id int32, err error) {
	ctx, span := tracer.Start(ctx, "db.getGameIdByIgdbId")
	defer span.End()

	const q = `
		SELECT id
		FROM games
		WHERE igdb_id = $1`

	if err = pgxscan.Get(ctx, s.db, &id, q, igdbID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, apperr.NewNotFoundError("game", igdbID)
		}
		return 0, err
	}

	return id, nil
}

// CreateGame creates new game
func (s *Storage) CreateGame(ctx context.Context, cg model.CreateGame) (id int32, err error) {
	ctx, span := tracer.Start(ctx, "db.createGame")
	defer span.End()

	const q = `
		INSERT INTO games
    		(name, developers, publishers, release_date, genres, logo_url, summary,
    		 platforms, screenshots, websites, slug, igdb_rating, igdb_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7,
		        $8, $9, $10, $11::varchar(50), $12, $13, $14)
		RETURNING id`

	err = s.db.QueryRow(ctx, q, cg.Name, cg.DevelopersIDs, cg.PublishersIDs, cg.ReleaseDate, cg.GenresIDs, cg.LogoURL, cg.Summary,
		cg.PlatformsIDs, cg.Screenshots, cg.Websites, cg.Slug, cg.IGDBRating, cg.IGDBID, time.Now()).
		Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("inserting game %s: %w", cg.Name, err)
	}

	return id, nil
}

// UpdateGame updates game
// If game does not exist returns apperr.Error with NotFound status code
func (s *Storage) UpdateGame(ctx context.Context, id int32, ug model.UpdateGameData) error {
	ctx, span := tracer.Start(ctx, "db.updateGame")
	defer span.End()

	const q = `
		UPDATE games
		SET name = $2, developers = $3, publishers = $4, release_date = $5, genres = $6, logo_url = $7, summary = $8,
		    platforms = $9, screenshots = $10, websites = $11, slug = $12, igdb_rating = $13, igdb_id = $14, updated_at = $15
		WHERE id = $1`

	releaseDate, err := types.ParseDate(ug.ReleaseDate)
	if err != nil {
		return fmt.Errorf("invalid date %v: %v", releaseDate, err)
	}
	res, err := s.db.Exec(ctx, q, id, ug.Name, ug.Developers, ug.Publishers, releaseDate.String(), ug.Genres, ug.LogoURL, ug.Summary,
		ug.PlatformsIDs, ug.Screenshots, ug.Websites, ug.Slug, ug.IGDBRating, ug.IGDBID, time.Now())
	if err != nil {
		return fmt.Errorf("updating game %d: %v", id, err)
	}

	return checkRowsAffected(res, "game", id)
}

// UpdateGameRating updates game rating
// If game does not exist returns apperr.Error with NotFound status code
func (s *Storage) UpdateGameRating(ctx context.Context, id int32) error {
	ctx, span := tracer.Start(ctx, "db.updateGameRating")
	defer span.End()

	const q = `
		UPDATE games
		SET rating = (
			SELECT COALESCE(SUM(rating)::numeric / COUNT(rating), 0)
			FROM ratings
			WHERE game_id = $1),
			updated_at = $2
		WHERE id = $1`

	res, err := s.db.Exec(ctx, q, id, time.Now())
	if err != nil {
		return fmt.Errorf("updating game %d rating: %v", id, err)
	}

	return checkRowsAffected(res, "game", id)
}

// DeleteGame deletes game by id.
// If game does not exist returns apperr.Error with NotFound status code
func (s *Storage) DeleteGame(ctx context.Context, id int32) error {
	ctx, span := tracer.Start(ctx, "db.deleteGame")
	defer span.End()

	const q = `
		DELETE FROM games
		WHERE id = $1`
	res, err := s.db.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("deleting game %d: %v", id, err)
	}
	return checkRowsAffected(res, "game", id)
}
