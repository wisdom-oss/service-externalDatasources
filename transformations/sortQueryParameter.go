package transformations

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SortQueryParameter is a type alias on a string array. However, on this
// type alias the interface TransformationFunction is implemented. This leads
// to any variable having this type to execute the transformation
type SortQueryParameter []string

func (sqp SortQueryParameter) ApplyBefore(r *http.Request, data []byte) error {
	if err := json.Unmarshal(data, &sqp); err != nil {
		return err
	}
	fmt.Println(sqp)
	return nil
}

func (sqp SortQueryParameter) ApplyAfter(r *http.Response, data []byte) error {
	return nil
}
