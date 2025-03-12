package model

import (
	"strings"

	"github.com/OutOfStack/game-library/pkg/types"
)

// Game represents database game model
type Game struct {
	ID            int32      `db:"id"`
	Name          string     `db:"name"`
	DevelopersIDs []int32    `db:"developers"`
	PublishersIDs []int32    `db:"publishers"`
	ReleaseDate   types.Date `db:"release_date"`
	GenresIDs     []int32    `db:"genres"`
	LogoURL       string     `db:"logo_url"`
	Rating        float64    `db:"rating"`
	Summary       string     `db:"summary"`
	Slug          string     `db:"slug"`
	PlatformsIDs  []int32    `db:"platforms"`
	Screenshots   []string   `db:"screenshots"`
	Websites      []string   `db:"websites"`
	IGDBRating    float64    `db:"igdb_rating"`
	IGDBID        int64      `db:"igdb_id"`
	Weight        float64    `db:"weight"` // Readonly field
}

// CreateGame represents data for creating game
type CreateGame struct {
	Name          string
	DevelopersIDs []int32
	PublishersIDs []int32
	ReleaseDate   string
	GenresIDs     []int32
	LogoURL       string
	Summary       string
	Slug          string
	PlatformsIDs  []int32
	Screenshots   []string
	Websites      []string
	IGDBRating    float64
	IGDBID        int64
	Developer     string // helper field
	Publisher     string // helper field
}

// UpdateGameData represents data for updating game
type UpdateGameData struct {
	Name         string
	Developers   []int32
	Publishers   []int32
	ReleaseDate  string
	Genres       []int32
	LogoURL      string
	Summary      string
	Slug         string
	PlatformsIDs []int32
	Screenshots  []string
	Websites     []string
	IGDBRating   float64
	IGDBID       int64
}

// UpdatedGame - updated game field
type UpdatedGame struct {
	Name         *string
	Developer    *string
	ReleaseDate  *string
	GenresIDs    *[]int32
	LogoURL      *string
	Summary      *string
	PlatformsIDs *[]int32
	Screenshots  *[]string
	Websites     *[]string
}

// GamesFilter games filter
type GamesFilter struct {
	Name        string
	DeveloperID int32
	PublisherID int32
	GenreID     int32
	OrderBy     OrderBy
}

// GetGameSlug - returns game slug by name
func GetGameSlug(name string) string {
	return strings.ReplaceAll(strings.ToLower(strings.ToValidUTF8(name, "")), " ", "-")
}

// MapToUpdateGameData maps Game to UpdateGateData
func (g Game) MapToUpdateGameData(upd UpdatedGame) UpdateGameData {
	update := UpdateGameData{
		Name:         g.Name,
		Developers:   g.DevelopersIDs,
		Publishers:   g.PublishersIDs,
		ReleaseDate:  g.ReleaseDate.String(),
		Genres:       g.GenresIDs,
		LogoURL:      g.LogoURL,
		Summary:      g.Summary,
		Slug:         g.Slug,
		PlatformsIDs: g.PlatformsIDs,
		Screenshots:  g.Screenshots,
		Websites:     g.Websites,
		IGDBRating:   g.IGDBRating,
		IGDBID:       g.IGDBID,
	}

	if upd.Name != nil {
		update.Name = *upd.Name
		update.Slug = GetGameSlug(*upd.Name)
	}
	if upd.ReleaseDate != nil {
		update.ReleaseDate = *upd.ReleaseDate
	}
	if upd.GenresIDs != nil {
		update.Genres = *upd.GenresIDs
	}
	if upd.LogoURL != nil && *upd.LogoURL != "" {
		update.LogoURL = *upd.LogoURL
	}
	if upd.Summary != nil {
		update.Summary = *upd.Summary
	}
	if upd.PlatformsIDs != nil {
		update.PlatformsIDs = *upd.PlatformsIDs
	}
	if upd.Screenshots != nil {
		update.Screenshots = *upd.Screenshots
	}
	if upd.Websites != nil {
		update.Websites = *upd.Websites
	}

	return update
}
