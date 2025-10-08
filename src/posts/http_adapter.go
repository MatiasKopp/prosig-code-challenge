package posts

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MatiasKopp/prosig-code-challenge/httputil"
	"github.com/go-chi/chi/v5"
)

var (
	errMapper = map[error]int{
		ErrBlogPostNotFound: http.StatusNotFound,
	}
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
		httputil.HandlerHTTPError(w, "unexpected error getting all blog posts", err, errMapper)
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
		httputil.HandlerHTTPError(w, msg, err, errMapper)
		return
	}

	httputil.HandlerHTTPResponse(w, http.StatusOK, post)
}

// CreatePost Creates new post.
func (a *httpAdapter) CreatePost(w http.ResponseWriter, r *http.Request) {
	var requestBody CreatePostRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		httputil.HandlerHTTPError(w, "unexpected error reading post creation body", err, errMapper)
		return
	}

	postID, err := a.Service.CreateBlogPost(requestBody.Title, requestBody.Content)
	if err != nil {
		httputil.HandlerHTTPError(w, "unexpected error creating post", err, errMapper)
		return
	}

	httputil.HandlerHTTPResponse(w, http.StatusCreated, map[string]any{"blog_post_id": postID})
}

// CreateComment Creates new comment for specific post.
func (a *httpAdapter) CreateComment(w http.ResponseWriter, r *http.Request) {
	blogPostID := chi.URLParam(r, "id")

	var requestBody CreateCommentRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		httputil.HandlerHTTPError(w, "unexpected error reading comment creation body", err, errMapper)
		return
	}

	commentID, err := a.Service.CreateComment(blogPostID, requestBody.Text)
	if err != nil {
		httputil.HandlerHTTPError(w, "unexpected error creating comment", err, errMapper)
		return
	}

	httputil.HandlerHTTPResponse(w, http.StatusCreated, map[string]any{"comment_id": commentID})
}
