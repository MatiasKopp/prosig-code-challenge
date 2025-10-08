package httputil

import (
	"encoding/json"
	"net/http"
)

// HTTPErrorResponse
type HTTPErrorResponse struct {
	Message string `json:"message"`
	Cause   string `json:"cause"`
}

// handlerHTTPError Translates service errors into HTTP errors.
func HandlerHTTPError(w http.ResponseWriter, msg string, err error, errMapper map[error]int) {
	httpErr := HTTPErrorResponse{
		Message: msg,
		Cause:   err.Error(),
	}

	data, err := json.Marshal(httpErr)
	if err != nil {
		// TODO: Log warning here.
		return
	}

	statusCode, exists := errMapper[err]
	if !exists {
		statusCode = http.StatusInternalServerError
	}

	w.Header().Add("Content-Type", "application/problem+json")
	w.WriteHeader(statusCode)
	w.Write(data)
}

// handlerHTTPResponse Translates service errors into HTTP errors.
func HandlerHTTPResponse(w http.ResponseWriter, statusCode int, response any) {
	w.WriteHeader(statusCode)

	if response != nil {
		data, err := json.Marshal(response)
		if err != nil {
			// TODO: Log warning here.
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(data)
		return
	}

}
