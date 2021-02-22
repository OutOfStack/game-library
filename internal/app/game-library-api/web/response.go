package web

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// Respond marshals a value to JSON and sends it to client
func Respond(c *gin.Context, val interface{}, statusCode int) error {
	data, err := json.Marshal(val)
	if err != nil {
		return errors.Wrap(err, "marshalling value to json")
	}
	c.Header("content-type", "application/json;charset=utf-8")
	c.Status(statusCode)
	_, err = c.Writer.Write(data)
	if err != nil {
		return errors.Wrap(err, "writing to client")
	}

	return nil
}

// RespondError handles otugoing errors
func RespondError(c *gin.Context, err error) error {
	webErr, ok := err.(*Error)
	if ok {
		response := ErrorResponse{
			Error: webErr.Err.Error(),
		}
		return Respond(c, response, webErr.Status)
	}

	response := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}
	return Respond(c, response, http.StatusInternalServerError)
}
