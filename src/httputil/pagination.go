package httputil

import (
	"net/http"
	"strconv"
)

// Pagination Basic pagination structure.
type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Page   int `json:"page"`
}

// GetPaginationParams Parses pagination params from HTTP request.
func GetPaginationParams(r *http.Request) Pagination {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	return Pagination{Limit: limit, Offset: offset, Page: page}
}
