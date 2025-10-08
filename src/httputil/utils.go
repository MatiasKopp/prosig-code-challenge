package httputil

import (
	"encoding/json"
	"errors"
	"net/http"
)

// HTTPErrorResponse
type HTTPErrorResponse struct {
	Message string `json:"message"`
	Cause   string `json:"cause"`
}

// handlerHTTPError Translates service errors into HTTP errors.
func HandlerHTTPError(w http.ResponseWriter, msg string, serviceErr error, errMapper map[error]int) {
	httpErr := HTTPErrorResponse{
		Message: msg,
		Cause:   serviceErr.Error(),
	}

	data, err := json.Marshal(httpErr)
	if err != nil {
		// TODO: Log warning here.
		return
	}

	statusCode := http.StatusInternalServerError
	for errMap, status := range errMapper {
		if errors.Is(serviceErr, errMap) {
			statusCode = status
		}
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
