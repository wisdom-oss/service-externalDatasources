package transformations

import (
	"errors"
	"net/http"
)

const AddQueryParametersKey = "add-query-parameter"

// AddQueryParameters is a type representing a transformation function
// that adds query parameters to a URL.
type AddQueryParameters struct{}

var ErrValueNotString = errors.New("query parameter value is not string")

func (t AddQueryParameters) ApplyBefore(r *http.Request, options map[string]interface{}) error {
	queryParameters := r.URL.Query()
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
