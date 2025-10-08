package posts

// BlogPost Represents blogpost data.
type BlogPost struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Comments []Comment `json:"comments"`
}

// Comment Blogpost comment.
type Comment struct {
	ID          string `json:"id"`
	CommentText string `json:"comment_text"`
}
