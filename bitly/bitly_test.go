package bitly

import (
	"net/http"
	"testing"
)

func TestClient_NewRequest(t *testing.T) {
	testCases := []struct {
		baseURL   string
		url       string
		wantURL   string
		wantError string
	}{
		{
			baseURL:   "http://example.com",
			url:       "/foo",
			wantURL:   "http://example.com/foo",
			wantError: "",
		},
		{
			baseURL:   "http://example.com/bar",
			url:       "/foo",
			wantURL:   "http://example.com/bar/foo",
			wantError: "",
		},
	}
	for _, tc := range testCases {
		c := NewClient(http.DefaultClient)
		c.BaseURL = tc.baseURL
		req, err := c.NewRequest("GET", tc.url, nil)
		if tc.wantError == "" && err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tc.wantError != "" && (err == nil || tc.wantError != err.Error()) {
			t.Fatalf("want error %v got %v", tc.wantError, err)
		}
		if rawURL := req.URL.String(); rawURL != tc.wantURL {
			t.Fatalf("got url %v want %v", rawURL, tc.wantURL)
		}
	}
}
