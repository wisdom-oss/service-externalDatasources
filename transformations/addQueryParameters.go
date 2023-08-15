package transformations

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// AddQueryParameters is a type alias on a map between a string and an array of
// string values. The type alias also implements the TransformationFunction
// interface.
//
// This transformation is only applicable before a request is sent to the
// remote host.
type AddQueryParameters map[string][]string

func (parameterMapping AddQueryParameters) ApplyBefore(r *http.Request, data []byte) error {
	// try to parse the configuration for the transformation
	if err := json.Unmarshal(data, &parameterMapping); err != nil {
		return err
	}
	// now build the additional query parameter raw string
	var additionalRawQueryParameters string
	for parameter, values := range parameterMapping {
		for _, value := range values {
			additionalRawQueryParameters += fmt.Sprintf("&%s=%s",
				url.QueryEscape(parameter), url.QueryEscape(value))
		}
	}
	// now append the additional query parameters to the original query
	r.URL.RawQuery += additionalRawQueryParameters
	return nil
}

func (o AddQueryParameters) ApplyAfter(r *http.Response, data []byte) error {
	return errors.New("function is pre-proxy transformation")
}
