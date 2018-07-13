package bitly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	Version = "0.1.0"

	defaultBaseURL   = "https://api-ssl.bitly.com"
	defaultUserAgent = "go-bitly/" + Version

	apiVersion = "v4"
)

type Client struct {
	httpClient *http.Client
	BaseURL    string
	UserAgent  string
	Debug      bool
	Groups     GroupsService
}

func NewClient(httpClient *http.Client) *Client {
	c := &Client{httpClient: httpClient, BaseURL: defaultBaseURL}
	c.UserAgent = defaultUserAgent
	c.Groups = &GroupsClient{client: c}
	return c
}

func (c *Client) NewRequest(method, path string, payload interface{}) (*http.Request, error) {
	url := c.BaseURL + path
	body := new(bytes.Buffer)
	if payload != nil {
		err := json.NewEncoder(body).Encode(payload)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", c.UserAgent)

	return req, nil
}

func (c *Client) Do(req *http.Request, obj interface{}) (*http.Response, error) {
	if c.Debug {
		log.Printf("Executing request (%v): %#v", req.URL, req)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if c.Debug {
		log.Printf("Response received: %#v", resp)
	}

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	// If obj implements the io.Writer,
	// the response body is decoded into v.
	if obj != nil {
		if w, ok := obj.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(obj)
		}
	}

	return resp, err
}

type ErrorJSON struct {
	Field     string `json:"field"`
	Message   string `json:"message"`
	ErrorCode string `json:"error_code"`
}

type ErrorResponse struct {
	response    *http.Response
	Message     string      `json:"message"`
	Errors      []ErrorJSON `json:"errors"`
	Resource    string      `json:"resource"`
	Description string      `json:"description"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %v %v",
		r.response.Request.Method, r.response.Request.URL,
		r.response.StatusCode, r.Message)
}

func CheckResponse(resp *http.Response) error {
	if code := resp.StatusCode; 200 <= code && code <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{}
	errorResponse.response = resp
	if err := json.NewDecoder(resp.Body).Decode(errorResponse); err != nil {
		return err
	}
	return errorResponse
}

func (c *Client) sendRequest(path string, payload, obj interface{}, method string) (*http.Response, error) {
	req, err := c.NewRequest(method, path, payload)
	if err != nil {
		return nil, err
	}

	return c.Do(req, obj)
}

func (c *Client) get(path string, obj interface{}) (*http.Response, error) {
	return c.sendRequest(path, nil, obj, "GET")
}

func (c *Client) post(path string, payload, obj interface{}) (*http.Response, error) {
	return c.sendRequest(path, payload, obj, "POST")
}

func (c *Client) put(path string, payload, obj interface{}) (*http.Response, error) {
	return c.sendRequest(path, payload, obj, "PUT")
}

func (c *Client) patch(path string, payload, obj interface{}) (*http.Response, error) {
	return c.sendRequest(path, payload, obj, "PATCH")
}

func (c *Client) delete(path string, payload interface{}, obj interface{}) (*http.Response, error) {
	return c.sendRequest(path, payload, obj, "DELETE")
}

func versioned(path string) string {
	return fmt.Sprintf("/%s/%s", apiVersion, strings.Trim(path, "/"))
}
