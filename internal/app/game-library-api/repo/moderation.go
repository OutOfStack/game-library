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

// SetModerationRecordResult sets moderation result for moderation record
func (s *Storage) SetModerationRecordResult(ctx context.Context, id int32, res model.UpdateModerationResult) error {
	ctx, span := tracer.Start(ctx, "setModerationRecordResult")
	defer span.End()

	const q = `
        UPDATE game_moderation
        SET status = $2, details = $3, error = $4, updated_at = $5
        WHERE id = $1`

	execRes, err := s.querier(ctx).Exec(ctx, q, id, res.ResultStatus, res.Details, res.Error, time.Now())
	if err != nil {
		return fmt.Errorf("setting moderation %d result: %w", id, err)
	}

	return checkRowsAffected(execRes, "moderation", id)
}

// GetModerationRecordByID returns moderation record by id
func (s *Storage) GetModerationRecordByID(ctx context.Context, id int32) (m model.Moderation, err error) {
	ctx, span := tracer.Start(ctx, "getModerationRecordByID")
	defer span.End()

	const q = `
        SELECT id, game_id, status, details, error, game_data, created_at, updated_at
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

// GetModerationRecordsByGameID returns moderation records for a game ordered by newest first
func (s *Storage) GetModerationRecordsByGameID(ctx context.Context, gameID int32) (list []model.Moderation, err error) {
	ctx, span := tracer.Start(ctx, "getModerationRecordsByGameID")
	defer span.End()

	const q = `
        SELECT id, game_id, status, details, error, game_data, created_at, updated_at
        FROM game_moderation
        WHERE game_id = $1
        ORDER BY id DESC`

	if err = pgxscan.Select(ctx, s.querier(ctx), &list, q, gameID); err != nil {
		return nil, err
	}
	return list, nil
}
