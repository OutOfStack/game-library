package taskprocessor

import (
	"testing"

	"github.com/OutOfStack/game-library/internal/client/igdbapi"
	"github.com/stretchr/testify/require"
)

func TestMapTrailerURLs_FiltersDuplicatesAndCapsAtTwo(t *testing.T) {
	t.Parallel()

	videos := []igdbapi.Video{
		{VideoID: ""},
		{VideoID: " first "},
		{VideoID: "first"},
		{VideoID: "second"},
		{VideoID: "third"},
	}

	require.Equal(t, []string{
		"https://www.youtube.com/watch?v=first",
		"https://www.youtube.com/watch?v=second",
	}, mapTrailerURLs(videos))
}

func TestResolveRelatedGameIDs_PreservesAvailableUniqueOrderAndCapsAtSix(t *testing.T) {
	t.Parallel()

	gameIDs := map[int64]int32{
		10: 100,
		1:  101,
		2:  102,
		3:  103,
		4:  104,
		5:  105,
		6:  106,
		7:  107,
	}

	require.Equal(t, []int32{101, 102, 103, 104, 105, 106}, resolveRelatedGameIDs(
		10,
		[]int64{10, 1, 99, 2, 1, 3, 4, 5, 6, 7},
		gameIDs,
	))
}

func TestFranchiseGameIDs_FlattensMainAndAdditionalFranchises(t *testing.T) {
	t.Parallel()

	require.Equal(t, []int64{1, 2, 3, 4}, franchiseGameIDs(
		igdbapi.Franchise{Games: []int64{1, 2}},
		[]igdbapi.Franchise{{Games: []int64{3}}, {Games: []int64{4}}},
	))
}
