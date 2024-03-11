package transformations

import (
	"net/http"
	"net/url"
	"slices"
)

const SortQueryParametersKey = "sort-query-parameter"

type SortQueryParameters struct{}

func (t SortQueryParameters) ApplyBefore(r *http.Request, options map[string]interface{}) error {
	originalQueryParameter := r.URL.Query()
	var queryParameterKeys []string
	for key, _ := range originalQueryParameter {
		queryParameterKeys = append(queryParameterKeys, key)
	}
	slices.Sort(queryParameterKeys)
	queryParameters := new(url.Values)
	for _, key := range queryParameterKeys {
		parameterArray := originalQueryParameter[key]
		for _, value := range parameterArray {
			queryParameters.Add(key, value)
		}
	}
	r.URL.RawQuery = queryParameters.Encode()
	return nil
}
