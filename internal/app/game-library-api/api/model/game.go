package model

import (
	"strings"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/api/validation"
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

// ValidateWith validates CreateGameRequest
func (r *CreateGameRequest) ValidateWith(v *validation.Validator) (bool, []web.FieldError) {
	var validationErrors []web.FieldError

	if r.Name == "" {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "name",
			Error: v.ErrRequiredMsg(),
		})
	}

	if r.Developer == "" {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "developer",
			Error: v.ErrRequiredMsg(),
		})
	}

	if r.ReleaseDate == "" {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "releaseDate",
			Error: v.ErrRequiredMsg(),
		})
	} else if !v.ValidateDate(r.ReleaseDate) {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "releaseDate",
			Error: v.ErrInvalidDateMsg(),
		})
	}

	if len(r.GenresIDs) == 0 {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "genresIds",
			Error: v.ErrRequiredMsg(),
		})
	} else if !v.ValidatePositive(r.GenresIDs) {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "genresIds",
			Error: v.ErrNonPositiveValuesMsg(),
		})
	}

	if r.LogoURL == "" {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "logoUrl",
			Error: v.ErrRequiredMsg(),
		})
	} else if !v.ValidateImageURLs([]string{r.LogoURL}) {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "logoUrl",
			Error: v.ErrInvalidImageURLMsg(),
		})
	}

	if r.Summary == "" {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "summary",
			Error: v.ErrRequiredMsg(),
		})
	}

	if len(r.PlatformsIDs) == 0 {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "platformsIds",
			Error: v.ErrRequiredMsg(),
		})
	} else if !v.ValidatePositive(r.PlatformsIDs) {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "platformsIds",
			Error: v.ErrNonPositiveValuesMsg(),
		})
	}

	if len(r.Screenshots) == 0 {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "screenshots",
			Error: v.ErrRequiredMsg(),
		})
	} else if !v.ValidateImageURLs(r.Screenshots) {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "screenshots",
			Error: v.ErrInvalidImageURLsMsg(),
		})
	}

	if !v.ValidateWebsiteURLs(r.Websites) {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "websites",
			Error: v.ErrInvalidWebsitesURLMsg(),
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
	r.GenresIDs = validation.RemoveDuplicates(r.GenresIDs)
	r.PlatformsIDs = validation.RemoveDuplicates(r.PlatformsIDs)
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

// ValidateWith validates UpdateGameRequest
func (r *UpdateGameRequest) ValidateWith(v *validation.Validator) (bool, []web.FieldError) {
	var validationErrors []web.FieldError

	if r.Name != nil && *r.Name == "" {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "name",
			Error: v.ErrRequiredMsg(),
		})
	}

	if r.Developer != nil && *r.Developer == "" {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "developer",
			Error: v.ErrRequiredMsg(),
		})
	}

	if r.ReleaseDate != nil {
		if *r.ReleaseDate == "" {
			validationErrors = append(validationErrors, web.FieldError{
				Field: "releaseDate",
				Error: v.ErrRequiredMsg(),
			})
		} else if !v.ValidateDate(*r.ReleaseDate) {
			validationErrors = append(validationErrors, web.FieldError{
				Field: "releaseDate",
				Error: v.ErrInvalidDateMsg(),
			})
		}
	}

	if r.GenresIDs != nil {
		if len(*r.GenresIDs) == 0 {
			validationErrors = append(validationErrors, web.FieldError{
				Field: "genresIds",
				Error: v.ErrRequiredMsg(),
			})
		} else if !v.ValidatePositive(*r.GenresIDs) {
			validationErrors = append(validationErrors, web.FieldError{
				Field: "genresIds",
				Error: v.ErrNonPositiveValuesMsg(),
			})
		}
	}

	if r.LogoURL != nil {
		if *r.LogoURL == "" {
			validationErrors = append(validationErrors, web.FieldError{
				Field: "logoUrl",
				Error: v.ErrRequiredMsg(),
			})
		} else if !v.ValidateImageURLs([]string{*r.LogoURL}) {
			validationErrors = append(validationErrors, web.FieldError{
				Field: "logoUrl",
				Error: v.ErrInvalidImageURLMsg(),
			})
		}
	}

	if r.Summary != nil && *r.Summary == "" {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "summary",
			Error: v.ErrRequiredMsg(),
		})
	}

	if r.PlatformsIDs != nil {
		if len(*r.PlatformsIDs) == 0 {
			validationErrors = append(validationErrors, web.FieldError{
				Field: "platformsIds",
				Error: v.ErrRequiredMsg(),
			})
		} else if !v.ValidatePositive(*r.PlatformsIDs) {
			validationErrors = append(validationErrors, web.FieldError{
				Field: "platformsIds",
				Error: v.ErrNonPositiveValuesMsg(),
			})
		}
	}

	if r.Screenshots != nil {
		if len(*r.Screenshots) == 0 {
			validationErrors = append(validationErrors, web.FieldError{
				Field: "screenshots",
				Error: v.ErrRequiredMsg(),
			})
		} else if !v.ValidateImageURLs(*r.Screenshots) {
			validationErrors = append(validationErrors, web.FieldError{
				Field: "screenshots",
				Error: v.ErrInvalidImageURLsMsg(),
			})
		}
	}

	if r.Websites != nil && !v.ValidateWebsiteURLs(*r.Websites) {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "websites",
			Error: v.ErrInvalidWebsitesURLMsg(),
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
		*r.GenresIDs = validation.RemoveDuplicates(*r.GenresIDs)
	}
	if r.PlatformsIDs != nil {
		*r.PlatformsIDs = validation.RemoveDuplicates(*r.PlatformsIDs)
	}
}
