package globals

import "github.com/qustavo/dotsql"

// This file contains globally shared variables (e.g., service name, sql queries)

// ServiceName contains the global identifier for the service
const ServiceName = "external-apis"

// SqlQueries contains the prepared sql queries from the resources folder
var SqlQueries *dotsql.DotSql

// RequiredUserGroup contains the user group read from the scope file needed
// to access this service
var RequiredUserGroup string
