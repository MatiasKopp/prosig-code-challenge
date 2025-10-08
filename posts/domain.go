package posts

// BlogPost Represents blogpost data.
type BlogPost struct {
	ID       string
	Title    string
	Content  string
	Comments []Comment
}

// Comment Blogpost comment.
type Comment struct {
	Text string
}
