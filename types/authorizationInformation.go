package types

type AuthorizationDataLocation string

const (
	LOCATION_HEADER AuthorizationDataLocation = "header"
	LOCATION_QUERY  AuthorizationDataLocation = "query"
)

type AuthorizationInformation struct {
	Location AuthorizationDataLocation
	Key      string
	Value    string
}

// TODO: implement conversion/parsing functions
