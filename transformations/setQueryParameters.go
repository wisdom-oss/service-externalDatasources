package transformations

import (
	"encoding/json"
	"errors"
	"net/http"
)

type SetQueryParameters map[string][]string

func (parameterMapping SetQueryParameters) ApplyBefore(r *http.Request, data []byte) error {
	// try to parse the configuration for the transformation
	if err := json.Unmarshal(data, &parameterMapping); err != nil {
		return err
	}
	// get the original query parameters
	parameters := r.URL.Query()
	// now iterate above the configured parameters and replace/add the values
	// to the mapping
	for key, values := range parameterMapping {
		if parameters.Has(key) {
			parameters[key] = values
		} else {
			for _, value := range values {
				parameters.Add(key, value)
			}
		}
	}
	r.URL.RawQuery = parameters.Encode()
	return nil
}

func (o SetQueryParameters) ApplyAfter(r http.Response, data []byte) error {
	return errors.New("function is pre-proxy transformation")
}
