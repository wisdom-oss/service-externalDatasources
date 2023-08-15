package enums

import (
	wisdomUtils "github.com/wisdom-oss/microservice-utils"
)

type TransformationType string

const (
	TRANSFORM_SORT_QUERY_PARAMETERS TransformationType = "sortQueryParameters"
	TRANSFORM_SET_QUERY_PARAMETERS  TransformationType = "setQueryParameters"
	TRANSFORM_ADD_QUERY_PARAMETERS  TransformationType = "addQueryParameters"
)

// preRequestTransformations contains all transformations that will happen
// before sending the request to the target
var preRequestTransformations = []TransformationType{
	TRANSFORM_ADD_QUERY_PARAMETERS,
	TRANSFORM_SORT_QUERY_PARAMETERS,
	TRANSFORM_SET_QUERY_PARAMETERS,
}

// postRequestTransformations contains all transformations that will happen
// after sending the request to the target
var postRequestTransformations = []TransformationType{}

// IsPre checks if the [TransformationType] is to be executed before the request
// gets sent to its destination
func (tt TransformationType) IsPre() bool {
	return wisdomUtils.ArrayContains(preRequestTransformations, tt)
}

// IsPost checks if the [TransformationType] is to be executed after the request
// gets sent to its destination
func (tt TransformationType) IsPost() bool {
	return wisdomUtils.ArrayContains(postRequestTransformations, tt)
}
