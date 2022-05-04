package transport

import (
	"context"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	token := ""
	logger := logrus.New()
	c := NewClient(logger, &token)
	assert.Equal(t, c.BaseURL.Host, "api.github.com")
}

func TestRequestWithoutToken(t *testing.T) {
	token := ""
	logger := logrus.New()
	ctx := context.Background()

	c := NewClient(logger, &token)
	req, _ := c.NewRequest("GET", "", "", "")
	req = req.WithContext(ctx)
	res, _ := c.HttpClient.Do(req)

	assert.Equal(t, res.StatusCode, 200)
	assert.Equal(t, req.Header["authorization"], []string([]string(nil)))

}

func TestRequestWithBadToken(t *testing.T) {
	token := ""
	logger := logrus.New()
	ctx := context.Background()

	c := NewClient(logger, &token)
	req, _ := c.NewRequest("GET", "", "", "BAD_TOKEN")
	req = req.WithContext(ctx)
	res, _ := c.HttpClient.Do(req)

	assert.Equal(t, res.StatusCode, 401)
	assert.Equal(t, req.Header["authorization"], []string([]string(nil)))
}

// You must have a token defined in the env variable GOOD_TOKEN for this test to pass
func TestRequestWithGoodToken(t *testing.T) {
	token := ""
	logger := logrus.New()
	ctx := context.Background()

	c := NewClient(logger, &token)
	req, _ := c.NewRequest("GET", "", "", os.Getenv("GOOD_TOKEN"))
	req = req.WithContext(ctx)
	res, _ := c.HttpClient.Do(req)
	assert.Equal(t, res.StatusCode, 200)
	assert.NotEqual(t, req.Header["authorization"], nil)
}
