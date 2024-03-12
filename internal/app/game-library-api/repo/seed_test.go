package repo_test

import (
	"testing"

	"github.com/OutOfStack/game-library/internal/app/game-library-manage/schema"
	"github.com/stretchr/testify/require"
)

// TestSeed_ShouldBeNoError tests case when seed data from seed.sql, and there should be no error
func TestSeed_ShouldBeNoError(t *testing.T) {
	setup(t)
	defer teardown(t)

	err := schema.Seed(db)
	require.NoError(t, err)
}
