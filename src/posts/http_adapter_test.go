package posts

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_httpAdapter_GetAllPosts(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() *httpAdapter
		request    *http.Request
		wantStatus int
		wantBody   string
	}{
		{
			name: "success_200_empty_list",
			setup: func() *httpAdapter {
				service := NewMocksService(t)

				service.EXPECT().GetAllBlogPosts(10, 0).Return([]BlogPost{}, nil)

				return &httpAdapter{
					Service: service,
				}
			},
			request:    httptest.NewRequest(http.MethodGet, "/posts", nil),
			wantStatus: http.StatusOK,
			wantBody:   "{\"blog_posts\":[],\"pagination\":{\"limit\":10,\"offset\":0,\"page\":1}}",
		},
		{
			name: "success_200_empty_list",
			setup: func() *httpAdapter {
				service := NewMocksService(t)

				service.EXPECT().GetAllBlogPosts(10, 0).Return([]BlogPost{
					{
						ID:       "1",
						Title:    "First Post",
						Content:  "This is the body of the first post",
						Comments: nil,
					},
				}, nil)

				return &httpAdapter{
					Service: service,
				}
			},
			request:    httptest.NewRequest(http.MethodGet, "/posts", nil),
			wantStatus: http.StatusOK,
			wantBody:   "{\"blog_posts\":[{\"id\":\"1\",\"title\":\"First Post\",\"content\":\"This is the body of the first post\",\"comments\":null}],\"pagination\":{\"limit\":10,\"offset\":0,\"page\":1}}",
		},
		{
			name: "service_error_500",
			setup: func() *httpAdapter {
				service := NewMocksService(t)

				service.EXPECT().GetAllBlogPosts(10, 0).Return(nil, errors.New("internal error"))

				return &httpAdapter{
					Service: service,
				}
			},
			request:    httptest.NewRequest(http.MethodGet, "/posts", nil),
			wantStatus: http.StatusInternalServerError,
			wantBody:   "{\"message\":\"unexpected error getting all blog posts\",\"cause\":\"internal error\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.setup()
			recorder := httptest.NewRecorder()
			a.GetAllPosts(recorder, tt.request)

			if recorder.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", recorder.Code, tt.wantStatus)
			}
			if tt.wantBody != "" {
				body, _ := io.ReadAll(recorder.Body)
				if string(body) != tt.wantBody {
					t.Errorf("got body %q, want %q", string(body), tt.wantBody)
				}
			}
		})
	}
}

func Test_httpAdapter_GetPost(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() *httpAdapter
		request    *http.Request
		wantStatus int
		wantBody   string
		id         string
	}{
		{
			name: "success_200",
			setup: func() *httpAdapter {
				service := NewMocksService(t)

				service.EXPECT().GetBlogPost("1").Return(&BlogPost{
					ID:       "1",
					Title:    "First Post",
					Content:  "This is the body of the first post",
					Comments: nil,
				}, nil)

				return &httpAdapter{
					Service: service,
				}
			},
			request:    httptest.NewRequest(http.MethodGet, "/posts/1", nil),
			wantStatus: http.StatusOK,
			wantBody:   "{\"id\":\"1\",\"title\":\"First Post\",\"content\":\"This is the body of the first post\",\"comments\":[]}",
			id:         "1",
		},
		{
			name: "service_error_500",
			setup: func() *httpAdapter {
				service := NewMocksService(t)

				service.EXPECT().GetBlogPost("1").Return(nil, errors.New("internal error"))

				return &httpAdapter{
					Service: service,
				}
			},
			request:    httptest.NewRequest(http.MethodGet, "/posts/1", nil),
			wantStatus: http.StatusInternalServerError,
			wantBody:   "{\"message\":\"unexpected error getting post with ID (1)\",\"cause\":\"internal error\"}",
			id:         "1",
		},
		{
			name: "service_not_found_404",
			setup: func() *httpAdapter {
				service := NewMocksService(t)

				service.EXPECT().GetBlogPost("1").Return(nil, ErrBlogPostNotFound)

				return &httpAdapter{
					Service: service,
				}
			},
			request:    httptest.NewRequest(http.MethodGet, "/posts/1", nil),
			wantStatus: http.StatusNotFound,
			wantBody:   "{\"message\":\"unexpected error getting post with ID (1)\",\"cause\":\"blog post not found\"}",
			id:         "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.setup()
			recorder := httptest.NewRecorder()

			// Inject chi route context with id param
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.id)
			ctx := context.WithValue(tt.request.Context(), chi.RouteCtxKey, rctx)
			reqWithCtx := tt.request.WithContext(ctx)

			a.GetPost(recorder, reqWithCtx)

			if recorder.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", recorder.Code, tt.wantStatus)
			}
			if tt.wantBody != "" {
				body, _ := io.ReadAll(recorder.Body)
				if string(body) != tt.wantBody {
					t.Errorf("got body %q, want %q", string(body), tt.wantBody)
				}
			}
		})
	}
}

func Test_httpAdapter_CreatePost(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() *httpAdapter
		request    *http.Request
		wantStatus int
		wantBody   string
	}{
		{
			name: "success_201",
			setup: func() *httpAdapter {
				service := NewMocksService(t)
				service.EXPECT().CreateBlogPost("some_title", "some_content").Return(1, nil)
				return &httpAdapter{Service: service}
			},
			request:    httptest.NewRequest(http.MethodPost, "/posts", io.NopCloser(strings.NewReader(`{"title":"some_title","content":"some_content"}`))),
			wantStatus: http.StatusCreated,
			wantBody:   "{\"blog_post_id\":1}",
		},
		{
			name: "validation_error_400",
			setup: func() *httpAdapter {
				service := NewMocksService(t)
				return &httpAdapter{Service: service}
			},
			request:    httptest.NewRequest(http.MethodPost, "/posts", io.NopCloser(strings.NewReader(`{"title":"","content":""}`))),
			wantStatus: http.StatusBadRequest,
			wantBody:   "{\"message\":\"missing title or content\",\"cause\":\"bad request\"}",
		},
		{
			name: "service_error_500",
			setup: func() *httpAdapter {
				service := NewMocksService(t)
				service.EXPECT().CreateBlogPost("some_title", "some_content").Return(0, errors.New("internal error"))
				return &httpAdapter{Service: service}
			},
			request:    httptest.NewRequest(http.MethodPost, "/posts", io.NopCloser(strings.NewReader(`{"title":"some_title","content":"some_content"}`))),
			wantStatus: http.StatusInternalServerError,
			wantBody:   "{\"message\":\"unexpected error creating post\",\"cause\":\"internal error\"}",
		},
		{
			name: "json_unmarshal_error",
			setup: func() *httpAdapter {
				service := NewMocksService(t)
				return &httpAdapter{Service: service}
			},
			request:    httptest.NewRequest(http.MethodPost, "/posts", io.NopCloser(strings.NewReader(`{"title":"bad json"`))), // missing closing brace
			wantStatus: http.StatusInternalServerError,
			wantBody:   "{\"message\":\"unexpected error reading post creation body\",\"cause\":\"unexpected EOF\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.setup()
			recorder := httptest.NewRecorder()
			a.CreatePost(recorder, tt.request)

			if recorder.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", recorder.Code, tt.wantStatus)
			}
			if tt.wantBody != "" {
				body, _ := io.ReadAll(recorder.Body)
				if string(body) != tt.wantBody {
					t.Errorf("got body %q, want %q", string(body), tt.wantBody)
				}
			}
		})
	}
}

func Test_httpAdapter_CreateComment(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() *httpAdapter
		request    *http.Request
		wantStatus int
		wantBody   string
		id         string
	}{
		{
			name: "success_201",
			setup: func() *httpAdapter {
				service := NewMocksService(t)
				service.EXPECT().CreateComment("1", "some comment").Return(2, nil)
				return &httpAdapter{Service: service}
			},
			request:    httptest.NewRequest(http.MethodPost, "/posts/1/comments", io.NopCloser(strings.NewReader(`{"text":"some comment"}`))),
			wantStatus: http.StatusCreated,
			wantBody:   "{\"comment_id\":2}",
			id:         "1",
		},
		{
			name: "validation_error_400",
			setup: func() *httpAdapter {
				service := NewMocksService(t)
				return &httpAdapter{Service: service}
			},
			request:    httptest.NewRequest(http.MethodPost, "/posts/1/comments", io.NopCloser(strings.NewReader(`{"text":""}`))),
			wantStatus: http.StatusBadRequest,
			wantBody:   "{\"message\":\"missing comment text\",\"cause\":\"bad request\"}",
			id:         "1",
		},
		{
			name: "service_error_500",
			setup: func() *httpAdapter {
				service := NewMocksService(t)
				service.EXPECT().CreateComment("1", "some comment").Return(0, errors.New("internal error"))
				return &httpAdapter{Service: service}
			},
			request:    httptest.NewRequest(http.MethodPost, "/posts/1/comments", io.NopCloser(strings.NewReader(`{"text":"some comment"}`))),
			wantStatus: http.StatusInternalServerError,
			wantBody:   "{\"message\":\"unexpected error creating comment\",\"cause\":\"internal error\"}",
			id:         "1",
		},
		{
			name: "json_unmarshal_error",
			setup: func() *httpAdapter {
				service := NewMocksService(t)
				return &httpAdapter{Service: service}
			},
			request:    httptest.NewRequest(http.MethodPost, "/posts/1/comments", io.NopCloser(strings.NewReader(`{"text":"bad json"`))), // missing closing brace
			wantStatus: http.StatusInternalServerError,
			wantBody:   "{\"message\":\"unexpected error reading comment creation body\",\"cause\":\"unexpected EOF\"}",
			id:         "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.setup()
			recorder := httptest.NewRecorder()

			// Inject chi route context with id param
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.id)
			ctx := context.WithValue(tt.request.Context(), chi.RouteCtxKey, rctx)
			reqWithCtx := tt.request.WithContext(ctx)

			a.CreateComment(recorder, reqWithCtx)

			if recorder.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", recorder.Code, tt.wantStatus)
			}
			if tt.wantBody != "" {
				body, _ := io.ReadAll(recorder.Body)
				if string(body) != tt.wantBody {
					t.Errorf("got body %q, want %q", string(body), tt.wantBody)
				}
			}
		})
	}
}
