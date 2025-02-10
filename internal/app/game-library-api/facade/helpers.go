package facade

import (
	"strconv"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
)

// Cache keys
const (
	gamesKey        = "games"
	gameKey         = "game"
	gamesCountKey   = "games-count"
	userRatingsKey  = "user-ratings"
	companiesKey    = "companies"
	topCompaniesKey = "top-companies"
	genresKey       = "genres"
	topGenresKey    = "top-genres"
	platformsKey    = "platforms"
)

func getGamesKey(pageSize, page int64, filter model.GamesFilter) string {
	return gamesKey + "|" + strconv.FormatInt(pageSize, 10) + "|" + strconv.FormatInt(page, 10) + "|" +
		filter.OrderBy.Field + "|" + filter.Name + "|" + strconv.FormatInt(int64(filter.GenreID), 10) + "|" +
		strconv.FormatInt(int64(filter.DeveloperID), 10) + "|" + strconv.FormatInt(int64(filter.PublisherID), 10)
}

func getGameKey(id int32) string {
	return gameKey + "|" + strconv.FormatInt(int64(id), 10)
}

func getGamesCountKey(filter model.GamesFilter) string {
	return gamesCountKey + "|" + filter.Name + "|" + strconv.FormatInt(int64(filter.GenreID), 10) + "|" +
		strconv.FormatInt(int64(filter.DeveloperID), 10) + "|" + strconv.FormatInt(int64(filter.PublisherID), 10)
}

func getUserRatingsKey(userID string) string {
	return userRatingsKey + "|" + userID
}

func getCompaniesKey() string {
	return companiesKey
}

func getTopCompaniesKey(companyType string, limit int64) string {
	return topCompaniesKey + "|" + companyType + "|" + strconv.FormatInt(limit, 10)
}

func getGenresKey() string {
	return genresKey
}

func getTopGenresKey(limit int64) string {
	return topGenresKey + "|" + strconv.FormatInt(limit, 10)
}

func getPlatformsKey() string {
	return platformsKey
}
