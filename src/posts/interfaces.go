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
	// CreateBlogPost Creates a new blog post and returns its generated ID.
	CreateBlogPost(title, content string) (int64, error)
	// CreateComment Creates a new comment and associates it with a blog post.
	CreateComment(blogPostID, text string) (int64, error)
}

// Repository Posts repository interface.
type Repository interface {
	// GetAllBlogPosts Returns all existing blog posts paginated.
	GetAllBlogPosts(page, limit int) ([]BlogPost, error)
	// GetBlogPost Returns single blog post with provided ID.
	GetBlogPost(id string) (*BlogPost, error)
	// CreateBlogPost Creates a new blog post and returns its generated ID.
	CreateBlogPost(title, content string) (int64, error)
	// CreateComment Creates a new comment and associates it with a blog post.
	CreateComment(blogPostID, text string) (int64, error)
}

// CreatePostRequest Structure used in new post request.
type CreatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// CreateCommentRequest Structure used in new comment request.
type CreateCommentRequest struct {
	Text string `json:"text"`
}

// GetAllResponse Get all blog posts response
type GetAllResponse struct {
	BlogPosts  []BlogPost          `json:"blog_posts"`
	Pagination httputil.Pagination `json:"pagination"`
}
