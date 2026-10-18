package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/OutOfStack/game-library/internal/api/validation"
	"github.com/OutOfStack/game-library/internal/web"
	"github.com/microcosm-cc/bluemonday"
)

const (
	fieldName         = "name"
	fieldDeveloper    = "developer"
	fieldReleaseDate  = "releaseDate"
	fieldGenresIDs    = "genresIds"
	fieldLogoURL      = "logoUrl"
	fieldPlatformsIDs = "platformsIds"
	fieldScreenshots  = "screenshots"
)

// validatable is a request with validation rules
type validatable interface {
	ValidateWith(rules *validation.Rules) (bool, []web.FieldError)
}

// sanitizable is a request that cleans up its input
type sanitizable interface {
	Sanitize()
}

// validateRequest checks request against its validation rules and sanitizes it if valid
func (p *Provider) validateRequest(req validatable) error {
	if valid, fields := req.ValidateWith(p.rules); !valid {
		return web.NewErrorWithFields(errors.New("validation error"), http.StatusBadRequest, fields)
	}
	if s, ok := req.(sanitizable); ok {
		s.Sanitize()
	}
	return nil
}

// ValidateWith validates CreateGameRequest
func (r *CreateGameRequest) ValidateWith(rules *validation.Rules) (bool, []web.FieldError) {
	var validationErrors []web.FieldError

	if r.Name == "" {
		validationErrors = append(validationErrors, web.FieldError{
			Field: fieldName,
			Error: rules.ErrRequiredMsg(),
		})
	}

	if r.Developer == "" {
		validationErrors = append(validationErrors, web.FieldError{
			Field: fieldDeveloper,
			Error: rules.ErrRequiredMsg(),
		})
	}

	if r.ReleaseDate == "" {
		validationErrors = append(validationErrors, web.FieldError{
			Field: fieldReleaseDate,
			Error: rules.ErrRequiredMsg(),
		})
	} else if !rules.ValidateDate(r.ReleaseDate) {
		validationErrors = append(validationErrors, web.FieldError{
			Field: fieldReleaseDate,
			Error: rules.ErrInvalidDateMsg(),
		})
	}

	if len(r.GenresIDs) == 0 {
		validationErrors = append(validationErrors, web.FieldError{
			Field: fieldGenresIDs,
			Error: rules.ErrRequiredMsg(),
		})
	} else if !rules.ValidatePositive(r.GenresIDs) {
		validationErrors = append(validationErrors, web.FieldError{
			Field: fieldGenresIDs,
			Error: rules.ErrNonPositiveValuesMsg(),
		})
	}

	if r.LogoURL == "" {
		validationErrors = append(validationErrors, web.FieldError{
			Field: fieldLogoURL,
			Error: rules.ErrRequiredMsg(),
		})
	} else if !rules.ValidateImageURLs([]string{r.LogoURL}) {
		validationErrors = append(validationErrors, web.FieldError{
			Field: fieldLogoURL,
			Error: rules.ErrInvalidImageURLMsg(),
		})
	}

	if r.Summary == "" {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "summary",
			Error: rules.ErrRequiredMsg(),
		})
	}

	if len(r.PlatformsIDs) == 0 {
		validationErrors = append(validationErrors, web.FieldError{
			Field: fieldPlatformsIDs,
			Error: rules.ErrRequiredMsg(),
		})
	} else if !rules.ValidatePositive(r.PlatformsIDs) {
		validationErrors = append(validationErrors, web.FieldError{
			Field: fieldPlatformsIDs,
			Error: rules.ErrNonPositiveValuesMsg(),
		})
	}

	if len(r.Screenshots) == 0 {
		validationErrors = append(validationErrors, web.FieldError{
			Field: fieldScreenshots,
			Error: rules.ErrRequiredMsg(),
		})
	} else if !rules.ValidateImageURLs(r.Screenshots) {
		validationErrors = append(validationErrors, web.FieldError{
			Field: fieldScreenshots,
			Error: rules.ErrInvalidImageURLsMsg(),
		})
	}

	if !rules.ValidateWebsiteURLs(r.Websites) {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "websites",
			Error: rules.ErrInvalidWebsitesURLMsg(),
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

// ValidateWith validates UpdateGameRequest
func (r *UpdateGameRequest) ValidateWith(rules *validation.Rules) (bool, []web.FieldError) {
	var validationErrors []web.FieldError

	if r.Name != nil && *r.Name == "" {
		validationErrors = append(validationErrors, web.FieldError{
			Field: fieldName,
			Error: rules.ErrRequiredMsg(),
		})
	}

	if r.Developer != nil && *r.Developer == "" {
		validationErrors = append(validationErrors, web.FieldError{
			Field: fieldDeveloper,
			Error: rules.ErrRequiredMsg(),
		})
	}

	if r.ReleaseDate != nil {
		if *r.ReleaseDate == "" {
			validationErrors = append(validationErrors, web.FieldError{
				Field: fieldReleaseDate,
				Error: rules.ErrRequiredMsg(),
			})
		} else if !rules.ValidateDate(*r.ReleaseDate) {
			validationErrors = append(validationErrors, web.FieldError{
				Field: fieldReleaseDate,
				Error: rules.ErrInvalidDateMsg(),
			})
		}
	}

	if r.GenresIDs != nil {
		if len(*r.GenresIDs) == 0 {
			validationErrors = append(validationErrors, web.FieldError{
				Field: fieldGenresIDs,
				Error: rules.ErrRequiredMsg(),
			})
		} else if !rules.ValidatePositive(*r.GenresIDs) {
			validationErrors = append(validationErrors, web.FieldError{
				Field: fieldGenresIDs,
				Error: rules.ErrNonPositiveValuesMsg(),
			})
		}
	}

	if r.LogoURL != nil {
		if *r.LogoURL == "" {
			validationErrors = append(validationErrors, web.FieldError{
				Field: fieldLogoURL,
				Error: rules.ErrRequiredMsg(),
			})
		} else if !rules.ValidateImageURLs([]string{*r.LogoURL}) {
			validationErrors = append(validationErrors, web.FieldError{
				Field: fieldLogoURL,
				Error: rules.ErrInvalidImageURLMsg(),
			})
		}
	}

	if r.Summary != nil && *r.Summary == "" {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "summary",
			Error: rules.ErrRequiredMsg(),
		})
	}

	if r.PlatformsIDs != nil {
		if len(*r.PlatformsIDs) == 0 {
			validationErrors = append(validationErrors, web.FieldError{
				Field: fieldPlatformsIDs,
				Error: rules.ErrRequiredMsg(),
			})
		} else if !rules.ValidatePositive(*r.PlatformsIDs) {
			validationErrors = append(validationErrors, web.FieldError{
				Field: fieldPlatformsIDs,
				Error: rules.ErrNonPositiveValuesMsg(),
			})
		}
	}

	if r.Screenshots != nil {
		if len(*r.Screenshots) == 0 {
			validationErrors = append(validationErrors, web.FieldError{
				Field: fieldScreenshots,
				Error: rules.ErrRequiredMsg(),
			})
		} else if !rules.ValidateImageURLs(*r.Screenshots) {
			validationErrors = append(validationErrors, web.FieldError{
				Field: fieldScreenshots,
				Error: rules.ErrInvalidImageURLsMsg(),
			})
		}
	}

	if r.Websites != nil && !rules.ValidateWebsiteURLs(*r.Websites) {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "websites",
			Error: rules.ErrInvalidWebsitesURLMsg(),
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

// ValidateWith validates CreateRatingRequest
func (r *CreateRatingRequest) ValidateWith(rules *validation.Rules) (bool, []web.FieldError) {
	var validationErrors []web.FieldError

	if r.Rating > 5 {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "rating",
			Error: rules.ErrInvalidRatingMsg(),
		})
	}

	return len(validationErrors) == 0, validationErrors
}

// ValidateWith validates GetUserRatingsRequest
func (r *GetUserRatingsRequest) ValidateWith(rules *validation.Rules) (bool, []web.FieldError) {
	var validationErrors []web.FieldError

	if !rules.ValidatePositive(r.GameIDs) {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "gameIds",
			Error: rules.ErrNonPositiveValuesMsg(),
		})
	}

	return len(validationErrors) == 0, validationErrors
}
