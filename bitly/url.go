package bitly

import "net/url"

func buildURL(rawURL string, params url.Values) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	u.RawQuery = params.Encode()
	return u.String(), nil
}
