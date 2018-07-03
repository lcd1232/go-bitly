package bitly

import (
	"reflect"
	"testing"
)

func TestBasicAuth(t *testing.T) {
	testcase := []struct {
		username string
		password string
		want     string
	}{
		{
			username: "test",
			password: "test",
			want:     "dGVzdDp0ZXN0",
		},
		{
			username: "admin",
			password: "super_password",
			want:     "YWRtaW46c3VwZXJfcGFzc3dvcmQ=",
		},
	}
	for _, tc := range testcase {
		got := basicAuth(tc.username, tc.password)
		if got != tc.want {
			t.Fatalf("got %v want %v", got, tc.want)
		}
	}
}

func TestHttpBasicCredentials_Headers(t *testing.T) {
	testcase := []struct {
		username string
		password string
		want     map[string]string
	}{
		{
			username: "test",
			password: "test",
			want:     map[string]string{httpHeaderAuthorization: "Basic dGVzdDp0ZXN0"},
		},
		{
			username: "admin",
			password: "super_password",
			want:     map[string]string{httpHeaderAuthorization: "Basic YWRtaW46c3VwZXJfcGFzc3dvcmQ="},
		},
	}

	for _, tc := range testcase {
		credentials := NewHTTPBasicCredentials(tc.username, tc.password)
		got := credentials.Headers()
		if !reflect.DeepEqual(got, tc.want) {
			t.Fatalf("got %v want %v", got, tc.want)
		}
	}
}

func TestOauthTokenCredentials_Headers(t *testing.T) {
	testcase := []struct {
		token string
		want  map[string]string
	}{
		{
			token: "test",
			want:  map[string]string{httpHeaderAuthorization: "Bearer test"},
		},
		{
			token: "secret_token",
			want:  map[string]string{httpHeaderAuthorization: "Bearer secret_token"},
		},
	}
	for _, tc := range testcase {
		credentials := NewOauthTokenCredentials(tc.token)
		got := credentials.Headers()
		if !reflect.DeepEqual(got, tc.want) {
			t.Fatalf("got %v want %v", got, tc.want)
		}
	}
}
