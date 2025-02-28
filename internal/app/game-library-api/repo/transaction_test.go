package repo_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
