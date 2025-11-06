package repo_test

import (
	"context"
	"errors"
	"testing"

	"github.com/OutOfStack/game-library/internal/repo"
	mocks "github.com/OutOfStack/game-library/internal/repo/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestWithTx_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTx := mocks.NewMockTx(ctrl)
	ctx := t.Context()

	ctxWithTx := repo.WithTx(ctx, mockTx)

	retrievedTx, ok := repo.TxFromContext(ctxWithTx)
	require.True(t, ok)
	require.Equal(t, mockTx, retrievedTx)
}

func TestTxFromContext_NoTransaction(t *testing.T) {
	ctx := t.Context()

	tx, ok := repo.TxFromContext(ctx)
	require.False(t, ok)
	require.Nil(t, tx)
}

func TestStorage_RunWithTx_Success(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()
	executed := false

	err := s.RunWithTx(ctx, func(ctx context.Context) error {
		tx, ok := repo.TxFromContext(ctx)
		require.True(t, ok)
		require.NotNil(t, tx)
		executed = true
		return nil
	})

	require.NoError(t, err)
	require.True(t, executed)
}

func TestStorage_RunWithTx_RollbackOnError(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()
	expectedErr := errors.New("test error")

	err := s.RunWithTx(ctx, func(ctx context.Context) error {
		tx, ok := repo.TxFromContext(ctx)
		require.True(t, ok)
		require.NotNil(t, tx)
		return expectedErr
	})

	require.Error(t, err)
	require.Equal(t, expectedErr, err)
}

// TestCommit_Success tests committing tx
func TestCommit_Success(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()
	tx, err := s.BeginTx(ctx)
	require.NoError(t, err)

	err = tx.Commit(ctx)
	require.NoError(t, err)
}

// TestRollback_Success tests rollback of tx
func TestRollback_Success(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	ctx := t.Context()
	tx, err := s.BeginTx(ctx)
	require.NoError(t, err)

	err = tx.Rollback(ctx)
	require.NoError(t, err)
}
