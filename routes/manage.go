package routes

import (
	"encoding/json"
	"external-api-service/globals"
	"external-api-service/types"
	"fmt"
	"github.com/blockloop/scan/v2"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
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

// NewExternalAPI allows the creation of an external data source
func NewExternalAPI(w http.ResponseWriter, r *http.Request) {
	// get the error channels
	nativeErrorChannel := r.Context().Value("nativeErrorChannel").(chan error)
	nativeErrorHandled := r.Context().Value("nativeErrorHandled").(chan bool)
	wisdomErrorChannel := r.Context().Value("wisdomErrorChannel").(chan string)
	wisdomErrorHandled := r.Context().Value("wisdomErrorHandled").(chan bool)
	// since the data is expected as form-data parse the form to start handling it
	err := r.ParseMultipartForm(8589934592)
	if err != nil {
		// handle the error that occurred
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}

	// now check the multipart form for the name for the new external data source
	nameList, ok := r.MultipartForm.Value["name"]
	if !ok {
		wisdomErrorChannel <- "NO_DATASOURCE_NAME"
		<-wisdomErrorHandled
		return
	}
	if len(nameList) > 1 {
		log.Warn().Msg("multiple names found in request, only using the first")
	}

	datasourceName := strings.TrimSpace(nameList[0])
	if datasourceName == "" {
		wisdomErrorChannel <- "EMPTY_DATASOURCE_NAME"
		<-wisdomErrorHandled
		return
	}
	// now check the multipart form for the name for the new external data source
	var datasourceDescription *string
	descriptionList, ok := r.MultipartForm.Value["description"]
	if ok {
		s := strings.TrimSpace(descriptionList[0])

		if s == "" {
			wisdomErrorChannel <- "EMPTY_DATASOURCE_NAME"
			<-wisdomErrorHandled
			return
		}
		datasourceDescription = &s
	}
	if len(descriptionList) > 1 {
		log.Warn().Msg("multiple names found in request, only using the first")
	}

	// now check the multipart form for the metadata
	var metadata *types.Metadata
	metadataList, metadataSet := r.MultipartForm.Value["metadata"]
	if !metadataSet {
		log.Warn().Msg("no metadata sent in creation request")
		metadata = nil
	} else {
		if len(metadataList) > 1 {
			log.Warn().Msg("multiple metadata objects in request. only using first")
		}
		// now get the first metadata object
		rawMetadata := strings.TrimSpace(metadataList[0])

		// now check if the object even contains text
		if rawMetadata == "" {
			wisdomErrorChannel <- "METADATA_MISSING"
			<-wisdomErrorHandled
			return
		}
		// now try to parse the text as json
		decoder := json.NewDecoder(strings.NewReader(rawMetadata))
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&metadata)
		switch err.(type) {
		case *json.InvalidUnmarshalError:
			nativeErrorChannel <- err
			<-nativeErrorHandled
			return
		case *json.SyntaxError:
			wisdomErrorChannel <- "INVALID_METADATA_JSON"
			<-wisdomErrorHandled
			return
		case *json.UnmarshalTypeError:
			wisdomErrorChannel <- "INVALID_METADATA"
			<-wisdomErrorHandled
			return
		default:
			if err != nil {
				// now check if the error message indicates that a field is not valid
				var illegalField string
				fmt.Println(err.Error())
				_, parseError := fmt.Sscanf(err.Error(), `json: unknown field %s`, &illegalField)
				if parseError != nil {
					nativeErrorChannel <- err
					<-nativeErrorHandled
				} else {
					wisdomErrorChannel <- "UNKNOWN_METADATA_FIELD"
					<-wisdomErrorHandled
				}
				return
			}
		}
	}

	// now check the multipart form for the metadata
	var api *types.Metadata
	apiList, apiSet := r.MultipartForm.Value["api"]
	if !apiSet {
		log.Warn().Msg("no api configuration sent in creation request")
		api = nil
	} else {
		if len(apiList) > 1 {
			log.Warn().Msg("multiple api configuration objects in request. only using first")
		}
		// now get the first api config object
		rawApi := strings.TrimSpace(metadataList[0])

		// now check if the object even contains text
		if rawApi == "" {
			wisdomErrorChannel <- "API_CONFIGURATION_MISSING"
			<-wisdomErrorHandled
			return
		}

		// now try to parse the text as json
		decoder := json.NewDecoder(strings.NewReader(rawApi))
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&api)
		switch err.(type) {
		case *json.InvalidUnmarshalError:
			nativeErrorChannel <- err
			<-nativeErrorHandled
			return
		case *json.SyntaxError:
		case *json.UnmarshalTypeError:
			wisdomErrorChannel <- "INVALID_API_CONFIGURATION"
			<-wisdomErrorHandled
			return
		default:
			if err != nil {
				// now check if the error message indicates that a field is not valid
				var illegalField string
				fmt.Println(err.Error())
				_, parseError := fmt.Sscanf(err.Error(), `json: unknown field %s`, &illegalField)
				if parseError != nil {
					nativeErrorChannel <- err
					<-nativeErrorHandled
				} else {
					wisdomErrorChannel <- "UNKNOWN_API_CONFIGURATION_FIELD"
					<-wisdomErrorHandled
				}
				return
			}
		}
	}

	transaction, err := globals.Db.BeginTx(r.Context(), nil)
	if err != nil {
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}

	result, err := globals.SqlQueries.QueryRow(transaction, "add-base-data", datasourceName, datasourceDescription)
	if err != nil {
		rollbackErr := transaction.Rollback()
		if rollbackErr != nil {
			nativeErrorChannel <- err
			<-nativeErrorHandled
			return
		}
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}

	var uuid string
	err = result.Scan(&uuid)
	if err != nil {
		rollbackErr := transaction.Rollback()
		if rollbackErr != nil {
			nativeErrorChannel <- err
			<-nativeErrorHandled
			return
		}
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}

	_, err = globals.SqlQueries.Exec(transaction, "add-metadata",
		uuid,
		metadata.Reference,
		metadata.Origin,
		metadata.DistinctiveFeatures,
		metadata.UsageRights,
		metadata.UsageDuties,
		metadata.RealEntities,
		metadata.LocalExpert,
		metadata.Documentation,
		metadata.UpdateRate,
		metadata.Languages,
		metadata.Billing,
		metadata.Provision,
		metadata.DerivedFrom,
		metadata.Recent,
		metadata.Validity,
		metadata.Duplicates,
		metadata.Errors,
		metadata.Precision,
		metadata.Reputation,
		metadata.DataObjectivity,
		metadata.UsualSurveyMethod,
		metadata.Density,
		metadata.Coverage,
		metadata.RepresentationConsistency,
		metadata.LogicalConsistency,
		metadata.DataDelay,
		metadata.DelayInformationTransmission,
		metadata.PerformanceLimitations,
		metadata.Availability,
		metadata.GDPRCompliant,
	)
	if err != nil {
		rollbackErr := transaction.Rollback()
		if rollbackErr != nil {
			nativeErrorChannel <- err
			<-nativeErrorHandled
			return
		}
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}

	err = transaction.Commit()
	if err != nil {
		rollbackErr := transaction.Rollback()
		if rollbackErr != nil {
			nativeErrorChannel <- err
			<-nativeErrorHandled
			return
		}
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}

}
