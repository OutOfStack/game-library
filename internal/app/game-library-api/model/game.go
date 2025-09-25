package model

import (
	"database/sql"
	"strings"

	"github.com/OutOfStack/game-library/pkg/types"
)

// Game - db game model
type Game struct {
	ID               int32            `db:"id"`
	Name             string           `db:"name"`
	DevelopersIDs    []int32          `db:"developers"`
	PublishersIDs    []int32          `db:"publishers"`
	ReleaseDate      types.Date       `db:"release_date"`
	GenresIDs        []int32          `db:"genres"`
	LogoURL          string           `db:"logo_url"`
	Rating           float64          `db:"rating"`
	Summary          string           `db:"summary"`
	Slug             string           `db:"slug"`
	PlatformsIDs     []int32          `db:"platforms"`
	Screenshots      []string         `db:"screenshots"`
	Websites         []string         `db:"websites"`
	IGDBRating       float64          `db:"igdb_rating"`
	IGDBRatingCount  int32            `db:"igdb_rating_count"`
	IGDBID           int64            `db:"igdb_id"`
	ModerationStatus ModerationStatus `db:"moderation_status"`
	ModerationID     sql.NullInt32    `db:"moderation_id"`
	TrendingIndex    float64          `db:"trending_index"`
}

// CreateGameData - data for creating game in db
type CreateGameData struct {
	Name             string
	DevelopersIDs    []int32
	PublishersIDs    []int32
	ReleaseDate      string
	GenresIDs        []int32
	LogoURL          string
	Summary          string
	Slug             string
	PlatformsIDs     []int32
	Screenshots      []string
	Websites         []string
	IGDBRating       float64
	IGDBRatingCount  int32
	IGDBID           int64
	ModerationStatus ModerationStatus
}

// CreateGame - create game data
type CreateGame struct {
	Name         string
	ReleaseDate  string
	GenresIDs    []int32
	LogoURL      string
	Summary      string
	Slug         string
	PlatformsIDs []int32
	Screenshots  []string
	Websites     []string
	Developer    string // helper field
	Publisher    string // helper field
}

// UpdateGameData - data for updating game in db
type UpdateGameData struct {
	Name             string
	DevelopersIDs    []int32
	PublishersIDs    []int32
	ReleaseDate      string
	GenresIDs        []int32
	LogoURL          string
	Summary          string
	Slug             string
	PlatformsIDs     []int32
	Screenshots      []string
	Websites         []string
	ModerationStatus ModerationStatus
}

// UpdateGameIGDBData - igdb data for updating game in db
type UpdateGameIGDBData struct {
	Name            string
	PlatformsIDs    []int32
	Websites        []string
	IGDBRating      float64
	IGDBRatingCount int32
}

// UpdateGame - update game fields
type UpdateGame struct {
	Name         *string
	Developer    *string
	Publisher    string
	ReleaseDate  *string
	GenresIDs    *[]int32
	LogoURL      *string
	Summary      *string
	PlatformsIDs *[]int32
	Screenshots  *[]string
	Websites     *[]string
}

// GamesFilter - games filter
type GamesFilter struct {
	Name        string
	DeveloperID int32
	PublisherID int32
	GenreID     int32
	OrderBy     OrderBy
}

// GameTrendingData contains data needed for trending index calculation
type GameTrendingData struct {
	Year            int     `db:"release_year"`
	Month           int     `db:"release_month"`
	IGDBRating      float64 `db:"igdb_rating"`
	IGDBRatingCount int32   `db:"igdb_rating_count"`
	Rating          float64 `db:"rating"`
	RatingCount     int32   `db:"rating_count"`
}

// GetGameSlug - returns game slug by name
func GetGameSlug(name string) string {
	return strings.ReplaceAll(strings.ToLower(strings.ToValidUTF8(name, "")), " ", "-")
}

// MapToUpdateGameData maps UpdateGame and Game to UpdateGameData
func (ug UpdateGame) MapToUpdateGameData(g Game, developersIDs []int32) UpdateGameData {
	update := UpdateGameData{
		Name:          g.Name,
		DevelopersIDs: g.DevelopersIDs,
		PublishersIDs: g.PublishersIDs,
		ReleaseDate:   g.ReleaseDate.String(),
		GenresIDs:     g.GenresIDs,
		LogoURL:       g.LogoURL,
		Summary:       g.Summary,
		Slug:          g.Slug,
		PlatformsIDs:  g.PlatformsIDs,
		Screenshots:   g.Screenshots,
		Websites:      g.Websites,
	}

	update.DevelopersIDs = developersIDs
	update.ModerationStatus = ModerationStatusRecheck

	if ug.Name != nil {
		update.Name = *ug.Name
		update.Slug = GetGameSlug(*ug.Name)
	}
	if ug.ReleaseDate != nil {
		update.ReleaseDate = *ug.ReleaseDate
	}
	if ug.GenresIDs != nil {
		update.GenresIDs = *ug.GenresIDs
	}
	if ug.LogoURL != nil && *ug.LogoURL != "" {
		update.LogoURL = *ug.LogoURL
	}
	if ug.Summary != nil {
		update.Summary = *ug.Summary
	}
	if ug.PlatformsIDs != nil {
		update.PlatformsIDs = *ug.PlatformsIDs
	}
	if ug.Screenshots != nil {
		update.Screenshots = *ug.Screenshots
	}
	if ug.Websites != nil {
		update.Websites = *ug.Websites
	}

	return update
}

// MapToCreateGameData maps CreateGame to CreateGameData
func (cg CreateGame) MapToCreateGameData(publisherID, developerID int32) CreateGameData {
	return CreateGameData{
		Name:             cg.Name,
		DevelopersIDs:    []int32{developerID},
		PublishersIDs:    []int32{publisherID},
		ReleaseDate:      cg.ReleaseDate,
		GenresIDs:        cg.GenresIDs,
		LogoURL:          cg.LogoURL,
		Summary:          cg.Summary,
		Slug:             cg.Slug,
		PlatformsIDs:     cg.PlatformsIDs,
		Screenshots:      cg.Screenshots,
		Websites:         cg.Websites,
		ModerationStatus: ModerationStatusCheck,
	}
}
