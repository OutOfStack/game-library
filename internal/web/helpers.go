package web

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GetIDParam returns url id param
func GetIDParam(r *http.Request) (int32, error) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil || id <= 0 {
		return 0, NewErrorFromMessage("invalid id", http.StatusBadRequest)
	}
	return int32(id), nil
}
