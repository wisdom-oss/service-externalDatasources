package transformations

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// SortQueryParameter is a type alias on a string array. However, on this
// type alias the interface TransformationFunction is implemented. This leads
// to any variable having this type to execute the transformation
type SortQueryParameter []string

// ApplyBefore implements a part of the TransformationFunction interface
//
// In this particular instance the function sorts the request parameters as
// configured and adds the remaining unsorted query parameters back at the end
func (sqp SortQueryParameter) ApplyBefore(r *http.Request, data []byte) error {
	if err := json.Unmarshal(data, &sqp); err != nil {
		return err
	}
	// get the request parameters and their original values
	originalParameters := r.URL.Query()
	// now create an empty string which will contain the new sorted raw query
	// parameters
	var rawSortedQueryParameters string
	// now iterate through the parsed query parameters
	for _, param := range sqp {
		if !originalParameters.Has(param) {
			continue
		}
		// retrieve the original parameters values
		values := originalParameters[param]
		// and iterate over them to add them to the sorted parameter string
		for _, value := range values {
			rawSortedQueryParameters += fmt.Sprintf("&%s=%s", url.QueryEscape(param), url.QueryEscape(value))
		}
		originalParameters.Del(param)
	}
	// now add all the remaining parameters as they were parsed when receiving
	// the request
	for originalParam, originalValues := range originalParameters {
		// iterate over them to add them to the sorted parameter string
		for _, value := range originalValues {
			rawSortedQueryParameters += fmt.Sprintf("&%s=%s", url.QueryEscape(originalParam), url.QueryEscape(value))
		}
	}
	r.URL.RawQuery = rawSortedQueryParameters
	return nil
}

func (sqp SortQueryParameter) ApplyAfter(r *http.Response, data []byte) error {
	return nil
}
