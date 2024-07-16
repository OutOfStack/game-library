package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Respond marshals a value to JSON and sends it to client
func Respond(c *gin.Context, val interface{}, statusCode int) {
	if statusCode == http.StatusNoContent {
		c.Writer.WriteHeader(statusCode)
		return
	}

	if val == nil {
		val = struct{}{}
	}
	data, err := json.Marshal(val)
	if err != nil {
		Err(c, fmt.Errorf("marshalling value to json: %w", err))
		return
	}
	c.Header("content-type", "application/json;charset=utf-8")
	c.Status(statusCode)
	_, err = c.Writer.Write(data)
	if err != nil {
		Err(c, fmt.Errorf("writing to client: %v", err))
		return
	}
}

// RespondError handles outgoing errors
func RespondError(c *gin.Context, err error) {
	var webErr *Error
	if ok := errors.As(err, &webErr); ok {
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
