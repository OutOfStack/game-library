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

// OrderGamesBy options
var (
	OrderGamesByDefault = model.OrderBy{
		Field: "trending_index",
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
	OrderGamesByRating = model.OrderBy{
		Field: "rating",
		Order: model.DescendingSortOrder,
	}
)

// GetGames returns list of games with specified pageSize at specified page
func (s *Storage) GetGames(ctx context.Context, pageSize, page uint32, filter model.GamesFilter) (list []model.Game, err error) {
	ctx, span := tracer.Start(ctx, "getGames")
	defer span.End()

	query := psql.Select("id", "name", "release_date", "logo_url", "rating", "summary", "genres", "platforms",
		"screenshots", "developers", "publishers", "websites", "slug", "igdb_rating", "igdb_rating_count", "igdb_id", "trending_index").
		From("games").
		Where(sq.Eq{"moderation_status": model.ModerationStatusReady}).
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

	if err = pgxscan.Select(ctx, s.querier(ctx), &list, q, args...); err != nil {
		return nil, err
	}

	return list, nil
}

// GetGamesCount returns games count
func (s *Storage) GetGamesCount(ctx context.Context, filter model.GamesFilter) (count uint64, err error) {
	ctx, span := tracer.Start(ctx, "getGamesCount")
	defer span.End()

	query := psql.Select("COUNT(id)").
		From("games").
		Where(sq.Eq{"moderation_status": model.ModerationStatusReady})

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

	if err = pgxscan.Get(ctx, s.querier(ctx), &count, q, args...); err != nil {
		return 0, err
	}

	return count, nil
}

// GetGameByID returns game by id.
// If game does not exist returns apperr.Error with NotFound status code
func (s *Storage) GetGameByID(ctx context.Context, id int32) (game model.Game, err error) {
	ctx, span := tracer.Start(ctx, "getGameById")
	defer span.End()

	const q = `
		SELECT id, name, developers, publishers, release_date, genres, logo_url, rating, summary, platforms,
       		screenshots, websites, slug, igdb_rating, igdb_rating_count, igdb_id, moderation_status, moderation_id, trending_index
		FROM games
		WHERE id = $1
		FOR UPDATE`

	if err = pgxscan.Get(ctx, s.querier(ctx), &game, q, id); err != nil {
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
	ctx, span := tracer.Start(ctx, "getGameIdByIgdbId")
	defer span.End()

	const q = `
		SELECT id
		FROM games
		WHERE igdb_id = $1`

	if err = pgxscan.Get(ctx, s.querier(ctx), &id, q, igdbID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, apperr.NewNotFoundError("game", igdbID)
		}
		return 0, err
	}

	return id, nil
}

// CreateGame creates new game
func (s *Storage) CreateGame(ctx context.Context, cg model.CreateGameData) (id int32, err error) {
	ctx, span := tracer.Start(ctx, "createGame")
	defer span.End()

	const q = `
		INSERT INTO games
    		(name, developers, publishers, release_date, genres, logo_url, summary,
    		 platforms, screenshots, websites, slug, igdb_rating, igdb_rating_count, igdb_id, moderation_status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7,
		        $8, $9, $10, $11::varchar(50), $12, $13, $14, $15, $16)
		RETURNING id`

	err = s.querier(ctx).QueryRow(ctx, q, cg.Name, cg.DevelopersIDs, cg.PublishersIDs, cg.ReleaseDate, cg.GenresIDs, cg.LogoURL, cg.Summary,
		cg.PlatformsIDs, cg.Screenshots, cg.Websites, cg.Slug, cg.IGDBRating, cg.IGDBRatingCount, cg.IGDBID, cg.ModerationStatus, time.Now()).
		Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("inserting game %s: %w", cg.Name, err)
	}

	return id, nil
}

// UpdateGame updates game
// If game does not exist returns apperr.Error with NotFound status code
func (s *Storage) UpdateGame(ctx context.Context, id int32, ug model.UpdateGameData) error {
	ctx, span := tracer.Start(ctx, "updateGame")
	defer span.End()

	const q = `
		UPDATE games
		SET name = $2, developers = $3, publishers = $4, release_date = $5, genres = $6, logo_url = $7, summary = $8,
		    platforms = $9, screenshots = $10, websites = $11, slug = $12, moderation_status = $13, updated_at = $14
		WHERE id = $1`

	releaseDate, err := types.ParseDate(ug.ReleaseDate)
	if err != nil {
		return fmt.Errorf("invalid date %v: %v", releaseDate, err)
	}
	res, err := s.querier(ctx).Exec(ctx, q, id,
		ug.Name, ug.DevelopersIDs, ug.PublishersIDs, releaseDate.String(), ug.GenresIDs, ug.LogoURL, ug.Summary,
		ug.PlatformsIDs, ug.Screenshots, ug.Websites, ug.Slug, ug.ModerationStatus, time.Now())
	if err != nil {
		return fmt.Errorf("updating game %d: %v", id, err)
	}

	return checkRowsAffected(res, "game", id)
}

// UpdateGameRating updates game rating
// If game does not exist returns apperr.Error with NotFound status code
func (s *Storage) UpdateGameRating(ctx context.Context, id int32) error {
	ctx, span := tracer.Start(ctx, "updateGameRating")
	defer span.End()

	const q = `
		UPDATE games
		SET rating = (
			SELECT COALESCE(SUM(rating)::numeric / COUNT(rating), 0)
			FROM ratings
			WHERE game_id = $1),
			updated_at = $2
		WHERE id = $1`

	res, err := s.querier(ctx).Exec(ctx, q, id, time.Now())
	if err != nil {
		return fmt.Errorf("updating game %d rating: %v", id, err)
	}

	return checkRowsAffected(res, "game", id)
}

// UpdateGameIGDBInfo updates game igdb info
// If game does not exist returns apperr.Error with NotFound status code
func (s *Storage) UpdateGameIGDBInfo(ctx context.Context, id int32, ug model.UpdateGameIGDBData) error {
	ctx, span := tracer.Start(ctx, "updateGameIgdbInfo")
	defer span.End()

	const q = `
		UPDATE games
		SET name = $2, platforms = $3, websites = $4, igdb_rating = $5, igdb_rating_count = $6, updated_at = $7
		WHERE id = $1`

	res, err := s.querier(ctx).Exec(ctx, q, id, ug.Name, ug.PlatformsIDs, ug.Websites, ug.IGDBRating, ug.IGDBRatingCount, time.Now())
	if err != nil {
		return fmt.Errorf("updating game %d igdb info: %v", id, err)
	}

	return checkRowsAffected(res, "game", id)
}

// DeleteGame deletes game by id.
// If game does not exist returns apperr.Error with NotFound status code
func (s *Storage) DeleteGame(ctx context.Context, id int32) error {
	ctx, span := tracer.Start(ctx, "deleteGame")
	defer span.End()

	const q = `
		DELETE FROM games
		WHERE id = $1`
	res, err := s.querier(ctx).Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("deleting game %d: %v", id, err)
	}
	return checkRowsAffected(res, "game", id)
}

// GetPublisherGamesCount returns the number of games created by a publisher in the specified date range
func (s *Storage) GetPublisherGamesCount(ctx context.Context, publisherID int32, startDate, endDate time.Time) (count int, err error) {
	ctx, span := tracer.Start(ctx, "getPublisherGamesCount")
	defer span.End()

	const q = `
		SELECT COUNT(id)
		FROM games
		WHERE $1 = ANY(publishers)
		AND created_at >= $2
		AND created_at <= $3`

	if err = pgxscan.Get(ctx, s.querier(ctx), &count, q, publisherID, startDate, endDate); err != nil {
		return 0, err
	}

	return count, nil
}

// UpdateGameTrendingIndex updates the trending index for a specific game
func (s *Storage) UpdateGameTrendingIndex(ctx context.Context, gameID int32, trendingIndex float64) error {
	ctx, span := tracer.Start(ctx, "updateGameTrendingIndex")
	defer span.End()

	const q = `
		UPDATE games
		SET trending_index = $2, updated_at = $3
		WHERE id = $1`

	res, err := s.querier(ctx).Exec(ctx, q, gameID, trendingIndex, time.Now())
	if err != nil {
		return fmt.Errorf("updating trending index for game %d: %w", gameID, err)
	}

	return checkRowsAffected(res, "game", gameID)
}

// UpdateGameModerationID updates game moderation id
func (s *Storage) UpdateGameModerationID(ctx context.Context, gameID, moderationID int32) error {
	ctx, span := tracer.Start(ctx, "updateGameModerationId")
	defer span.End()

	const q = `
		UPDATE games
		SET moderation_id = $2
		WHERE id = $1`

	res, err := s.querier(ctx).Exec(ctx, q, gameID, moderationID)
	if err != nil {
		return fmt.Errorf("updating moderation id for game %d: %w", gameID, err)
	}

	return checkRowsAffected(res, "game", gameID)
}

// GetGameTrendingData retrieves data needed for trending index calculation
func (s *Storage) GetGameTrendingData(ctx context.Context, gameID int32) (model.GameTrendingData, error) {
	ctx, span := tracer.Start(ctx, "getGameTrendingData")
	defer span.End()

	const q = `
		SELECT 
			EXTRACT(year FROM release_date)::int as release_year,
			EXTRACT(month FROM release_date)::int as release_month,
			igdb_rating,
			igdb_rating_count,
			rating,
			COALESCE((SELECT COUNT(*) FROM ratings WHERE game_id = $1), 0) as rating_count
		FROM games
		WHERE id = $1`

	var data model.GameTrendingData
	err := pgxscan.Get(ctx, s.querier(ctx), &data, q, gameID)
	if err != nil {
		return model.GameTrendingData{}, fmt.Errorf("querying game trending data: %w", err)
	}

	return data, nil
}

// GetGamesIDsAfterID returns games ids after provided id
func (s *Storage) GetGamesIDsAfterID(ctx context.Context, lastID int32, batchSize int) ([]int32, error) {
	ctx, span := tracer.Start(ctx, "getGamesForTrendingIndexUpdate")
	defer span.End()

	const q = `
		SELECT id
		FROM games
		WHERE id > $1
		ORDER BY id
		LIMIT $2`

	var gamesIDs []int32
	err := pgxscan.Select(ctx, s.querier(ctx), &gamesIDs, q, lastID, batchSize)
	if err != nil {
		return nil, fmt.Errorf("querying games ids after id %d: %w", lastID, err)
	}

	return gamesIDs, nil
}

// GetGamesByPublisherID returns games created by a specific publisher company id
func (s *Storage) GetGamesByPublisherID(ctx context.Context, publisherID int32) (list []model.Game, err error) {
	ctx, span := tracer.Start(ctx, "getGamesByPublisherID")
	defer span.End()

	const q = `
        SELECT id, name, developers, publishers, release_date, genres, logo_url, rating, summary, platforms,
               screenshots, websites, slug, igdb_rating, igdb_rating_count, igdb_id, moderation_status
        FROM games
        WHERE $1 = ANY(publishers)
        ORDER BY id DESC`

	if err = pgxscan.Select(ctx, s.querier(ctx), &list, q, publisherID); err != nil {
		return nil, err
	}

	return list, nil
}
