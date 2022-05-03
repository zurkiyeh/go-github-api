package transport

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func validateResponse(resp *http.Response) error {

	switch status := resp.StatusCode; {
	case status == 401:
		return fmt.Errorf("error validating response: %s", ErrUnAuthorized)

	case (status > 200 && status < 299):
		var apiError APIError
		body, _ := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &apiError); err != nil {
			return fmt.Errorf("can not unmarshal reponse JSON")
		}
		return fmt.Errorf("an error ocurred while sending request: %s ", ErrBadRequest)
	default:
		return nil
	}
}
