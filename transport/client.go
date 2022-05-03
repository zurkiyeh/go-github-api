package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	defaultRateLimit  = 60
)

var (
	// ErrBadRequest is a any response that isnt >200 and <299
	ErrBadRequest = errors.New("bad request")

	// ErrBadRequest is a 401 http error.
	ErrUnAuthorized = errors.New("bad request: Unauthorized")

	// ErrNotFound is a 404 http error.
	ErrNotFound = errors.New("not found")
)

// A Client struct will contain logic to handle communication with Github API search endpoint.
type Client struct {
	BaseURL *url.URL

	Logger     *logrus.Logger
	HttpClient *http.Client
	token      *string
}

func NewClient(logger *logrus.Logger, token *string) *Client {
	return &Client{
		BaseURL: defaultBaseURL,
		Logger:  logger,
		token:   token,
		HttpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// newRequest generates a http.Request based on the method
func (c *Client) NewRequest(method, path string, query string, token string) (*http.Request, error) {
	req, err := http.NewRequest(method, c.BaseURL.String()+path+query, nil)

	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if token != "" {
		c.Logger.Debug("Using Basic auth token")
		req.Header.Set("authorization", fmt.Sprintf("Bearer %s", token))
	} else {
		c.Logger.Warning(fmt.Sprintf("No auth token was set. Be aware that Github API restricts non-authenticated accounts to %d calls/hr", defaultRateLimit))
		c.Logger.Warning("Refer to README for more info on how to generate a token")
	}
	return req, nil
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
