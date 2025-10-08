package repo_test

import (
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
	md := model.ModerationData{
		Name:    cg.Name,
		Summary: cg.Summary,
	}
	mid, err := s.CreateModerationRecord(ctx, model.CreateModeration{GameID: gameID, GameData: md})
	require.NoError(t, err)
	require.NotZero(t, mid)

	// Get by ID
	got, err := s.GetModerationRecordByID(ctx, mid)
	require.NoError(t, err)
	require.Equal(t, mid, got.ID)
	require.Equal(t, gameID, got.GameID)
	require.Empty(t, got.Status)
	require.Empty(t, got.Details)
	require.Equal(t, cg.Name, got.GameData.Name)
	require.Equal(t, cg.Summary, got.GameData.Summary)
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

func TestModeration_SetResultByGameID_ShouldUpdateRecord(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	// Create game and moderation
	cg := getCreateGameData()
	gameID, err := s.CreateGame(ctx, cg)
	require.NoError(t, err)

	mid, err := s.CreateModerationRecord(ctx, model.CreateModeration{
		GameID:   gameID,
		GameData: model.ModerationData{Name: td.String()},
		Status:   model.ModerationStatusPending,
	})
	require.NoError(t, err)

	// Link moderation to game
	_, err = db.Exec(ctx, "UPDATE games SET moderation_id = $1 WHERE id = $2", mid, gameID)
	require.NoError(t, err)

	// Update result to declined
	err = s.SetModerationRecordResultByGameID(ctx, gameID, model.UpdateModerationResult{
		ResultStatus: model.ModerationStatusDeclined,
		Details:      "invalid logo url",
	})
	require.NoError(t, err)

	got, err := s.GetModerationRecordByGameID(ctx, gameID)
	require.NoError(t, err)
	require.Equal(t, string(model.ModerationStatusDeclined), got.Status)
	require.Equal(t, "invalid logo url", got.Details)
	require.Equal(t, int32(1), got.Attempts)

	// Update result to ready
	err = s.SetModerationRecordResultByGameID(ctx, gameID, model.UpdateModerationResult{
		ResultStatus: model.ModerationStatusReady,
		Details:      "approved",
	})
	require.NoError(t, err)

	got, err = s.GetModerationRecordByGameID(ctx, gameID)
	require.NoError(t, err)
	require.Equal(t, string(model.ModerationStatusReady), got.Status)
	require.Equal(t, "approved", got.Details)
	require.Equal(t, int32(1), got.Attempts)
}

func TestModeration_GetRecordByID_NotFound_ShouldReturnError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	id := int32(td.Uint32())
	_, err := s.GetModerationRecordByID(t.Context(), id)
	require.ErrorIs(t, err, apperr.NewNotFoundError("moderation", id))
}

func TestModeration_GetRecordByGameID_ShouldReturnRecord(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	cg := getCreateGameData()
	gameID, err := s.CreateGame(ctx, cg)
	require.NoError(t, err)

	md := model.ModerationData{Name: td.String(), Summary: td.String()}
	mid, err := s.CreateModerationRecord(ctx, model.CreateModeration{
		GameID:   gameID,
		GameData: md,
		Status:   model.ModerationStatusPending,
	})
	require.NoError(t, err)

	_, err = db.Exec(ctx, "UPDATE games SET moderation_id = $1 WHERE id = $2", mid, gameID)
	require.NoError(t, err)

	got, err := s.GetModerationRecordByGameID(ctx, gameID)
	require.NoError(t, err)
	require.Equal(t, mid, got.ID)
	require.Equal(t, gameID, got.GameID)
	require.Equal(t, string(model.ModerationStatusPending), got.Status)
	require.Equal(t, md.Name, got.GameData.Name)
	require.Equal(t, md.Summary, got.GameData.Summary)
	require.Equal(t, int32(0), got.Attempts)
}

func TestModeration_GetRecordByGameID_NotFound_ShouldReturnError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	gameID := int32(td.Uint32())
	_, err := s.GetModerationRecordByGameID(t.Context(), gameID)
	require.ErrorIs(t, err, apperr.NewNotFoundError("moderation_by_game_id", gameID))
}

func TestModeration_GetPendingGameIDs_ShouldReturnPendingOnly(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	cg1 := getCreateGameData()
	gameID1, err := s.CreateGame(ctx, cg1)
	require.NoError(t, err)

	cg2 := getCreateGameData()
	gameID2, err := s.CreateGame(ctx, cg2)
	require.NoError(t, err)

	cg3 := getCreateGameData()
	gameID3, err := s.CreateGame(ctx, cg3)
	require.NoError(t, err)

	mid1, err := s.CreateModerationRecord(ctx, model.CreateModeration{
		GameID:   gameID1,
		GameData: model.ModerationData{Name: td.String()},
		Status:   model.ModerationStatusPending,
	})
	require.NoError(t, err)

	mid2, err := s.CreateModerationRecord(ctx, model.CreateModeration{
		GameID:   gameID2,
		GameData: model.ModerationData{Name: td.String()},
		Status:   model.ModerationStatusPending,
	})
	require.NoError(t, err)

	_, err = s.CreateModerationRecord(ctx, model.CreateModeration{
		GameID:   gameID3,
		GameData: model.ModerationData{Name: td.String()},
		Status:   model.ModerationStatusReady,
	})
	require.NoError(t, err)

	records, err := s.GetPendingModerationGameIDs(ctx, 10)
	require.NoError(t, err)
	require.Len(t, records, 2)
	require.Equal(t, mid1, records[0].ModerationID)
	require.Equal(t, gameID1, records[0].GameID)
	require.Equal(t, mid2, records[1].ModerationID)
	require.Equal(t, gameID2, records[1].GameID)
}

func TestModeration_GetPendingGameIDs_ShouldRespectLimit(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	for range 5 {
		cg := getCreateGameData()
		gameID, err := s.CreateGame(ctx, cg)
		require.NoError(t, err)

		_, err = s.CreateModerationRecord(ctx, model.CreateModeration{
			GameID:   gameID,
			GameData: model.ModerationData{Name: td.String()},
			Status:   model.ModerationStatusPending,
		})
		require.NoError(t, err)
	}

	records, err := s.GetPendingModerationGameIDs(ctx, 3)
	require.NoError(t, err)
	require.Len(t, records, 3)
}

func TestModeration_SetModerationRecordStatus_ShouldUpdateStatus(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	cg1 := getCreateGameData()
	gameID1, err := s.CreateGame(ctx, cg1)
	require.NoError(t, err)

	cg2 := getCreateGameData()
	gameID2, err := s.CreateGame(ctx, cg2)
	require.NoError(t, err)

	mid1, err := s.CreateModerationRecord(ctx, model.CreateModeration{
		GameID:   gameID1,
		GameData: model.ModerationData{Name: td.String()},
		Status:   model.ModerationStatusPending,
	})
	require.NoError(t, err)

	mid2, err := s.CreateModerationRecord(ctx, model.CreateModeration{
		GameID:   gameID2,
		GameData: model.ModerationData{Name: td.String()},
		Status:   model.ModerationStatusPending,
	})
	require.NoError(t, err)

	err = s.SetModerationRecordStatus(ctx, []int32{mid1, mid2}, model.ModerationStatusInProgress)
	require.NoError(t, err)

	m1, err := s.GetModerationRecordByID(ctx, mid1)
	require.NoError(t, err)
	require.Equal(t, string(model.ModerationStatusInProgress), m1.Status)

	m2, err := s.GetModerationRecordByID(ctx, mid2)
	require.NoError(t, err)
	require.Equal(t, string(model.ModerationStatusInProgress), m2.Status)
}

func TestModeration_SetModerationRecordStatus_ShouldIncrementAttemptsWhenSettingToPending(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()

	cg := getCreateGameData()
	gameID, err := s.CreateGame(ctx, cg)
	require.NoError(t, err)

	mid, err := s.CreateModerationRecord(ctx, model.CreateModeration{
		GameID:   gameID,
		GameData: model.ModerationData{Name: td.String()},
		Status:   model.ModerationStatusInProgress,
	})
	require.NoError(t, err)

	err = s.SetModerationRecordStatus(ctx, []int32{mid}, model.ModerationStatusPending)
	require.NoError(t, err)

	got, err := s.GetModerationRecordByID(ctx, mid)
	require.NoError(t, err)
	require.Equal(t, string(model.ModerationStatusPending), got.Status)
	require.Equal(t, int32(1), got.Attempts)
}

func TestModeration_SetModerationRecordStatus_WithEmptySlice_ShouldNotError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	err := s.SetModerationRecordStatus(t.Context(), []int32{}, model.ModerationStatusInProgress)
	require.NoError(t, err)
}
