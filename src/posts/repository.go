package posts

import (
	"database/sql"
	"errors"
	"fmt"
	"maps"
	"slices"
)

var (
	// ErrBlogPostNotFound Blog post not found error.
	ErrBlogPostNotFound = errors.New("blog post not found")
)

// repository Simple productive repository pointing to sqlite db.
type repository struct {
	db *sql.DB
}

// NewRepository Returns new productive repository implementation.
func NewRepository(db *sql.DB) (Repository, error) {
	return &repository{db: db}, nil
}

// blogPostComment Internal struct to flatten blogpost-comment relationship.
type blogPostComment struct {
	BlogPostID      string
	BlogPostTitle   string
	BlogPostContent string
	CommentID       sql.NullString
	CommentText     sql.NullString
}

// readBlogPosts Internal reusable function that retrieves blog posts and comments.
// If `id` is non-empty, it fetches a single post. If not, it fetches all (optionally paginated).
func (r *repository) readBlogPosts(id string, limit, offset int) ([]BlogPost, error) {
	query := `
		SELECT 
			a.id, 
			a.title,
			a.content,
			c.id,
			c.comment_text
		FROM blog_posts a
			LEFT JOIN blog_posts_comments b
				ON a.id = b.blog_post_id
			LEFT JOIN comments c
				ON b.comment_id = c.id
	`
	args := []any{}

	if id != "" {
		query += " WHERE a.id = ?"
		args = append(args, id)
	}

	query += " ORDER BY a.id"

	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}
	if offset > 0 {
		query += " OFFSET ?"
		args = append(args, offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query blog posts: %w", err)
	}
	defer rows.Close()

	var blogPostComments []blogPostComment
	for rows.Next() {
		var i blogPostComment
		if err := rows.Scan(
			&i.BlogPostID,
			&i.BlogPostTitle,
			&i.BlogPostContent,
			&i.CommentID,
			&i.CommentText,
		); err != nil {
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

		if item.CommentID.Valid {
			blogPost.Comments = append(blogPost.Comments, Comment{
				ID:          item.CommentID.String,
				CommentText: item.CommentText.String,
			})
		}

		res[item.BlogPostID] = blogPost
	}

	return slices.Collect(maps.Values(res)), nil
}

// GetAllBlogPosts Returns all existing blog posts paginated.
func (r *repository) GetAllBlogPosts(limit, offset int) ([]BlogPost, error) {
	return r.readBlogPosts("", limit, offset)
}

// GetBlogPost Returns a single blog post with its comments.
func (r *repository) GetBlogPost(id string) (*BlogPost, error) {
	posts, err := r.readBlogPosts(id, 0, 0)
	if err != nil {
		return nil, err
	}
	if len(posts) == 0 {
		return nil, ErrBlogPostNotFound
	}
	return &posts[0], nil
}

// CreateBlogPost Creates a new blog post and returns its generated ID.
func (r *repository) CreateBlogPost(title, content string) (int64, error) {
	res, err := r.db.Exec(`
		INSERT INTO blog_posts (title, content)
		VALUES (?, ?)`,
		title, content)
	if err != nil {
		return 0, fmt.Errorf("failed to create blog post: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get inserted blog post ID: %w", err)
	}

	return id, nil
}

// CreateComment Creates a new comment and associates it with a blog post.
func (r *repository) CreateComment(blogPostID, text string) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to start tx: %w", err)
	}
	defer tx.Rollback()

	// Insert comment
	res, err := tx.Exec(`
		INSERT INTO comments (comment_text)
		VALUES (?)`,
		text)
	if err != nil {
		return 0, fmt.Errorf("failed to insert comment: %w", err)
	}

	commentID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get comment ID: %w", err)
	}

	// Associate with blog post
	_, err = tx.Exec(`
		INSERT INTO blog_posts_comments (blog_post_id, comment_id)
		VALUES (?, ?)`,
		blogPostID, commentID)
	if err != nil {
		return 0, fmt.Errorf("failed to link comment: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit tx: %w", err)
	}

	return commentID, nil
}
