package web

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// Decode unmarshalls JSON request body
func Decode(c *gin.Context, val interface{}) error {
	err := json.NewDecoder(c.Request.Body).Decode(val)
	if err != nil {
		return errors.Wrap(err, "decoding request body")
	}

	return nil
}
