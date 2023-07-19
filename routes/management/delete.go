package managementRoutes

import (
	"external-api-service/globals"
	"github.com/go-chi/chi/v5"
	"net/http"
	"regexp"
)

func DeleteExternalDataSource(w http.ResponseWriter, r *http.Request) {
	// access the error channels for handling errors
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
	_, err := globals.SqlQueries.Query(globals.Db, "delete-datasource", uuid)
	if err != nil {
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}

	// since there is no content available, return a 204
	w.WriteHeader(http.StatusNoContent)
}
