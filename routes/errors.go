package routes

import wisdomType "github.com/wisdom-oss/commonTypes/v2"

// ErrInvalidUUID is returned if a request parameter (path or query) is expected
// to be a UUID but is not a valid argument for pgtype.UUID's Scan function
var ErrInvalidUUID = wisdomType.WISdoMError{
	Type:   "",
	Status: 400,
	Title:  "Invalid UUID",
	Detail: "The UUID provided as parameter is not in a valid format",
}

// ErrDatasourceNotFound is returned if the specified data source is not stored
// in the system
var ErrDatasourceNotFound = wisdomType.WISdoMError{
	Type:   "",
	Status: 404,
	Title:  "No Datasource Found",
	Detail: "The datasource specified in this request is not present in the system",
}

var ErrInvalidDatasourceAPI = wisdomType.WISdoMError{
	Type:   "",
	Status: 400,
	Title:  "Invalid Datasource API Configuration",
	Detail: "The configuration stored in the database is not valid for usage with the proxy. Please check the configuration",
}
