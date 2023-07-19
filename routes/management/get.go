package managementRoutes

import (
	"encoding/json"
	"external-api-service/globals"
	"external-api-service/types"
	"github.com/blockloop/scan/v2"
	"net/http"
)

// AllExternalDataSources takes a new http request and returns all external
// data sources currently stored in the database with their complete metadata
// and api proxy endpoint
func AllExternalDataSources(w http.ResponseWriter, r *http.Request) {
	// access the native error channel for handling query errors
	nativeErrorChannel := r.Context().Value("nativeErrorChannel").(chan error)
	nativeErrorHandled := r.Context().Value("nativeErrorHandled").(chan bool)

	// now query the whole info schema containing all information
	rows, err := globals.SqlQueries.Query(globals.Db, "get-all-information")
	if err != nil {
		// send the error into the error handling channel and wait until the
		// error has been processed
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}

	// now scan the query results into an array of the output structs
	var externalDataSources []types.ExternalDataSource
	err = scan.Rows(&externalDataSources, rows)
	if err != nil {
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}

	// now set the correct response content type
	w.Header().Set("Content-Type", "text/json")

	// now encode the array of data sources and send it back
	err = json.NewEncoder(w).Encode(externalDataSources)
	if err != nil {
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}
}
