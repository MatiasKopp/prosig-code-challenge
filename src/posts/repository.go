package posts

import (
	"database/sql"
	"maps"
	"slices"
)

// repository Simple productive repository pointing to sqlite db.
type repository struct {
	db *sql.DB
}

// NewRepository Returns new productive repository implementation.
func NewRepository(db *sql.DB) (Repository, error) {
	return &repository{db: db}, nil
}

// blogPostComment Internal struct to flatten blogpost comment relationship.
type blogPostComment struct {
	BlogPostID      string
	BlogPostTitle   string
	BlogPostContent string
	CommentID       string
	CommentText     string
}

// GetAllBlogPosts Returns all existing blog posts paginated.
func (r *repository) GetAllBlogPosts(limit, offset int) ([]BlogPost, error) {
	rows, err := r.db.Query(`
		SELECT 
			a.id, 
			a.title,
			a.content,
			c.id,
			c.comment_text

		FROM blog_posts a
			LEFT JOIN blog_posts_comments b
				ON a.id = b.blog_post_id
			INNER JOIN comments c
				ON b.comment_id = c.id

		ORDER BY a.id
		LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogPostComments []blogPostComment
	for rows.Next() {
		var i blogPostComment
		if err := rows.Scan(&i.BlogPostID, &i.BlogPostTitle, &i.BlogPostContent, &i.CommentID, &i.CommentText); err != nil {
			return nil, err
		}
		blogPostComments = append(blogPostComments, i)
	}

	res := map[string]BlogPost{}
	for _, item := range blogPostComments {
		blogPost, exists := res[item.BlogPostID]
		if !exists {
			blogPost = BlogPost{
				ID:      item.BlogPostID,
				Title:   item.BlogPostTitle,
				Content: item.BlogPostContent,
			}
		}

		if item.CommentID != "" {
			blogPost.Comments = append(blogPost.Comments, Comment{
				ID:          item.CommentID,
				CommentText: item.CommentText,
			})
		}

		res[item.BlogPostID] = blogPost
	}

	return slices.Collect(maps.Values(res)), nil
}

// GetBlogPost Returns single blog post with provided ID.
func (r *repository) GetBlogPost(id string) (*BlogPost, error) {
	return nil, nil
}

// CreateBlogPost Creates new blog post.
func (r *repository) CreateBlogPost(request *CreatePostRequest) error {
	return nil
}

// CreateComment Creates new comment associated with blog post.
func (r *repository) CreateComment(request *CreateCommentRequest) error {
	return nil
}
