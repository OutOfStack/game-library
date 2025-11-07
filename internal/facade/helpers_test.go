package facade

import (
	"testing"

	"github.com/OutOfStack/game-library/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestGetGamesKey(t *testing.T) {
	filter := model.GamesFilter{
		OrderBy:     model.OrderBy{Field: "name"},
		Name:        "test",
		GenreID:     1,
		DeveloperID: 2,
		PublisherID: 3,
	}

	expectedKey := "games|10|1|name|test|1|2|3"
	key := getGamesKey(10, 1, filter)

	assert.Equal(t, expectedKey, key)
}

func TestGetGameKey(t *testing.T) {
	expectedKey := "game|123"
	key := getGameKey(123)

	assert.Equal(t, expectedKey, key)
}

func TestGetGamesCountKey(t *testing.T) {
	filter := model.GamesFilter{
		Name:        "test",
		GenreID:     1,
		DeveloperID: 2,
		PublisherID: 3,
	}

	expectedKey := "games-count|test|1|2|3"
	key := getGamesCountKey(filter)

	assert.Equal(t, expectedKey, key)
}

func TestGetUserRatingsKey(t *testing.T) {
	expectedKey := "user-ratings|user123"
	key := getUserRatingsKey("user123")

	assert.Equal(t, expectedKey, key)
}

func TestGetCompaniesKey(t *testing.T) {
	expectedKey := "companies"
	key := getCompaniesKey()

	assert.Equal(t, expectedKey, key)
}

func TestGetTopCompaniesKey(t *testing.T) {
	expectedKey := "top-companies|developer|10"
	key := getTopCompaniesKey("developer", 10)

	assert.Equal(t, expectedKey, key)
}

func TestGetGenresKey(t *testing.T) {
	expectedKey := "genres"
	key := getGenresKey()

	assert.Equal(t, expectedKey, key)
}

func TestGetTopGenresKey(t *testing.T) {
	expectedKey := "top-genres|5"
	key := getTopGenresKey(5)

	assert.Equal(t, expectedKey, key)
}

func TestGetPlatformsKey(t *testing.T) {
	expectedKey := "platforms"
	key := getPlatformsKey()

	assert.Equal(t, expectedKey, key)
}
