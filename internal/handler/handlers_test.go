package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var handlers = MockHandlers()

func TestShowPosts(t *testing.T) {
	handlers.App.Repository.Posts.Insert(1, "title", "content")

	tests := []*struct {
		name     string
		path     string
		method   string
		wantCode int
	}{
		{
			path:     "/invalid-path",
			method:   "GET",
			wantCode: http.StatusNotFound,
		},
		{
			path:     "/post",
			method:   "POST",
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			path:     "/post",
			method:   "GET",
			wantCode: http.StatusBadRequest,
		},
		{
			path:     "/post?id=invalid-value",
			method:   "GET",
			wantCode: http.StatusBadRequest,
		},
		{
			path:     "/post?id=0001",
			method:   "GET",
			wantCode: http.StatusBadRequest,
		},
		{
			path:     "/post?id=-1",
			method:   "GET",
			wantCode: http.StatusNotFound,
		},
		{
			path:     "/post?id=1",
			method:   "GET",
			wantCode: http.StatusOK,
		},
	}

	for i, tt := range tests {
		tt.name = fmt.Sprintf("test%d", i+1)

		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			http.HandlerFunc(handlers.ShowPost).ServeHTTP(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tt.wantCode)
			}
		})
	}
}
