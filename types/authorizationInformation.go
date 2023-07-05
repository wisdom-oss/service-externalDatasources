package types

import (
	"fmt"
)

type AuthorizationDataLocation string

const (
	LOCATION_HEADER AuthorizationDataLocation = "header"
	LOCATION_QUERY  AuthorizationDataLocation = "query"
)

func (loc AuthorizationDataLocation) ToString() string {
	return fmt.Sprintf("%s", loc)
}

type AuthorizationInformation struct {
	Location AuthorizationDataLocation
	Key      string
	Value    string
}
