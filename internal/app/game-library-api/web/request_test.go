package web_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/api/validation"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/appconf"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// Mocked structures for testing
type ReqBody struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (r *ReqBody) ValidateWith(v *validation.Validator) (bool, []web.FieldError) {
	var validationErrors []web.FieldError

	if r.Name == "" {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "name",
			Error: v.ErrRequiredMsg(),
		})
	}
	if !strings.Contains(r.Email, "@") {
		validationErrors = append(validationErrors, web.FieldError{
			Field: "email",
			Error: "invalid email",
		})
	}

	return len(validationErrors) == 0, validationErrors
}

func (r *ReqBody) Sanitize() {
	r.Name = strings.TrimSpace(r.Name)
	r.Email = strings.TrimSpace(r.Email)
}

func TestDecodeChi_Success(t *testing.T) {
	decoder := web.NewDecoder(zap.NewNop(), &appconf.Cfg{})
	email := td.Email()
	input := ReqBody{
		Name:  td.String(),
		Email: " " + email + " ",
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	var decoded ReqBody
	err := decoder.Decode(req, &decoded)

	require.NoError(t, err)
	require.Equal(t, input.Name, decoded.Name)
	require.Equal(t, email, decoded.Email)
}

func TestDecodeChi_InvalidJSON(t *testing.T) {
	decoder := web.NewDecoder(zap.NewNop(), &appconf.Cfg{})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("{invalid_json}")))
	req.Header.Set("Content-Type", "application/json")

	var decoded ReqBody
	err := decoder.Decode(req, &decoded)

	require.Error(t, err)
	require.IsType(t, &web.Error{}, err)
	require.Equal(t, http.StatusBadRequest, err.(*web.Error).StatusCode) //nolint
}

func TestDecodeChi_ValidationError(t *testing.T) {
	decoder := web.NewDecoder(zap.NewNop(), &appconf.Cfg{})
	input := ReqBody{
		Name:  "",
		Email: "invalid-email",
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	var decoded ReqBody
	err := decoder.Decode(req, &decoded)

	require.Error(t, err)
	require.IsType(t, &web.Error{}, err)

	var validationErr *web.Error
	ok := errors.As(err, &validationErr)
	require.True(t, ok)
	require.Equal(t, http.StatusBadRequest, validationErr.StatusCode)
	require.Len(t, validationErr.Fields, 2)

	require.Equal(t, "name", validationErr.Fields[0].Field)
	require.Contains(t, validationErr.Fields[0].Error, "required")
	require.Equal(t, "email", validationErr.Fields[1].Field)
	require.Contains(t, validationErr.Fields[1].Error, "invalid email")
}

func TestDecodeChi_ValidationErrorEmptyBody(t *testing.T) {
	decoder := web.NewDecoder(zap.NewNop(), &appconf.Cfg{})
	input := ReqBody{} // Empty struct triggers validation errors
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	var decoded ReqBody
	err := decoder.Decode(req, &decoded)

	// Ensure the error is of type *Error
	require.Error(t, err)
	require.IsType(t, &web.Error{}, err)

	var validationErr *web.Error
	ok := errors.As(err, &validationErr)
	require.True(t, ok)
	require.Equal(t, http.StatusBadRequest, validationErr.StatusCode)

	// Ensure fields contain validation errors
	require.Len(t, validationErr.Fields, 2)
	require.Equal(t, "name", validationErr.Fields[0].Field)
	require.Contains(t, validationErr.Fields[0].Error, "required")
	require.Equal(t, "email", validationErr.Fields[1].Field)
	require.Contains(t, strings.ToLower(validationErr.Fields[1].Error), "email")
}
