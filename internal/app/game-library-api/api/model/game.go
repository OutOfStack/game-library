package model

// GetGamesQueryParams - get games query params
type GetGamesQueryParams struct {
	PageSize  int    `form:"pageSize" binding:"required,gt=0"`
	Page      int    `form:"page" binding:"required,gt=0"`
	OrderBy   string `form:"orderBy" binding:"oneof=default name releaseDate"`
	Name      string `form:"name"`
	Genre     int32  `form:"genre"`
	Developer int32  `form:"developer"`
	Publisher int32  `form:"publisher"`
}

// GameResponse - game response
type GameResponse struct {
	ID          int32      `json:"id"`
	Name        string     `json:"name"`
	Developers  []Company  `json:"developers"`
	Publishers  []Company  `json:"publishers"`
	ReleaseDate string     `json:"releaseDate"`
	Genres      []Genre    `json:"genres"`
	LogoURL     string     `json:"logoUrl,omitempty"`
	Rating      float64    `json:"rating"`
	Summary     string     `json:"summary,omitempty"`
	Slug        string     `json:"slug,omitempty"`
	Platforms   []Platform `json:"platforms"`
	Screenshots []string   `json:"screenshots"`
	Websites    []string   `json:"websites"`
}

// GamesResponse - games response
type GamesResponse struct {
	Games []GameResponse `json:"games"`
	Count uint64         `json:"count"`
}

// CreateGameRequest - create game request
type CreateGameRequest struct {
	Name         string   `json:"name" validate:"required"`
	Developer    string   `json:"developer" validate:"required"`
	ReleaseDate  string   `json:"releaseDate" validate:"date"`
	GenresIDs    []int32  `json:"genresIds"`
	LogoURL      string   `json:"logoUrl"`
	Summary      string   `json:"summary"`
	PlatformsIDs []int32  `json:"platformsIDs"`
	Screenshots  []string `json:"screenshots"`
	Websites     []string `json:"websites"`
}

// UpdateGameRequest - update game request. All fields are optional
type UpdateGameRequest struct {
	Name        *string   `json:"name"`
	Developer   *string   `json:"developer" validate:"omitempty"`
	ReleaseDate *string   `json:"releaseDate" validate:"omitempty,date"`
	GenresIDs   *[]int32  `json:"genresIds" validate:"omitempty"`
	LogoURL     *string   `json:"logoUrl"`
	Summary     *string   `json:"summary"`
	Platforms   *[]int32  `json:"platforms"`
	Screenshots *[]string `json:"screenshots"`
	Websites    *[]string `json:"websites"`
}
