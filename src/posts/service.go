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
	return s.Repository.GetBlogPost(id)
}

// CreateBlogPost Creates a new blog post and returns its generated ID.
func (s *service) CreateBlogPost(title, content string) (int64, error) {
	return s.Repository.CreateBlogPost(title, content)
}

// CreateComment Creates a new comment and associates it with a blog post.
func (s *service) CreateComment(blogPostID, text string) (int64, error) {
	_, err := s.Repository.GetBlogPost(blogPostID)
	if err != nil {
		return 0, err
	}

	return s.Repository.CreateComment(blogPostID, text)
}
