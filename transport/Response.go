package transport

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

// Checks for status errors
func validateResponse(resp *http.Response, logger *logrus.Logger) error {

	switch status := resp.StatusCode; {
	case status == 401:
		var apiError APIError
		body, _ := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &apiError); err != nil {
			logger.Debug("error while parsing API error response: ", err)
			return fmt.Errorf("can not unmarshal reponse JSON")
		}
		return fmt.Errorf("error validating response: %s", ErrUnAuthorized)

	case (status > 200 && status < 299):
		var apiError APIError
		body, _ := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &apiError); err != nil {
			logger.Debug("error while parsing API error response: ", err)
			return fmt.Errorf("can not unmarshal reponse JSON")
		}
		logger.Debug("error while parsing API error response: ", apiError.Errors, apiError.Message)
		return fmt.Errorf("an error ocurred while sending request: %s ", ErrBadRequest)
	default:
		return nil
	}
}
