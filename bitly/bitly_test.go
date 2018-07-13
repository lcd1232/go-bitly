package bitly

import (
	"net/http"
	"net/http/httptest"
)

var (
	mux    *http.ServeMux
	client *Client
	server *httptest.Server
)

func setupMockServer() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewClient(NewOauthTokenCredentials("bitly-token"))
	client.BaseURL = server.URL
}

func teardownMockServer() {
	server.Close()
}
