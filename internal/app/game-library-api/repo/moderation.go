package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/georgysavva/scany/v2/pgxscan"
)

// CreateModerationRecord creates a moderation record for a game
func (s *Storage) CreateModerationRecord(ctx context.Context, m model.CreateModeration) (id int32, err error) {
	ctx, span := tracer.Start(ctx, "createModerationRecord")
	defer span.End()

	const q = `
        INSERT INTO game_moderation (game_id, game_data, status, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id`

	if err = s.querier(ctx).QueryRow(ctx, q, m.GameID, m.GameData, m.Status, time.Now()).Scan(&id); err != nil {
		return 0, fmt.Errorf("creating moderation for game %d: %w", m.GameID, err)
	}

	return id, nil
}

// SetModerationRecordResultByGameID sets moderation result for moderation record by game id
func (s *Storage) SetModerationRecordResultByGameID(ctx context.Context, gameID int32, result model.UpdateModerationResult) error {
	ctx, span := tracer.Start(ctx, "setModerationRecordResultByGameID")
	defer span.End()

	const q = `
        UPDATE game_moderation
        SET status = $2, 
            details = $3, 
            attempts = CASE WHEN $2 != $4 THEN attempts + 1 ELSE attempts END, 
            updated_at = $5
        WHERE id = (
        	SELECT moderation_id 
        	FROM games 
        	WHERE id = $1
        	LIMIT 1
        )`

	res, err := s.querier(ctx).Exec(ctx, q, gameID, result.ResultStatus, result.Details, model.ModerationStatusReady, time.Now())
	if err != nil {
		return fmt.Errorf("setting moderation result for game id %d: %w", gameID, err)
	}

	return checkRowsAffected(res, "moderation_by_game_id", gameID)
}

// GetModerationRecordByID returns moderation record by id
func (s *Storage) GetModerationRecordByID(ctx context.Context, id int32) (m model.Moderation, err error) {
	ctx, span := tracer.Start(ctx, "getModerationRecordByID")
	defer span.End()

	const q = `
        SELECT id, game_id, status, details, attempts, game_data, created_at, updated_at
        FROM game_moderation
        WHERE id = $1`

	if err = pgxscan.Get(ctx, s.querier(ctx), &m, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Moderation{}, apperr.NewNotFoundError("moderation", id)
		}
		return model.Moderation{}, err
	}
	return m, nil
}

// GetModerationRecordByGameID returns current moderation record by game id
func (s *Storage) GetModerationRecordByGameID(ctx context.Context, gameID int32) (m model.Moderation, err error) {
	ctx, span := tracer.Start(ctx, "getModerationRecordByGameID")
	defer span.End()

	const q = `
        SELECT id, game_id, status, details, attempts, game_data, created_at, updated_at
        FROM game_moderation
        WHERE id = (
        	SELECT moderation_id 
        	FROM games 
        	WHERE id = $1
        	LIMIT 1
        )`

	if err = pgxscan.Get(ctx, s.querier(ctx), &m, q, gameID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Moderation{}, apperr.NewNotFoundError("moderation_by_game_id", gameID)
		}
		return model.Moderation{}, err
	}
	return m, nil
}

// GetModerationRecordsByGameID returns moderation records for a game ordered by newest first
func (s *Storage) GetModerationRecordsByGameID(ctx context.Context, gameID int32) (list []model.Moderation, err error) {
	ctx, span := tracer.Start(ctx, "getModerationRecordsByGameID")
	defer span.End()

	const q = `
        SELECT id, game_id, status, details, attempts, game_data, created_at, updated_at
        FROM game_moderation
        WHERE game_id = $1
        ORDER BY id DESC`

	if err = pgxscan.Select(ctx, s.querier(ctx), &list, q, gameID); err != nil {
		return nil, err
	}
	return list, nil
}

// GetPendingModerationGameIDs returns game IDs that have pending moderation status
func (s *Storage) GetPendingModerationGameIDs(ctx context.Context, limit int) ([]model.ModerationIDGameID, error) {
	ctx, span := tracer.Start(ctx, "getPendingModerationGameIDs")
	defer span.End()

	const q = `
        SELECT id, game_id
        FROM game_moderation
        WHERE status = $1
        ORDER BY id
        LIMIT $2
        FOR NO KEY UPDATE`

	var data []model.ModerationIDGameID
	if err := pgxscan.Select(ctx, s.querier(ctx), &data, q, model.ModerationStatusPending, limit); err != nil {
		return nil, fmt.Errorf("getting pending moderation game IDs: %w", err)
	}
	return data, nil
}

// SetModerationRecordStatus sets moderation status for games IDs
func (s *Storage) SetModerationRecordStatus(ctx context.Context, gameIDs []int32, status model.ModerationStatus) error {
	ctx, span := tracer.Start(ctx, "setModerationStatus")
	defer span.End()

	if len(gameIDs) == 0 {
		return nil
	}

	const q = `
        UPDATE game_moderation
        SET status = $2,
            attempts = CASE WHEN $2 = $3 THEN attempts + 1 ELSE attempts END,
            updated_at = $4
        WHERE id = ANY($1)`

	_, err := s.querier(ctx).Exec(ctx, q, gameIDs, status, model.ModerationStatusPending, time.Now())
	if err != nil {
		return fmt.Errorf("set moderation status to %s: %w", status, err)
	}
	return nil
}
