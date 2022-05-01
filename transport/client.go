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
func (c *Client) NewRequest(method, path string, params url.Values, payload io.Reader) (*http.Request, error) {
	url := c.getURL(path, params)
	req, err := http.NewRequest(method, url.String(), payload)

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
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		c.Logger.Error("an error occured while sending request", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode <= 200 || resp.StatusCode >= 299 {
		var apiError APIError
		body, _ := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &apiError); err != nil {
			fmt.Println("Can not unmarshal JSON")
		} else {
			// 	fmt.Printf("an error occured:  %s", apiError.Message)
			// 	for _, apiErr := range apiError.Errors {
			// 		fmt.Printf(apiErr.Code)
			// 	}
			// }
			// fmt.Printf("%s\n", body)
		}
		return err
	}
	var result Response
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	fmt.Printf("Your response is %s: ", body)
	return nil
}
