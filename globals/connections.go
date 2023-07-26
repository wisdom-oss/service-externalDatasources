package globals

import (
	"database/sql"
	"net/http"
)

// This file contains all globally shared connections (e.g., Databases)

// Db contains the globally available connection to the database
var Db *sql.DB

// HttpClient contains the globally used http client in the proxy switch
var HttpClient = &http.Client{}
