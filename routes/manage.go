package routes

import (
	"encoding/json"
	"external-api-service/globals"
	"external-api-service/types"
	"github.com/blockloop/scan/v2"
	"net/http"
)

func GetAllExternalAPIs(w http.ResponseWriter, r *http.Request) {
	// get the error channels
	nativeErrorChannel := r.Context().Value("nativeErrorChannel").(chan error)
	nativeErrorHandled := r.Context().Value("nativeErrorHandled").(chan bool)
	//wisdomErrorChannel := r.Context().Value("wisdomErrorChannel").(chan string)
	// query the whole database
	rows, err := globals.SqlQueries.Query(globals.Db, "get-all-information")
	if err != nil {
		nativeErrorChannel <- err
		return
	}

	var externalSources []types.ExternalDataSource
	err = scan.Rows(&externalSources, rows)
	if err != nil {
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}
	w.Header().Set("Content-Type", "text/json")
	_ = json.NewEncoder(w).Encode(externalSources)
	return
}
