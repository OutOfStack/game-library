package repo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestCommit_Success tests committing tx
func TestCommit_Success(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	tx, err := s.BeginTx(context.Background())
	require.NoError(t, err)

	err = tx.Commit()
	require.NoError(t, err)
}

// TestRollback_Success tests rollback of tx
func TestRollback_Success(t *testing.T) {
	s := setup(t)
	defer teardown(t)

	tx, err := s.BeginTx(context.Background())
	require.NoError(t, err)

	err = tx.Rollback()
	require.NoError(t, err)
}
