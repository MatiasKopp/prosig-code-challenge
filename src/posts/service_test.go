package posts

import (
	"errors"
	"reflect"
	"testing"
)

func Test_service_GetAllBlogPosts(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(m *MocksRepository)
		limit   int
		offset  int
		want    []BlogPost
		wantErr bool
	}{
		{
			name: "success",
			setup: func(m *MocksRepository) {
				m.EXPECT().GetAllBlogPosts(10, 0).Return([]BlogPost{{ID: "1", Title: "A", Content: "B"}}, nil)
			},
			limit:   10,
			offset:  0,
			want:    []BlogPost{{ID: "1", Title: "A", Content: "B"}},
			wantErr: false,
		},
		{
			name: "repo error",
			setup: func(m *MocksRepository) {
				m.EXPECT().GetAllBlogPosts(10, 0).Return(nil, errors.New("fail"))
			},
			limit:   10,
			offset:  0,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMocksRepository(t)
			if tt.setup != nil {
				tt.setup(repo)
			}
			s := &service{Repository: repo}
			got, err := s.GetAllBlogPosts(tt.limit, tt.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllBlogPosts() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllBlogPosts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_GetBlogPost(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(m *MocksRepository)
		id      string
		want    *BlogPost
		wantErr bool
	}{
		{
			name: "success",
			setup: func(m *MocksRepository) {
				m.EXPECT().GetBlogPost("1").Return(&BlogPost{ID: "1", Title: "T", Content: "C"}, nil)
			},
			id:      "1",
			want:    &BlogPost{ID: "1", Title: "T", Content: "C"},
			wantErr: false,
		},
		{
			name: "repo error",
			setup: func(m *MocksRepository) {
				m.EXPECT().GetBlogPost("1").Return(nil, errors.New("fail"))
			},
			id:      "1",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMocksRepository(t)
			if tt.setup != nil {
				tt.setup(repo)
			}
			s := &service{Repository: repo}
			got, err := s.GetBlogPost(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBlogPost() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBlogPost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_CreateBlogPost(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(m *MocksRepository)
		title   string
		content string
		want    int64
		wantErr bool
	}{
		{
			name: "success",
			setup: func(m *MocksRepository) {
				m.EXPECT().CreateBlogPost("T", "C").Return(int64(42), nil)
			},
			title:   "T",
			content: "C",
			want:    42,
			wantErr: false,
		},
		{
			name: "repo error",
			setup: func(m *MocksRepository) {
				m.EXPECT().CreateBlogPost("T", "C").Return(int64(0), errors.New("fail"))
			},
			title:   "T",
			content: "C",
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMocksRepository(t)
			if tt.setup != nil {
				tt.setup(repo)
			}
			s := &service{Repository: repo}
			got, err := s.CreateBlogPost(tt.title, tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateBlogPost() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("CreateBlogPost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_CreateComment(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(m *MocksRepository)
		blogPostID string
		text       string
		want       int64
		wantErr    bool
	}{
		{
			name: "success",
			setup: func(m *MocksRepository) {
				m.EXPECT().GetBlogPost("1").Return(&BlogPost{ID: "1"}, nil)
				m.EXPECT().CreateComment("1", "comment").Return(int64(99), nil)
			},
			blogPostID: "1",
			text:       "comment",
			want:       99,
			wantErr:    false,
		},
		{
			name: "get post error",
			setup: func(m *MocksRepository) {
				m.EXPECT().GetBlogPost("1").Return(nil, errors.New("not found"))
			},
			blogPostID: "1",
			text:       "comment",
			want:       0,
			wantErr:    true,
		},
		{
			name: "create comment error",
			setup: func(m *MocksRepository) {
				m.EXPECT().GetBlogPost("1").Return(&BlogPost{ID: "1"}, nil)
				m.EXPECT().CreateComment("1", "comment").Return(int64(0), errors.New("fail"))
			},
			blogPostID: "1",
			text:       "comment",
			want:       0,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMocksRepository(t)
			if tt.setup != nil {
				tt.setup(repo)
			}
			s := &service{Repository: repo}
			got, err := s.CreateComment(tt.blogPostID, tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateComment() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("CreateComment() = %v, want %v", got, tt.want)
			}
		})
	}
}
