package transformations

import (
	"net/http"
	"net/url"
)

const SetQueryParametersKey = "replace-query-parameter"

type SetQueryParameters struct{}

func (t SetQueryParameters) ApplyBefore(r *http.Request, options map[string]interface{}) error {
	queryParameters := new(url.Values)
	for key, value := range options {
		stringValue, isString := value.(string)
		if !isString {
			return ErrValueNotString
		}
		queryParameters.Add(key, stringValue)
	}
	r.URL.RawQuery = queryParameters.Encode()
	return nil
}
