package globals

import (
	"github.com/qustavo/dotsql"

	"external-api-service/types"
)

// This file contains globally shared variables (e.g., service name, sql queries)

// ServiceName contains the global identifier for the service
const ServiceName = "external-apis"

// SqlQueries contains the prepared sql queries from the resources folder
var SqlQueries *dotsql.DotSql

// RequiredUserGroup contains the user group read from the scope file needed
// to access this service
var RequiredUserGroup string

// Environment contains a mapping between the environment variables and the values
// they were set to. However, this variable only contains the configured environment
// variables
var Environment map[string]string = make(map[string]string)

// Errors contains all errors that have been predefined in the "errors.json" file.
var Errors []types.WISdoMError
