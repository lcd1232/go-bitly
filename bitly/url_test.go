package bitly

import (
	"net/url"
	"testing"
)

func TestBuildURL(t *testing.T) {
	testCases := []struct {
		desc       string
		url        string
		params     url.Values
		wantResult string
		wantError  string
	}{
		{
			desc:       "url as top domain",
			url:        "http://example.com",
			params:     url.Values{"a": []string{"b"}},
			wantResult: "http://example.com?a=b",
			wantError:  "",
		},
		{
			desc:       "params at url",
			url:        "http://example.com/path/to/random/place",
			params:     url.Values{"force": []string{"true"}},
			wantResult: "http://example.com/path/to/random/place?force=true",
			wantError:  "",
		},
		{
			desc:       "multiple params",
			url:        "http://example.com/path/to/random/place",
			params:     url.Values{"force": []string{"true"}, "sort": []string{"desc"}},
			wantResult: "http://example.com/path/to/random/place?force=true&sort=desc",
			wantError:  "",
		},
		{
			desc:       "empty params",
			url:        "http://example.com/path",
			params:     url.Values{},
			wantResult: "http://example.com/path",
			wantError:  "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := buildURL(tc.url, tc.params)
			if tc.wantError == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantError != "" && (err == nil || tc.wantError != err.Error()) {
				t.Fatalf("want error %v got %v", tc.wantError, err)
			}
			if got != tc.wantResult {
				t.Fatalf("got %v want %v", got, tc.wantResult)
			}
		})
	}
}
