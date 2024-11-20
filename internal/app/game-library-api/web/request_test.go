package web_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

// Mocked structures for testing
type ReqBody struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type MockTranslator struct{}

func (m *MockTranslator) Translate(err validator.FieldError) string {
	return err.Error() // Simple mock translation
}

func TestDecodeChi_Success(t *testing.T) {
	input := ReqBody{
		Name:  td.String(),
		Email: td.Email(),
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	var decoded ReqBody
	err := web.Decode(req, &decoded)

	require.NoError(t, err)
	require.Equal(t, input, decoded)
}

func TestDecodeChi_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("{invalid_json}")))
	req.Header.Set("Content-Type", "application/json")

	var decoded ReqBody
	err := web.Decode(req, &decoded)

	require.Error(t, err)
	require.IsType(t, &web.Error{}, err)
	require.Equal(t, http.StatusBadRequest, err.(*web.Error).StatusCode) //nolint
}

func TestDecodeChi_ValidationError(t *testing.T) {
	input := ReqBody{
		Name:  "",
		Email: "invalid-email",
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	var decoded ReqBody
	err := web.Decode(req, &decoded)

	require.Error(t, err)
	require.IsType(t, &web.Error{}, err)

	var validationErr *web.Error
	ok := errors.As(err, &validationErr)
	require.True(t, ok)
	require.Equal(t, http.StatusBadRequest, validationErr.StatusCode)
	require.Len(t, validationErr.Fields, 2)

	require.Equal(t, "Name", validationErr.Fields[0].Field)
	require.Contains(t, validationErr.Fields[0].Error, "required")
	require.Equal(t, "Email", validationErr.Fields[1].Field)
	require.Contains(t, validationErr.Fields[1].Error, "email")
}

func TestDecodeChi_ValidationErrorEmptyBody(t *testing.T) {
	input := ReqBody{} // Empty struct triggers validation errors
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	var decoded ReqBody
	err := web.Decode(req, &decoded)

	// Ensure the error is of type *Error
	require.Error(t, err)
	require.IsType(t, &web.Error{}, err)

	var validationErr *web.Error
	ok := errors.As(err, &validationErr)
	require.True(t, ok)
	require.Equal(t, http.StatusBadRequest, validationErr.StatusCode)

	// Ensure fields contain validation errors
	require.Len(t, validationErr.Fields, 2)
	require.Equal(t, "Name", validationErr.Fields[0].Field)
	require.Contains(t, validationErr.Fields[0].Error, "required")
	require.Equal(t, "Email", validationErr.Fields[1].Field)
	require.Contains(t, strings.ToLower(validationErr.Fields[1].Error), "email")
}
