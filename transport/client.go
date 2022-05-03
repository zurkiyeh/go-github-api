package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
)

// This package will store all the transport layer logic

var (
	defaultBaseURLStr = "https://api.github.com/"
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

func NewClient(logger *logrus.Logger) *Client {
	return &Client{
		BaseURL: defaultBaseURL,
		Logger:  logger,
		HttpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// newRequest generates a http.Request based on the method
func (c *Client) NewRequest(method, path string, query string) (*http.Request, error) {
	req, err := http.NewRequest(method, c.BaseURL.String()+path+query, nil)

	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (c *Client) getURL(path string, params url.Values) *url.URL {
	return &url.URL{
		Scheme:   c.BaseURL.Scheme,
		Host:     c.BaseURL.Host,
		Path:     fmt.Sprintf("/%s", path),
		RawQuery: "q=" + params.Encode(),
	}
}

// do performs a roundtrip using the underlying client
func (c *Client) Do(ctx context.Context, req *http.Request) error {
	if ctx == nil {
		return errors.New("context must be non-nil")
	}
	req = req.WithContext(ctx)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		c.Logger.Errorf("an error occured while sending request: %s", err)
		return err
	}
	defer resp.Body.Close()

	if ok := validateResponse(resp); ok != nil {
		return ok

	}

	var result Response
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		c.Logger.Error("Can not unmarshal JSON")
	}

	c.Logger.Debug("Request has returned ", len(result.Items), " items")
	for i, pr := range result.Items {
		c.Logger.Debug(" Pull-request #", i, " has title: ", pr.Title)
	}
	return nil
}
