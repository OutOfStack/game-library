package handler

// GameResponse represents game response
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

// CreateGameRequest represents create game request
type CreateGameRequest struct {
	Name         string   `json:"name" validate:"required"`
	Developer    string   `json:"developer" validate:"required"`
	ReleaseDate  string   `json:"releaseDate" validate:"date"`
	GenresIDs    []int32  `json:"genresIds"`
	LogoURL      string   `json:"logoUrl"`
	Summary      string   `json:"summary"`
	Slug         string   `json:"slug"`
	PlatformsIDs []int32  `json:"platformsIDs"`
	Screenshots  []string `json:"screenshots"`
	Websites     []string `json:"websites"`
}

// UpdateGameRequest represents update game request
// All fields are optional
type UpdateGameRequest struct {
	Name        *string   `json:"name"`
	Developer   *string   `json:"developer" validate:"omitempty"`
	ReleaseDate *string   `json:"releaseDate" validate:"omitempty,date"`
	GenresIDs   *[]int32  `json:"genresIds" validate:"omitempty"`
	LogoURL     *string   `json:"logoUrl"`
	Summary     *string   `json:"summary"`
	Slug        *string   `json:"slug"`
	Platforms   *[]int32  `json:"platforms"`
	Screenshots *[]string `json:"screenshots"`
	Websites    *[]string `json:"websites"`
}

// CreateRatingRequest represents create rating request
type CreateRatingRequest struct {
	Rating uint8 `json:"rating" validate:"gte=1,lte=5"`
}

// RatingResponse represents rating response
type RatingResponse struct {
	GameID int32  `json:"gameId"`
	UserID string `json:"userId"`
	Rating uint8  `json:"rating"`
}

// GetUserRatingsRequest represents get user ratings request
type GetUserRatingsRequest struct {
	GameIDs []int32 `json:"gameIds"`
}

// IDResponse represents response with id
type IDResponse struct {
	ID int32 `json:"id"`
}

// Genre represents genre response
type Genre struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

// Company represents company response
type Company struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

// Platform represents platform response
type Platform struct {
	ID           int32  `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
}

// CountResponse represent count response
type CountResponse struct {
	Count uint64 `json:"count"`
}
