package managementRoutes

import (
	"database/sql"
	"encoding/json"
	"errors"
	"external-api-service/globals"
	"external-api-service/types"
	"github.com/blockloop/scan/v2"
	"github.com/go-chi/chi/v5"
	"net/http"
	"regexp"
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

func SingleDataSource(w http.ResponseWriter, r *http.Request) {
	// get the error channels
	nativeErrorChannel := r.Context().Value("nativeErrorChannel").(chan error)
	nativeErrorHandled := r.Context().Value("nativeErrorHandled").(chan bool)
	wisdomErrorChannel := r.Context().Value("wisdomErrorChannel").(chan string)
	wisdomErrorHandled := r.Context().Value("wisdomErrorHandled").(chan bool)

	// now try to get the uuid of the data source
	uuid := chi.URLParam(r, "dataSourceUUID")

	// now check if the uuid is a valid uuid using regex
	regex := regexp.MustCompile(`^[[:xdigit:]]{8}-?[[:xdigit:]]{4}-?[[:xdigit:]]{4}-?[[:xdigit:]]{4}-?[[:xdigit:]]{12}$`)
	if !regex.MatchString(uuid) {
		wisdomErrorChannel <- "INVALID_UUID_FORMAT"
		<-wisdomErrorHandled
		return
	}

	// now try to retrieve a data source from the database
	dataSourceRow, err := globals.SqlQueries.Query(globals.Db, "get-single-source", uuid)
	if err != nil {
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}

	// now try to parse the data source and handle errors
	var dataSource types.ExternalDataSource
	err = scan.Row(&dataSource, dataSourceRow)
	if errors.As(err, &sql.ErrNoRows) {
		wisdomErrorChannel <- "NO_DATASOURCE_FOUND"
		<-wisdomErrorHandled
		return
	}
	if err != nil {
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}

	// now set the correct output format and encode the found data source
	w.Header().Set("Content-Type", "text/json")
	err = json.NewEncoder(w).Encode(dataSource)
	if err != nil {
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}
}
