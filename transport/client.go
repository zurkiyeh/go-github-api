package transport

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
)

// This package will store all the transport layer logic

var (
	defaultBaseURLStr = "https://api.github.com/search/issues"
	defaultBaseURL, _ = url.Parse(defaultBaseURLStr)
	defaultTimeout    = 30 * time.Second
)

var (
	// ErrBadRequest is a 400 http error.
	ErrBadRequest = errors.New("bad request")
	// ErrNotFound is a 404 http error.
	ErrNotFound = errors.New("not found")
)

// A Client struct will contain logic to handle communication with Github API search endpoint.
type Client struct {
	BaseURL *url.URL

	Logger     *logrus.Logger
	HttpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		BaseURL: defaultBaseURL,
		Logger:  logrus.New(),
		HttpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}
