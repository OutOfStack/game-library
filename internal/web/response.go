package web

import (
	"encoding/json"
	"errors"
	"net/http"
)

// Respond marshals a value to JSON and writes it to response
func Respond(w http.ResponseWriter, val interface{}, statusCode int) {
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return
	}

	if val == nil {
		val = struct{}{}
	}
	data, err := json.Marshal(val)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		data, _ = json.Marshal(ErrorResponse{
			Error: "marshal value to json:" + err.Error(),
		})
		_, _ = w.Write(data)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(statusCode)

	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		data, _ = json.Marshal(ErrorResponse{
			Error: "write to client: " + err.Error(),
		})
		_, _ = w.Write(data)
		return
	}
}

// RespondError handles outgoing errors by marshaling an error response.
// if error is not of type Error, then err is ignored.
// if no err is provided, InternalServerError returns to client
func RespondError(w http.ResponseWriter, err error) {
	var statusCode = http.StatusInternalServerError
	response := ErrorResponse{
		Error: http.StatusText(statusCode),
	}

	var webErr *Error
	if errors.As(err, &webErr) {
		statusCode = webErr.StatusCode
		errMsg := webErr.Err.Error()
		if statusCode >= 500 {
			errMsg = http.StatusText(statusCode)
		}
		response = ErrorResponse{
			Error:  errMsg,
			Fields: webErr.Fields,
		}
	}

	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(statusCode)
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// Respond500 returns InternalServerError to client
func Respond500(w http.ResponseWriter) {
	RespondError(w, nil)
}
