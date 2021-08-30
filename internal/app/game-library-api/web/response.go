package web

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// Respond marshals a value to JSON and sends it to client
func Respond(c *gin.Context, val interface{}, statusCode int) {
	if statusCode == http.StatusNoContent {
		c.Writer.WriteHeader(statusCode)
		return
	}

	data, err := json.Marshal(val)
	if err != nil {
		c.Error(errors.Wrap(err, "marshalling value to json"))
		return
	}
	c.Header("content-type", "application/json;charset=utf-8")
	c.Status(statusCode)
	_, err = c.Writer.Write(data)
	if err != nil {
		c.Error(errors.Wrap(err, "writing to client"))
		return
	}
}

// RespondError handles outgoing errors
func RespondError(c *gin.Context, err error) {
	webErr, ok := errors.Cause(err).(*Error)
	if ok {
		response := ErrorResponse{
			Error:  webErr.Err.Error(),
			Fields: webErr.Fields,
		}
		Respond(c, response, webErr.Status)
		return
	}

	response := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}
	Respond(c, response, http.StatusInternalServerError)
}
