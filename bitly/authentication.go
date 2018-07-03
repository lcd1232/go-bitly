package bitly

import "encoding/base64"

const (
	httpHeaderAuthorization = "Authorization"
)

// Provides credentials that can be used for authenticating with Bitly.
//
// See http://dev.bitly.com/v4/#section/Authentication
type Credentials interface {
	Headers() map[string]string
}

type httpBasicCredentials struct {
	username string
	password string
}

func NewHTTPBasicCredentials(username, password string) *httpBasicCredentials {
	return &httpBasicCredentials{
		username: username,
		password: password,
	}
}

func (c *httpBasicCredentials) Headers() map[string]string {
	return map[string]string{httpHeaderAuthorization: "Basic " + basicAuth(c.username, c.password)}
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

type oauthTokenCredentials struct {
	oauthToken string
}

func NewOauthTokenCredentials(token string) *oauthTokenCredentials {
	return &oauthTokenCredentials{
		oauthToken: token,
	}
}

func (c *oauthTokenCredentials) Headers() map[string]string {
	return map[string]string{httpHeaderAuthorization: "Bearer " + c.oauthToken}
}

//type accessTokenCredentials struct {
//	basicCredentials *httpBasicCredentials
//	accessToken string
//}
//
//func NewAccessTokenCredentials(username, password string) *accessTokenCredentials {
//	return &accessTokenCredentials{
//		basicCredentials: NewHTTPBasicCredentials(username, password),
//	}
//}
//
//func (c *accessTokenCredentials) Headers() map[string]string {
//	if c.accessToken == "" {
//		//Logic how to get access_token
//	}
//	return map[string]string{}
//}
