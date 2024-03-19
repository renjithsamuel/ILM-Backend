package googlebooks

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"integrated-library-service/model"
)

// Client ...
type Client interface {
	GetGoogleBooks(request *model.GetAllNewBooksRequest) ([]model.Book, int, error)
}

// GoogleBooksClient
type GoogleBooksClient struct {
	url    string
	apiKey string
	client *http.Client 
}

// New returns new instance of NewGoogleService
func NewGoogleService(URL string, apiKey string, client *http.Client) *GoogleBooksClient {
	return &GoogleBooksClient{client: client, url: URL, apiKey: apiKey}
}

// GetClient returns new generated http.Client for Google.client
func GetClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   timeout,
	}
}

// parseURL returns prefixed URL with base google URL
func (googleBooksClient *GoogleBooksClient) parseURL(URL string) (*url.URL, error) {
	fullURL, err := url.Parse(fmt.Sprint(googleBooksClient.url, URL))
	if err != nil {
		return nil, err
	}

	return fullURL, nil
}

// do modify request by attaching baseURL and basic Authentication keys
func (googleBooksClient *GoogleBooksClient) do(request *http.Request) (*WebResponse, error) {
	fullURL, err := googleBooksClient.parseURL(request.URL.String())
	if err != nil {
		return nil, err
	}

	request.URL = fullURL
	if request.Body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	res, err := googleBooksClient.client.Do(request)
	if err != nil {
		return nil, err
	}

	return NewWebResponse(res), nil
}
