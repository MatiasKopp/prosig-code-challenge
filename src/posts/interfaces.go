package posts

import (
	"net/http"

	"github.com/MatiasKopp/prosig-code-challenge/httputil"
)

// HTTPAdapter Posts http adapter interface.
type HTTPAdapter interface {
	// GetAllPosts Returns all posts.
	GetAllPosts(http.ResponseWriter, *http.Request)
	// GetPost Returns single specific post.
	GetPost(http.ResponseWriter, *http.Request)
	// CreatePost Creates new post.
	CreatePost(http.ResponseWriter, *http.Request)
	// CreateComment Creates new comment for specific post.
	CreateComment(http.ResponseWriter, *http.Request)
}

// Service Posts services interface.
type Service interface {
	// GetAllBlogPosts Returns all existing blog posts paginated.
	GetAllBlogPosts(page, limit int) ([]BlogPost, error)
	// GetBlogPost Returns single blog post with provided ID.
	GetBlogPost(id string) (*BlogPost, error)
	// CreateBlogPost Creates new blog post.
	CreateBlogPost(request *CreatePostRequest) error
	// CreateComment Creates new comment associated with blog post.
	CreateComment(request *CreateCommentRequest) error
}

// Repository Posts repository interface.
type Repository interface {
	// GetAllBlogPosts Returns all existing blog posts paginated.
	GetAllBlogPosts(page, limit int) ([]BlogPost, error)
	// GetBlogPost Returns single blog post with provided ID.
	GetBlogPost(id string) (*BlogPost, error)
	// CreateBlogPost Creates new blog post.
	CreateBlogPost(request *CreatePostRequest) error
	// CreateComment Creates new comment associated with blog post.
	CreateComment(request *CreateCommentRequest) error
}

// CreatePostRequest Structure used in new post request.
type CreatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// CreateCommentRequest Structure used in new comment request.
type CreateCommentRequest struct {
	BlogPostID string `json:"blog_post_id"`
	Text       string `json:"text"`
}

// GetAllResponse Get all blog posts response
type GetAllResponse struct {
	BlogPosts  []BlogPost          `json:"blog_posts"`
	Pagination httputil.Pagination `json:"pagination"`
}
