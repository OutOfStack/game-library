package model

import (
	"strings"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/microcosm-cc/bluemonday"
)

// GetGamesQueryParams - get games query params
type GetGamesQueryParams struct {
	PageSize    int    `form:"pageSize"`
	Page        int    `form:"page"`
	OrderBy     string `form:"orderBy"`
	Name        string `form:"name"`
	GenreID     int32  `form:"genre"`
	DeveloperID int32  `form:"developer"`
	PublisherID int32  `form:"publisher"`
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
	Name         string   `json:"name"`
	Developer    string   `json:"developer"`
	ReleaseDate  string   `json:"releaseDate"`
	GenresIDs    []int32  `json:"genresIds"`
	LogoURL      string   `json:"logoUrl"`
	Summary      string   `json:"summary"`
	PlatformsIDs []int32  `json:"platformsIds"`
	Screenshots  []string `json:"screenshots"`
	Websites     []string `json:"websites"`
}

// Validate validates CreateGameRequest
func (r *CreateGameRequest) Validate() (bool, []web.FieldError) {
	var validationErrors []web.FieldError

	if r.Name == "" {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "name",
			Error: ErrRequiredMsg,
		})
	}

	if r.Developer == "" {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "developer",
			Error: ErrRequiredMsg,
		})
	}

	if r.ReleaseDate == "" {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "releaseDate",
			Error: ErrRequiredMsg,
		})
	} else if !validateDate(r.ReleaseDate) {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "releaseDate",
			Error: ErrInvalidDateMsg,
		})
	}

	if len(r.GenresIDs) == 0 {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "genresIds",
			Error: ErrRequiredMsg,
		})
	} else if !validatePositive(r.GenresIDs) {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "genresIds",
			Error: ErrNonPositiveValuesMsg,
		})
	}

	if r.LogoURL == "" {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "logoUrl",
			Error: ErrRequiredMsg,
		})
	} else if !validateImageURLs([]string{r.LogoURL}) {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "logoUrl",
			Error: ErrInvalidImageURLMsg,
		})
	}

	if r.Summary == "" {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "summary",
			Error: ErrRequiredMsg,
		})
	}

	if len(r.PlatformsIDs) == 0 {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "platformsIds",
			Error: ErrRequiredMsg,
		})
	} else if !validatePositive(r.PlatformsIDs) {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "platformsIds",
			Error: ErrNonPositiveValuesMsg,
		})
	}

	if len(r.Screenshots) == 0 {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "screenshots",
			Error: ErrRequiredMsg,
		})
	} else if !validateImageURLs(r.Screenshots) {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "screenshots",
			Error: ErrInvalidImageURLsMsg,
		})
	}

	if !validateWebsiteURLs(r.Websites) {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "websites",
			Error: ErrInvalidWebsitesURLMsg,
		})
	}

	return len(validationErrors) == 0, validationErrors
}

// Sanitize cleans up user input for CreateGameRequest
func (r *CreateGameRequest) Sanitize() {
	p := bluemonday.StrictPolicy()

	// sanitize strings
	r.Name = strings.TrimSpace(p.Sanitize(r.Name))
	r.Developer = strings.TrimSpace(p.Sanitize(r.Developer))
	r.Summary = strings.TrimSpace(p.Sanitize(r.Summary))

	// remove duplicates
	r.GenresIDs = removeDuplicates(r.GenresIDs)
	r.PlatformsIDs = removeDuplicates(r.PlatformsIDs)
}

// UpdateGameRequest - update game request. All fields are optional
type UpdateGameRequest struct {
	Name         *string   `json:"name"`
	Developer    *string   `json:"developer"`
	ReleaseDate  *string   `json:"releaseDate"`
	GenresIDs    *[]int32  `json:"genresIds"`
	LogoURL      *string   `json:"logoUrl"`
	Summary      *string   `json:"summary"`
	PlatformsIDs *[]int32  `json:"platformsIds"`
	Screenshots  *[]string `json:"screenshots"`
	Websites     *[]string `json:"websites"`
}

// Validate validates UpdateGameRequest
func (r *UpdateGameRequest) Validate() (bool, []web.FieldError) {
	var validationErrors []web.FieldError

	if r.Name != nil && *r.Name == "" {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "name",
			Error: ErrRequiredMsg,
		})
	}

	if r.Developer != nil && *r.Developer == "" {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "developer",
			Error: ErrRequiredMsg,
		})
	}

	if r.ReleaseDate != nil {
		if *r.ReleaseDate == "" {
			validationErrors = append(validationErrors, web.FieldError{
				Field: "releaseDate",
				Error: ErrRequiredMsg,
			})
		} else if !validateDate(*r.ReleaseDate) {
			validationErrors = append(validationErrors, web.FieldError{
				Field: "releaseDate",
				Error: ErrInvalidDateMsg,
			})
		}
	}

	if r.GenresIDs != nil {
		if len(*r.GenresIDs) == 0 {
			validationErrors = append(validationErrors, web.FieldError{
				Field: "genresIds",
				Error: ErrRequiredMsg,
			})
		} else if !validatePositive(*r.GenresIDs) {
			validationErrors = append(validationErrors, web.FieldError{
				Field: "genresIds",
				Error: ErrNonPositiveValuesMsg,
			})
		}
	}

	if r.LogoURL != nil {
		if *r.LogoURL == "" {
			validationErrors = append(validationErrors, web.FieldError{
				Field: "logoUrl",
				Error: ErrRequiredMsg,
			})
		} else if !validateImageURLs([]string{*r.LogoURL}) {
			validationErrors = append(validationErrors, web.FieldError{
				Field: "logoUrl",
				Error: ErrInvalidImageURLMsg,
			})
		}
	}

	if r.Summary != nil && *r.Summary == "" {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "summary",
			Error: ErrRequiredMsg,
		})
	}

	if r.PlatformsIDs != nil {
		if len(*r.PlatformsIDs) == 0 {
			validationErrors = append(validationErrors, web.FieldError{
				Field: "platformsIds",
				Error: ErrRequiredMsg,
			})
		} else if !validatePositive(*r.PlatformsIDs) {
			validationErrors = append(validationErrors, web.FieldError{
				Field: "platformsIds",
				Error: ErrNonPositiveValuesMsg,
			})
		}
	}

	if r.Screenshots != nil {
		if len(*r.Screenshots) == 0 {
			validationErrors = append(validationErrors, web.FieldError{
				Field: "screenshots",
				Error: ErrRequiredMsg,
			})
		} else if !validateImageURLs(*r.Screenshots) {
			validationErrors = append(validationErrors, web.FieldError{
				Field: "screenshots",
				Error: ErrInvalidImageURLsMsg,
			})
		}
	}

	if r.Websites != nil && !validateWebsiteURLs(*r.Websites) {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "websites",
			Error: ErrInvalidWebsitesURLMsg,
		})
	}

	return len(validationErrors) == 0, validationErrors
}

// Sanitize cleans up user input for CreateGameRequest
func (r *UpdateGameRequest) Sanitize() {
	p := bluemonday.StrictPolicy()

	// sanitize strings
	if r.Name != nil {
		*r.Name = strings.TrimSpace(p.Sanitize(*r.Name))
	}
	if r.Developer != nil {
		*r.Developer = strings.TrimSpace(p.Sanitize(*r.Developer))
	}
	if r.Summary != nil {
		*r.Summary = strings.TrimSpace(p.Sanitize(*r.Summary))
	}

	// remove duplicates
	if r.GenresIDs != nil {
		*r.GenresIDs = removeDuplicates(*r.GenresIDs)
	}
	if r.PlatformsIDs != nil {
		*r.PlatformsIDs = removeDuplicates(*r.PlatformsIDs)
	}
}
