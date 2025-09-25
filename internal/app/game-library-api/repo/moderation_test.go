package repo_test

import (
	"database/sql"
	"testing"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/pkg/apperr"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/stretchr/testify/require"
)

func TestModeration_CreateAndGetByID_ShouldMatch(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	// Create a game to reference
	cg := getCreateGameData()
	gameID, err := s.CreateGame(ctx, cg)
	require.NoError(t, err)

	// Create moderation record
	md := model.ModerationData{Name: cg.Name, Slug: cg.Slug}
	mid, err := s.CreateModerationRecord(ctx, model.CreateModeration{GameID: gameID, GameData: md})
	require.NoError(t, err)
	require.NotZero(t, mid)

	// Get by ID
	got, err := s.GetModerationRecordByID(ctx, mid)
	require.NoError(t, err)
	require.Equal(t, mid, got.ID)
	require.Equal(t, gameID, got.GameID)
	require.Empty(t, got.ResultStatus)
	require.Empty(t, got.Details)
	require.Equal(t, cg.Name, got.GameData.Name)
	require.Equal(t, cg.Slug, got.GameData.Slug)
}

func TestModeration_GetByGameID_ShouldBeOrderedDesc(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	// Create a game to reference
	cg := getCreateGameData()
	gameID, err := s.CreateGame(ctx, cg)
	require.NoError(t, err)

	// Create two moderation records
	id1, err := s.CreateModerationRecord(ctx, model.CreateModeration{GameID: gameID, GameData: model.ModerationData{Name: td.String()}})
	require.NoError(t, err)
	id2, err := s.CreateModerationRecord(ctx, model.CreateModeration{GameID: gameID, GameData: model.ModerationData{Name: td.String()}})
	require.NoError(t, err)

	list, err := s.GetModerationRecordsByGameID(ctx, gameID)
	require.NoError(t, err)
	require.Len(t, list, 2)
	require.Equal(t, id2, list[0].ID)
	require.Equal(t, id1, list[1].ID)
}

func TestModeration_SetResult_ShouldUpdateRecord(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	// Create game and moderation
	cg := getCreateGameData()
	gameID, err := s.CreateGame(ctx, cg)
	require.NoError(t, err)

	id, err := s.CreateModerationRecord(ctx, model.CreateModeration{GameID: gameID, GameData: model.ModerationData{Name: td.String()}})
	require.NoError(t, err)

	// Update result to declined with details and error
	e := td.String()
	err = s.SetModerationRecordResult(ctx, id, model.UpdateModerationResult{
		ResultStatus: model.ModerationStatusDeclined,
		Details:      "invalid logo url",
		Error:        sql.NullString{String: e, Valid: true},
	})
	require.NoError(t, err)

	got, err := s.GetModerationRecordByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, string(model.ModerationStatusDeclined), got.ResultStatus)
	require.Equal(t, "invalid logo url", got.Details)
	require.True(t, got.Error.Valid)
	require.Equal(t, e, got.Error.String)

	// Update result to ready and clear error
	err = s.SetModerationRecordResult(ctx, id, model.UpdateModerationResult{
		ResultStatus: model.ModerationStatusReady,
		Details:      "",
		Error:        sql.NullString{},
	})
	require.NoError(t, err)

	got, err = s.GetModerationRecordByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, string(model.ModerationStatusReady), got.ResultStatus)
	require.Empty(t, got.Details)
}

func TestModeration_GetRecordByID_NotFound_ShouldReturnError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	id := int32(td.Uint32())
	_, err := s.GetModerationRecordByID(t.Context(), id)
	require.ErrorIs(t, err, apperr.NewNotFoundError("moderation", id))
}
