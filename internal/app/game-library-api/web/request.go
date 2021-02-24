package web

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Decode unmarshalls JSON request body
func Decode(c *gin.Context, val interface{}) error {
	err := json.NewDecoder(c.Request.Body).Decode(val)
	if err != nil {
		return NewRequestError(err, http.StatusBadRequest)
	}

	return nil
}
