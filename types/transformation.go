package types

import (
	"external-api-service/enums"
	"external-api-service/interfaces"
	"external-api-service/transformations"
	"fmt"
	"net/http"
)

type TransformationDefinition struct {
	// ID contains the identifier of this transformation operation
	ID string `json:"id" db:"id"`
	// Datasource contains the UUID of the datasource on which the operation
	// shall be executed in case of a request
	Datasource string `json:"datasource" db:"datasource"`
	// Action specifies which type of transformation is happening to the
	// request or response
	Action enums.TransformationType `json:"action" db:"action"`
	// Priority specifies how important the transformation is and therefore its
	// precedence before other transformations. The higher, the better
	Priority int64 `json:"priority" db:"priority"`
	// Data contains the data needed for executing the transformation.
	// It needs to be marshalled into the needed type before executing the
	// transformation
	Data []byte `json:"actionData" db:"data"`
}

// transformationFunctions is a mapping that takes an [enums.TransformationType]
// as a key and a struct implementing the apply function from the
// TransformationFunction
var transformationFunctions = map[enums.TransformationType]interfaces.TransformationFunction{
	enums.TRANSFORM_SORT_QUERY_PARAMETERS: transformations.SortQueryParameter{},
	enums.TRANSFORM_ADD_QUERY_PARAMETERS:  transformations.AddQueryParameters{},
	enums.TRANSFORM_SET_QUERY_PARAMETERS:  transformations.SortQueryParameter{},
}

// ApplyBefore gets the struct mapped to the transformation type and runs the Apply
// function on the struct supplying the data directly from the struct
func (td TransformationDefinition) ApplyBefore(r *http.Request) error {
	// get the struct implementing the apply function
	targetStruct, ok := transformationFunctions[td.Action]
	if !ok {
		return fmt.Errorf("no struct mapped implementing 'ApplyBefore()' to %s", td.Action)
	}
	return targetStruct.ApplyBefore(r, td.Data)
}

// ApplyAfter gets the struct mapped to the transformation type and runs the Apply
// function on the struct supplying the data directly from the struct
func (td TransformationDefinition) ApplyAfter(r *http.Response) error {
	// get the struct implementing the apply function
	targetStruct, ok := transformationFunctions[td.Action]
	if !ok {
		return fmt.Errorf("no struct mapped implementing 'ApplyAfter()' to %s", td.Action)
	}
	return targetStruct.ApplyAfter(r, td.Data)
}
