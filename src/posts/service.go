package posts

type service struct {
	Repository Repository
}

// NewService Returns new productive blog post service implementation.
func NewService(repository Repository) (Service, error) {
	return &service{
		Repository: repository,
	}, nil
}

// GetAllBlogPosts Returns all existing blog posts paginated.
func (s *service) GetAllBlogPosts(limit, offset int) ([]BlogPost, error) {
	return s.Repository.GetAllBlogPosts(limit, offset)
}

// GetBlogPost Returns single blog post with provided ID.
func (s *service) GetBlogPost(id string) (*BlogPost, error) {
	return nil, nil
}

// CreateBlogPost Creates new blog post.
func (s *service) CreateBlogPost(request *CreatePostRequest) error {
	return nil
}

// CreateComment Creates new comment associated with blog post.
func (s *service) CreateComment(request *CreateCommentRequest) error {
	return nil
}
