package web

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// Respond marshals a value to JSON and sends it to client
func Respond(c *gin.Context, val interface{}, statusCode int) error {
	if statusCode == http.StatusNoContent {
		c.Writer.WriteHeader(statusCode)
		return nil
	}

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

// RespondError handles outgoing errors
func RespondError(c *gin.Context, err error) error {
	webErr, ok := errors.Cause(err).(*Error)
	if ok {
		response := ErrorResponse{
			Error:  webErr.Err.Error(),
			Fields: webErr.Fields,
		}
		if err = Respond(c, response, webErr.Status); err != nil {
			return err
		}
		return nil
	}

	response := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}
	if err := Respond(c, response, http.StatusInternalServerError); err != nil {
		return err
	}
	return nil
}
