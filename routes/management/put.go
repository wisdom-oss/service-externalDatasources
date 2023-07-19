package managementRoutes

import (
	"database/sql"
	"encoding/json"
	"errors"
	"external-api-service/globals"
	"external-api-service/types"
	"github.com/blockloop/scan/v2"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"net/http"
	"regexp"
	"strings"
)

// ReplaceDataSourceRepresentation takes an incoming request using a multipart
// form containing the name, description, metadata and api endpoint description
// and replaces the current representation completely with the data contained
// in the request
func ReplaceDataSourceRepresentation(w http.ResponseWriter, r *http.Request) {
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
	result, err := globals.SqlQueries.Query(globals.Db, "data-source-exists", uuid)
	if err != nil {
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}

	// now try to parse the data source and handle errors
	var dataSourceExists bool
	err = scan.Row(&dataSourceExists, result)
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

	if !dataSourceExists {
		wisdomErrorChannel <- "NO_DATASOURCE_FOUND"
		<-wisdomErrorHandled
		return
	}

	// now try to parse the multipart form that has been sent
	if err := r.ParseMultipartForm(MaxMultipartFormMemoryBytes); err != nil {
		// send the error into the error handling channel and wait until the
		// error has been processed
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}

	// now validate that the transmitted and parsed data contains a name for the
	// new data source
	names, ok := r.MultipartForm.Value["name"]
	if !ok {
		wisdomErrorChannel <- "NO_DATASOURCE_NAME"
		<-wisdomErrorHandled
		return
	}

	// now check that at least one name has been set and extract it
	var name string
	switch len(names) {
	case 0:
		// return an error indicating that no data source name has been set
		wisdomErrorChannel <- "NO_DATASOURCE_NAME"
		<-wisdomErrorHandled
		return
	case 1:
		// just get the data source name
		name = names[0]
		break
	default:
		// since at least two names were set, log an error and add an error
		// message to the headers
		log.Warn().Msg("multiple names set for data source. using first name")
		w.Header().Add("x-warning", "multiple data souce names")
		name = names[0]
		break
	}

	// now validate that the extracted data source name is not empty
	name = strings.TrimSpace(name)
	if name == "" {
		wisdomErrorChannel <- "EMPTY_DATASOURCE_NAME"
		<-wisdomErrorHandled
		return
	}

	// now check if a description was sent. if one was sent, also check that
	// it is not empty
	var description *string
	descriptions, ok := r.MultipartForm.Value["description"]
	if ok {
		// check that at least one description has been sent and extract it
		switch len(descriptions) {
		case 0:
			// return an error indicating that no data source description has been set
			wisdomErrorChannel <- "NO_DATASOURCE_DESCRIPTION"
			<-wisdomErrorHandled
			return
		case 1:
			description = &descriptions[0]
			break
		default:
			// since at least two descriptions were set, log an error and
			// add an error message to the headers
			log.Warn().Msg("multiple descriptions set for data source. using first description")
			w.Header().Add("x-warning", "multiple data source descriptions")
			description = &descriptions[0]
			break
		}

		// now check if the description is not empty
		*description = strings.TrimSpace(*description)
		if *description == "" {
			wisdomErrorChannel <- "EMPTY_DATASOURCE_DESCRIPTION"
			<-wisdomErrorHandled
			return
		}
	}

	// now do the same with the json encoded metadata
	var metadata types.Metadata
	metadataObjects, ok := r.MultipartForm.Value["metadata"]
	if ok {
		var rawMetadataObject string
		// check that at least one metadata object was sent
		switch len(metadataObjects) {
		case 0:
			wisdomErrorChannel <- "METADATA_MISSING"
			<-wisdomErrorHandled
			return
		case 1:
			rawMetadataObject = metadataObjects[0]
			break
		default:
			// since at least two metadata objects were set, log an error and
			// add an error message to the headers
			log.Warn().Msg("multiple metadata objects set for data source. using first object")
			w.Header().Add("x-warning", "multiple metadata objects")
			rawMetadataObject = metadataObjects[0]
			break
		}

		// now check that the metadata object actually contains a value
		rawMetadataObject = strings.TrimSpace(rawMetadataObject)
		if rawMetadataObject == "" {
			wisdomErrorChannel <- "METADATA_MISSING"
			<-wisdomErrorHandled
			return
		}

		// now try to decode the rawMetadataObject into the struct
		if err := json.Unmarshal([]byte(rawMetadataObject), &metadata); err != nil {
			// now check if an error happened during the unmarshalling
			var unmarshalError *json.InvalidUnmarshalError
			var syntaxError *json.SyntaxError
			var typeError *json.UnmarshalTypeError
			switch {
			case errors.As(err, &unmarshalError):
				nativeErrorChannel <- err
				<-nativeErrorHandled
				return
			case errors.As(err, &syntaxError):
				wisdomErrorChannel <- "INVALID_METADATA_JSON"
				<-wisdomErrorHandled
				return
			case errors.As(err, &typeError):
				wisdomErrorChannel <- "INVALID_METADATA"
				<-wisdomErrorHandled
				return
			default:
				nativeErrorChannel <- err
				<-nativeErrorHandled
				return
			}
		}

	}

	// and do the name again with the api endpoint data
	var apiEndpoint types.API
	apiObjects, ok := r.MultipartForm.Value["api"]
	if ok {
		var rawAPIObject string
		// check that at least one metadata object was sent
		switch len(apiObjects) {
		case 0:
			wisdomErrorChannel <- "API_CONFIGURATION_MISSING"
			<-wisdomErrorHandled
			return
		case 1:
			rawAPIObject = apiObjects[0]
			break
		default:
			// since at least two api objects were set, log an error and add an
			// error message to the headers
			log.Warn().Msg("multiple api objects set for data source. using first object")
			w.Header().Add("x-warning", "multiple api objects")
			rawAPIObject = apiObjects[0]
			break
		}

		// now check that the api object actually contains a value
		rawAPIObject = strings.TrimSpace(rawAPIObject)
		if rawAPIObject == "" {
			wisdomErrorChannel <- "API_CONFIGURATION_MISSING"
			<-wisdomErrorHandled
			return
		}

		// now try to decode the rawAPIObject into the struct
		if err := json.Unmarshal([]byte(rawAPIObject), &apiEndpoint); err != nil {
			// now check if an error happened during the unmarshalling
			var unmarshalError *json.InvalidUnmarshalError
			var syntaxError *json.SyntaxError
			var typeError *json.UnmarshalTypeError
			switch {
			case errors.As(err, &unmarshalError):
				nativeErrorChannel <- err
				<-nativeErrorHandled
				return
			case errors.As(err, &typeError), errors.As(err, &syntaxError):
				wisdomErrorChannel <- "INVALID_API_CONFIGURATION"
				<-wisdomErrorHandled
				return
			default:
				nativeErrorChannel <- err
				<-nativeErrorHandled
				return
			}
		}

	}

	// to secure the following insertion queries, a transaction will be used
	// to roll back the changes in case of an error
	tx, err := globals.Db.BeginTx(r.Context(), nil)
	if err != nil {
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}

	// now create the base entry in the database and retrieve the inserted uuid
	// from the query
	_, err = globals.SqlQueries.Exec(tx, "set-base-data",
		uuid, name, description)
	if err != nil {
		// since an error occurred, roll back the changes made in this
		// transaction
		rbErr := tx.Rollback()
		if rbErr != nil {
			nativeErrorChannel <- rbErr
			<-nativeErrorHandled
			return
		}
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}

	// now update the metadata
	_, err = globals.SqlQueries.Exec(tx, "set-metadata",
		uuid, metadata.Reference, metadata.Origin,
		metadata.DistinctiveFeatures, metadata.UsageRights,
		metadata.UsageDuties, metadata.RealEntities, metadata.LocalExpert,
		metadata.Documentation, metadata.UpdateRate, metadata.Languages,
		metadata.Billing, metadata.Provision, metadata.DerivedFrom,
		metadata.Recent, metadata.Validity, metadata.Duplicates,
		metadata.Errors, metadata.Precision, metadata.Reputation,
		metadata.DataObjectivity, metadata.UsualSurveyMethod,
		metadata.Density, metadata.Coverage,
		metadata.RepresentationConsistency, metadata.LogicalConsistency,
		metadata.DataDelay, metadata.DelayInformationTransmission,
		metadata.PerformanceLimitations, metadata.Availability,
		metadata.GDPRCompliant,
	)
	if err != nil {
		// since an error occurred, roll back the changes made in this
		// transaction
		rbErr := tx.Rollback()
		if rbErr != nil {
			nativeErrorChannel <- rbErr
			<-nativeErrorHandled
			return
		}
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}

	// now update the api data if it is not nil
	_, err = globals.SqlQueries.Exec(tx, "set-api",
		uuid, apiEndpoint.IsSecure, apiEndpoint.Host, apiEndpoint.Port,
		apiEndpoint.Path, apiEndpoint.AdditionalHeaders,
	)
	if err != nil {
		// since an error occurred, roll back the changes made in this
		// transaction
		rbErr := tx.Rollback()
		if rbErr != nil {
			nativeErrorChannel <- rbErr
			<-nativeErrorHandled
			return
		}
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}
	// now commit the changes made in the database
	if err = tx.Commit(); err != nil {
		// since an error occurred, roll back the changes made in this
		// transaction
		rbErr := tx.Rollback()
		if rbErr != nil {
			nativeErrorChannel <- rbErr
			<-nativeErrorHandled
			return
		}
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}

	// now try to retrieve the data source from the database
	dataSourceRow, err := globals.SqlQueries.Query(globals.Db, "get-single-source", uuid)
	if err != nil {
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}

	// now try to parse the data source and handle errors
	var dataSource types.ExternalDataSource
	err = scan.Row(&dataSource, dataSourceRow)
	if err != nil {
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}

	// now set the correct output format and encode the found data source
	w.Header().Set("Content-Type", "text/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(dataSource)
	if err != nil {
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}

}
