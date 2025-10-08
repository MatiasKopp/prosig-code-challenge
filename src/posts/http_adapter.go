package posts

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MatiasKopp/prosig-code-challenge/httputil"
	"github.com/go-chi/chi/v5"
)

// httpAdapter Productive post http adapter implementation
type httpAdapter struct {
	Service Service
}

// NewHTTPAdapter Returns new productive HTTP adapter implementation.
func NewHTTPAdapter(service Service) (HTTPAdapter, error) {
	return &httpAdapter{
		Service: service,
	}, nil
}

// GetAllPosts Returns all posts.
func (a *httpAdapter) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	p := httputil.GetPaginationParams(r)

	posts, err := a.Service.GetAllBlogPosts(p.Limit, p.Offset)
	if err != nil {
		httputil.HandlerHTTPError(w, "unexpected error getting all blog posts", err)
		return
	}

	if len(posts) == 0 {
		posts = []BlogPost{}
	}
	response := GetAllResponse{
		BlogPosts:  posts,
		Pagination: p,
	}

	httputil.HandlerHTTPResponse(w, http.StatusOK, response)
}

// GetPost Returns single specific post.
func (a *httpAdapter) GetPost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	post, err := a.Service.GetBlogPost(id)
	if err != nil {
		msg := fmt.Sprintf("unexpected error getting post with ID (%s)", id)
		httputil.HandlerHTTPError(w, msg, err)
		return
	}

	httputil.HandlerHTTPResponse(w, http.StatusOK, post)
}

// CreatePost Creates new post.
func (a *httpAdapter) CreatePost(w http.ResponseWriter, r *http.Request) {
	var requestBody CreatePostRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		httputil.HandlerHTTPError(w, "unexpected error reading post creation body", err)
		return
	}

	err = a.Service.CreateBlogPost(&requestBody)
	if err != nil {
		httputil.HandlerHTTPError(w, "unexpected error creating post", err)
		return
	}

	httputil.HandlerHTTPResponse(w, http.StatusCreated, nil)
}

// CreateComment Creates new comment for specific post.
func (a *httpAdapter) CreateComment(w http.ResponseWriter, r *http.Request) {
	var requestBody CreateCommentRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		httputil.HandlerHTTPError(w, "unexpected error reading comment creation body", err)
		return
	}

	err = a.Service.CreateComment(&requestBody)
	if err != nil {
		httputil.HandlerHTTPError(w, "unexpected error comment post", err)
		return
	}

	httputil.HandlerHTTPResponse(w, http.StatusCreated, nil)
}
