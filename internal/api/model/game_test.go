package model_test

import (
	"reflect"
	"testing"

	"github.com/OutOfStack/game-library/internal/api/model"
	"github.com/OutOfStack/game-library/internal/api/validation"
	"github.com/OutOfStack/game-library/internal/appconf"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCreateGameRequestValidation(t *testing.T) {
	cfg := getCfg()
	cfg.S3.CDNBaseURL = td.URL()
	log := zap.NewNop()
	v := validation.NewValidator(log, cfg)

	validImageURL := cfg.S3.CDNBaseURL + "/" + td.String() + ".jpg"
	validWebsiteURL := "https://steampowered.com/game/" + td.String()

	t.Run("Valid CreateGameRequest", func(t *testing.T) {
		request := model.CreateGameRequest{
			Name:         "Test Game",
			Developer:    "Test Developer",
			ReleaseDate:  "2023-01-01",
			GenresIDs:    []int32{1, 2, 3},
			LogoURL:      validImageURL,
			Summary:      "Test summary",
			PlatformsIDs: []int32{1, 2},
			Screenshots:  []string{validImageURL, validImageURL},
			Websites:     []string{validWebsiteURL},
		}

		valid, errors := request.ValidateWith(v)
		require.True(t, valid, "Expected valid request")
		require.Empty(t, errors, "Expected no validation errors")
	})

	t.Run("Missing Required Fields", func(t *testing.T) {
		request := model.CreateGameRequest{
			// All fields empty
		}

		valid, errors := request.ValidateWith(v)
		require.False(t, valid, "Expected invalid request")
		require.Len(t, errors, 8, "Expected 8 validation errors") // Name, Developer, ReleaseDate, GenresIDs, LogoURL, Summary, PlatformsIDs, Screenshots

		// Check specific error fields
		fields := make(map[string]bool)
		for _, err := range errors {
			fields[err.Field] = true
			if err.Field == "logoURL" {
				require.Equal(t, v.ErrInvalidImageURLMsg(), err.Error)
			} else {
				require.Equal(t, v.ErrRequiredMsg(), err.Error)
			}
		}

		require.True(t, fields["name"], "Expected error for name field")
		require.True(t, fields["developer"], "Expected error for developer field")
		require.True(t, fields["releaseDate"], "Expected error for releaseDate field")
		require.True(t, fields["genresIds"], "Expected error for genresIds field")
		require.True(t, fields["logoUrl"], "Expected error for logoUrl field")
		require.True(t, fields["summary"], "Expected error for summary field")
		require.True(t, fields["platformsIds"], "Expected error for platformsIds field")
		require.True(t, fields["screenshots"], "Expected error for screenshots field")
	})

	t.Run("Invalid Date Format", func(t *testing.T) {
		request := model.CreateGameRequest{
			Name:         "Test Game",
			Developer:    "Test Developer",
			ReleaseDate:  "01/01/2023", // Invalid format
			GenresIDs:    []int32{1},
			LogoURL:      validImageURL,
			Summary:      "Test summary",
			PlatformsIDs: []int32{1},
			Screenshots:  []string{validImageURL},
		}

		valid, errors := request.ValidateWith(v)
		require.False(t, valid, "Expected invalid request")

		hasDateError := false
		for _, err := range errors {
			if err.Field == "releaseDate" && err.Error == v.ErrInvalidDateMsg() {
				hasDateError = true
				break
			}
		}
		require.True(t, hasDateError, "Expected error for invalid date format")
	})

	t.Run("Non-positive Genre IDs", func(t *testing.T) {
		request := model.CreateGameRequest{
			Name:         "Test Game",
			Developer:    "Test Developer",
			ReleaseDate:  "2023-01-01",
			GenresIDs:    []int32{1, 0, 3}, // Contains non-positive ID
			LogoURL:      validImageURL,
			Summary:      "Test summary",
			PlatformsIDs: []int32{1},
			Screenshots:  []string{validImageURL},
		}

		valid, errors := request.ValidateWith(v)
		require.False(t, valid, "Expected invalid request")

		hasGenreError := false
		for _, err := range errors {
			if err.Field == "genresIds" && err.Error == v.ErrNonPositiveValuesMsg() {
				hasGenreError = true
				break
			}
		}
		require.True(t, hasGenreError, "Expected error for non-positive genre IDs")
	})

	t.Run("Invalid Image URL", func(t *testing.T) {
		request := model.CreateGameRequest{
			Name:         "Test Game",
			Developer:    "Test Developer",
			ReleaseDate:  "2023-01-01",
			GenresIDs:    []int32{1},
			LogoURL:      "https://invalid-domain.com/image.jpg", // Invalid domain
			Summary:      "Test summary",
			PlatformsIDs: []int32{1},
			Screenshots:  []string{validImageURL},
		}

		valid, errors := request.ValidateWith(v)
		require.False(t, valid, "Expected invalid request")

		hasLogoError := false
		for _, err := range errors {
			if err.Field == "logoUrl" && err.Error == v.ErrInvalidImageURLMsg() {
				hasLogoError = true
				break
			}
		}
		require.True(t, hasLogoError, "Expected error for invalid logo URL")
	})

	t.Run("Invalid Screenshots", func(t *testing.T) {
		request := model.CreateGameRequest{
			Name:         "Test Game",
			Developer:    "Test Developer",
			ReleaseDate:  "2023-01-01",
			GenresIDs:    []int32{1},
			LogoURL:      validImageURL,
			Summary:      "Test summary",
			PlatformsIDs: []int32{1},
			Screenshots:  []string{validImageURL, "https://invalid-domain.com/image.jpg"}, // One invalid URL
		}

		valid, errors := request.ValidateWith(v)
		require.False(t, valid, "Expected invalid request")

		hasScreenshotError := false
		for _, err := range errors {
			if err.Field == "screenshots" && err.Error == v.ErrInvalidImageURLsMsg() {
				hasScreenshotError = true
				break
			}
		}
		require.True(t, hasScreenshotError, "Expected error for invalid screenshot URLs")
	})

	t.Run("Invalid Website URLs", func(t *testing.T) {
		request := model.CreateGameRequest{
			Name:         "Test Game",
			Developer:    "Test Developer",
			ReleaseDate:  "2023-01-01",
			GenresIDs:    []int32{1},
			LogoURL:      validImageURL,
			Summary:      "Test summary",
			PlatformsIDs: []int32{1},
			Screenshots:  []string{validImageURL},
			Websites:     []string{"https://invalid-domain.com"}, // Invalid domain
		}

		valid, errors := request.ValidateWith(v)
		require.False(t, valid, "Expected invalid request")

		hasWebsiteError := false
		for _, err := range errors {
			if err.Field == "websites" && err.Error == v.ErrInvalidWebsitesURLMsg() {
				hasWebsiteError = true
				break
			}
		}
		require.True(t, hasWebsiteError, "Expected error for invalid website URLs")
	})
}

func TestCreateGameRequestSanitize(t *testing.T) {
	t.Run("Sanitizes HTML content", func(t *testing.T) {
		// Create policy to compare results
		p := bluemonday.StrictPolicy()

		htmlString := createHTMLString()
		expected := p.Sanitize(htmlString)

		request := model.CreateGameRequest{
			Name:         htmlString,
			Developer:    htmlString,
			Summary:      htmlString,
			GenresIDs:    []int32{1, 2, 3},
			PlatformsIDs: []int32{4, 5, 6},
		}

		request.Sanitize()

		require.Equal(t, expected, request.Name, "Name field was not sanitized correctly")
		require.Equal(t, expected, request.Developer, "Developer field was not sanitized correctly")
		require.Equal(t, expected, request.Summary, "Summary field was not sanitized correctly")
	})

	t.Run("Trims whitespace", func(t *testing.T) {
		request := model.CreateGameRequest{
			Name:         "  Game Name with spaces  ",
			Developer:    "\t Developer Name \n",
			Summary:      " \n Summary with whitespace \t ",
			GenresIDs:    []int32{1, 2, 3},
			PlatformsIDs: []int32{4, 5, 6},
		}

		request.Sanitize()

		require.Equal(t, "Game Name with spaces", request.Name, "Name field whitespace was not trimmed")
		require.Equal(t, "Developer Name", request.Developer, "Developer field whitespace was not trimmed")
		require.Equal(t, "Summary with whitespace", request.Summary, "Summary field whitespace was not trimmed")
	})

	t.Run("Removes duplicate IDs", func(t *testing.T) {
		request := model.CreateGameRequest{
			GenresIDs:    []int32{1, 2, 2, 3, 3, 3, 4},
			PlatformsIDs: []int32{5, 5, 6, 7, 7, 8},
		}

		request.Sanitize()

		require.False(t, containsDuplicates(request.GenresIDs), "Duplicates were not removed from GenresIDs")
		require.False(t, containsDuplicates(request.PlatformsIDs), "Duplicates were not removed from PlatformsIDs")

		// Check if all unique values are preserved
		genresMap := make(map[int32]bool)
		for _, id := range []int32{1, 2, 3, 4} {
			genresMap[id] = true
		}

		platformsMap := make(map[int32]bool)
		for _, id := range []int32{5, 6, 7, 8} {
			platformsMap[id] = true
		}

		for _, id := range request.GenresIDs {
			require.True(t, genresMap[id], "Expected genre ID not found after deduplication")
		}

		for _, id := range request.PlatformsIDs {
			require.True(t, platformsMap[id], "Expected platform ID not found after deduplication")
		}
	})

	t.Run("Handles combined sanitization", func(t *testing.T) {
		p := bluemonday.StrictPolicy()
		htmlString := createHTMLString()

		request := model.CreateGameRequest{
			Name:         "  " + htmlString + "  ",
			Developer:    "\t" + htmlString + "\n",
			Summary:      " \n" + htmlString + "\t ",
			GenresIDs:    []int32{1, 1, 2, 2, 3},
			PlatformsIDs: []int32{4, 4, 5, 5, 6},
		}

		request.Sanitize()

		expectedText := p.Sanitize(htmlString)
		require.Equal(t, expectedText, request.Name, "Name not sanitized correctly")
		require.Equal(t, expectedText, request.Developer, "Developer not sanitized correctly")
		require.Equal(t, expectedText, request.Summary, "Summary not sanitized correctly")

		require.Len(t, request.GenresIDs, 3, "GenresIDs not deduplicated correctly")
		require.Len(t, request.PlatformsIDs, 3, "PlatformsIDs not deduplicated correctly")
	})
}

func TestUpdateGameRequestValidation(t *testing.T) {
	cfg := getCfg()
	cfg.S3.CDNBaseURL = td.URL()
	log := zap.NewNop()
	v := validation.NewValidator(log, cfg)

	validImageURL := cfg.S3.CDNBaseURL + "/" + td.String() + ".jpg"
	validWebsiteURL := "https://steampowered.com/game/" + td.String()

	t.Run("Empty UpdateGameRequest is valid", func(t *testing.T) {
		request := model.UpdateGameRequest{
			// All fields nil (valid for update)
		}

		valid, errors := request.ValidateWith(v)
		require.True(t, valid, "Expected valid request with all nil fields")
		require.Empty(t, errors, "Expected no validation errors")
	})

	t.Run("Valid UpdateGameRequest with all fields", func(t *testing.T) {
		request := model.UpdateGameRequest{
			Name:         strPtr("Updated Game"),
			Developer:    strPtr("Updated Developer"),
			ReleaseDate:  strPtr("2023-02-01"),
			GenresIDs:    slicePtr([]int32{1, 2, 3}),
			LogoURL:      strPtr(validImageURL),
			Summary:      strPtr("Updated summary"),
			PlatformsIDs: slicePtr([]int32{1, 2}),
			Screenshots:  slicePtr([]string{validImageURL, validImageURL}),
			Websites:     slicePtr([]string{validWebsiteURL}),
		}

		valid, errors := request.ValidateWith(v)
		require.True(t, valid, "Expected valid request")
		require.Empty(t, errors, "Expected no validation errors")
	})

	t.Run("Empty string values", func(t *testing.T) {
		request := model.UpdateGameRequest{
			Name:      strPtr(""), // Empty string
			Developer: strPtr(""), // Empty string
			Summary:   strPtr(""), // Empty string
		}

		valid, errors := request.ValidateWith(v)
		require.False(t, valid, "Expected invalid request")
		require.Len(t, errors, 3, "Expected 3 validation errors")

		fields := make(map[string]bool)
		for _, err := range errors {
			fields[err.Field] = true
			require.Equal(t, v.ErrRequiredMsg(), err.Error)
		}

		require.True(t, fields["name"], "Expected error for empty name")
		require.True(t, fields["developer"], "Expected error for empty developer")
		require.True(t, fields["summary"], "Expected error for empty summary")
	})

	t.Run("Invalid date format", func(t *testing.T) {
		request := model.UpdateGameRequest{
			ReleaseDate: strPtr("01/01/2023"), // Invalid format
		}

		valid, errors := request.ValidateWith(v)
		require.False(t, valid, "Expected invalid request")

		hasDateError := false
		for _, err := range errors {
			if err.Field == "releaseDate" && err.Error == v.ErrInvalidDateMsg() {
				hasDateError = true
				break
			}
		}
		require.True(t, hasDateError, "Expected error for invalid date format")
	})

	t.Run("Empty arrays", func(t *testing.T) {
		request := model.UpdateGameRequest{
			GenresIDs:    slicePtr([]int32{}),  // Empty array
			PlatformsIDs: slicePtr([]int32{}),  // Empty array
			Screenshots:  slicePtr([]string{}), // Empty array
		}

		valid, errors := request.ValidateWith(v)
		require.False(t, valid, "Expected invalid request")
		require.Len(t, errors, 3, "Expected 3 validation errors")

		fields := make(map[string]bool)
		for _, err := range errors {
			fields[err.Field] = true
			require.Equal(t, v.ErrRequiredMsg(), err.Error)
		}

		require.True(t, fields["genresIds"], "Expected error for empty genresIds")
		require.True(t, fields["platformsIds"], "Expected error for empty platformsIds")
		require.True(t, fields["screenshots"], "Expected error for empty screenshots")
	})

	t.Run("Non-positive IDs", func(t *testing.T) {
		request := model.UpdateGameRequest{
			GenresIDs:    slicePtr([]int32{1, -1, 3}), // Contains negative ID
			PlatformsIDs: slicePtr([]int32{0, 2}),     // Contains zero ID
		}

		valid, errors := request.ValidateWith(v)
		require.False(t, valid, "Expected invalid request")
		require.Len(t, errors, 2, "Expected 2 validation errors")

		fields := make(map[string]bool)
		for _, err := range errors {
			fields[err.Field] = true
			require.Equal(t, v.ErrNonPositiveValuesMsg(), err.Error)
		}

		require.True(t, fields["genresIds"], "Expected error for non-positive genres Ids")
		require.True(t, fields["platformsIds"], "Expected error for non-positive platforms Ids")
	})

	t.Run("Invalid URLs", func(t *testing.T) {
		request := model.UpdateGameRequest{
			LogoURL:     strPtr("https://invalid-domain.com/image.jpg"),
			Screenshots: slicePtr([]string{"https://invalid-domain.com/image.jpg"}),
			Websites:    slicePtr([]string{"https://invalid-domain.com"}),
		}

		valid, errors := request.ValidateWith(v)
		require.False(t, valid, "Expected invalid request")
		require.Len(t, errors, 3, "Expected 3 validation errors")

		fields := make(map[string]string)
		for _, err := range errors {
			fields[err.Field] = err.Error
		}

		require.Equal(t, v.ErrInvalidImageURLMsg(), fields["logoUrl"], "Expected error for invalid logo URL")
		require.Equal(t, v.ErrInvalidImageURLsMsg(), fields["screenshots"], "Expected error for invalid screenshot URLs")
		require.Equal(t, v.ErrInvalidWebsitesURLMsg(), fields["websites"], "Expected error for invalid website URLs")
	})
}

func TestUpdateGameRequest_Sanitize(t *testing.T) {
	tests := []struct {
		name     string
		request  model.UpdateGameRequest
		expected model.UpdateGameRequest
	}{
		{
			name: "sanitize strings and remove duplicates",
			request: model.UpdateGameRequest{
				Name:         strPtr(" Game Name with <script>alert('XSS')</script> "),
				Developer:    strPtr(" Developer <b>Name</b> "),
				Summary:      strPtr(" Game summary with <iframe src='malicious.com'></iframe> "),
				GenresIDs:    slicePtr([]int32{1, 2, 2, 3, 3, 3}),
				PlatformsIDs: slicePtr([]int32{5, 5, 6, 7, 7}),
			},
			expected: model.UpdateGameRequest{
				Name:         strPtr("Game Name with"),
				Developer:    strPtr("Developer Name"),
				Summary:      strPtr("Game summary with"),
				GenresIDs:    slicePtr([]int32{1, 2, 3}),
				PlatformsIDs: slicePtr([]int32{5, 6, 7}),
			},
		},
		{
			name: "nil fields should not cause panic",
			request: model.UpdateGameRequest{
				Name:         nil,
				Developer:    nil,
				Summary:      nil,
				GenresIDs:    nil,
				PlatformsIDs: nil,
			},
			expected: model.UpdateGameRequest{
				Name:         nil,
				Developer:    nil,
				Summary:      nil,
				GenresIDs:    nil,
				PlatformsIDs: nil,
			},
		},
		{
			name: "empty strings should be trimmed",
			request: model.UpdateGameRequest{
				Name:         strPtr("   "),
				Developer:    strPtr("   "),
				Summary:      strPtr("   "),
				GenresIDs:    slicePtr([]int32{}),
				PlatformsIDs: slicePtr([]int32{}),
			},
			expected: model.UpdateGameRequest{
				Name:         strPtr(""),
				Developer:    strPtr(""),
				Summary:      strPtr(""),
				GenresIDs:    slicePtr([]int32{}),
				PlatformsIDs: slicePtr([]int32{}),
			},
		},
		{
			name: "complex HTML should be properly sanitized",
			request: model.UpdateGameRequest{
				Name:      strPtr("<p>Game <strong>Name</strong> with <a href='http://example.com'>Link</a></p>"),
				Developer: strPtr("<div>Developer <em>Studios</em></div>"),
				Summary:   strPtr("<ul><li>Feature 1</li><li>Feature 2</li></ul>"),
			},
			expected: model.UpdateGameRequest{
				Name:      strPtr("Game Name with Link"),
				Developer: strPtr("Developer Studios"),
				Summary:   strPtr("Feature 1Feature 2"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy of the request to avoid modifying the test data
			request := tt.request

			// Call the method being tested
			request.Sanitize()

			// Check string fields
			if !areStringPointersEqual(request.Name, tt.expected.Name) {
				t.Errorf("Name = '%v', expected '%v'", valueOrNil(request.Name), valueOrNil(tt.expected.Name))
			}

			if !areStringPointersEqual(request.Developer, tt.expected.Developer) {
				t.Errorf("Developer = '%v', expected '%v'", valueOrNil(request.Developer), valueOrNil(tt.expected.Developer))
			}

			if !areStringPointersEqual(request.Summary, tt.expected.Summary) {
				t.Errorf("Summary = '%v', expected '%v'", valueOrNil(request.Summary), valueOrNil(tt.expected.Summary))
			}

			// Check slice fields
			if !areInt32SlicePointersEqual(request.GenresIDs, tt.expected.GenresIDs) {
				t.Errorf("GenresIDs = %v, expected %v", valueOrNil(request.GenresIDs), valueOrNil(tt.expected.GenresIDs))
			}

			if !areInt32SlicePointersEqual(request.PlatformsIDs, tt.expected.PlatformsIDs) {
				t.Errorf("PlatformsIDs = %v, expected %v", valueOrNil(request.PlatformsIDs), valueOrNil(tt.expected.PlatformsIDs))
			}
		})
	}
}

// Helper function to get string pointer
func strPtr(s string) *string {
	return &s
}

// Helper function to get slice pointer
func slicePtr[T any](slice []T) *[]T {
	return &slice
}

// Helper to create HTML strings
func createHTMLString() string {
	return "<script>alert('XSS')</script><p>Normal text</p><a href='javascript:alert(\"evil\")'>Click me</a>"
}

// Helper for duplicate removal check
func containsDuplicates(slice []int32) bool {
	seen := make(map[int32]bool)
	for _, v := range slice {
		if seen[v] {
			return true
		}
		seen[v] = true
	}
	return false
}

// Helper functions for comparing pointer values
func areStringPointersEqual(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func areInt32SlicePointersEqual(a, b *[]int32) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	sliceA := *a
	sliceB := *b

	// If both slices are empty, consider them equal
	if len(sliceA) == 0 && len(sliceB) == 0 {
		return true
	}

	// Otherwise, compare with DeepEqual
	return reflect.DeepEqual(sliceA, sliceB)
}

// Helper to display pointer values in error messages
func valueOrNil(ptr interface{}) interface{} {
	if ptr == nil {
		return "nil"
	}
	v := reflect.ValueOf(ptr)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "nil"
		}
		return v.Elem().Interface()
	}
	return ptr
}

func getCfg() *appconf.Cfg {
	return &appconf.Cfg{}
}
