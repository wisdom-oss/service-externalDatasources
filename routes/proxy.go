package routes

import (
	"database/sql"
	"errors"
	"external-api-service/enums"
	"external-api-service/globals"
	"external-api-service/types"
	"fmt"
	"github.com/blockloop/scan/v2"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// ProxySwitch is the first handler for a new proxy request.
// It validates the provided data source id and checks if an api endpoint is
// currently configured for the data source.
// Afterward, the switch checks if any transformations are to be executed, and
// uses the appropriate request function to proxy the request.
func ProxySwitch(w http.ResponseWriter, r *http.Request) {
	// get the error handlers from the request context
	nativeErrorChannel := r.Context().Value("nativeErrorChannel").(chan error)
	nativeErrorHandled := r.Context().Value("nativeErrorHandled").(chan bool)
	wisdomErrorChannel := r.Context().Value("wisdomErrorChannel").(chan string)
	wisdomErrorHandled := r.Context().Value("wisdomErrorHandled").(chan bool)

	// now extract the data source id from the url
	uuid := chi.URLParam(r, "datasourceUUID")

	// now check if the uuid matches a configured data source
	regex := regexp.MustCompile(`^[[:xdigit:]]{8}-?[[:xdigit:]]{4}-?[[:xdigit:]]{4}-?[[:xdigit:]]{4}-?[[:xdigit:]]{12}$`)
	if !regex.MatchString(uuid) {
		wisdomErrorChannel <- "INVALID_UUID_FORMAT"
		<-wisdomErrorHandled
		return
	}

	// since the uuid is valid now query the database if the datasource exists
	result, err := globals.SqlQueries.Query(globals.Db, "data-source-exists", uuid)
	if err != nil {
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}

	// now parse the result into a boolean and check if at least one row was
	// returned
	var datasourceFound bool
	err = scan.Row(&datasourceFound, result)
	if errors.As(err, &sql.ErrNoRows) || !datasourceFound {
		wisdomErrorChannel <- "NO_DATASOURCE_FOUND"
		<-wisdomErrorHandled
		return
	}
	if err != nil {
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}

	// now get the api configuration from the database
	apiRow, err := globals.SqlQueries.Query(globals.Db, "get-api", uuid)
	if err != nil {
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}

	// now try to parse the result into the api configuration
	var apiConfiguration types.API
	err = scan.Row(&apiConfiguration, apiRow)
	if errors.As(err, &sql.ErrNoRows) {
		// there is no api configuration for this datasource
		wisdomErrorChannel <- "NO_API_CONFIGURED"
		<-wisdomErrorHandled
		return
	}
	if err != nil {
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}

	// now check if the api configuration is usable
	if !apiConfiguration.IsValid() {
		wisdomErrorChannel <- "STORED_API_CONFIGURATION_INVALID"
		<-wisdomErrorHandled
		return
	}

	// now build the uri parts according to the configuration
	var uriSchema, host, path string
	var port uint16

	switch *apiConfiguration.IsSecure {
	case true:
		uriSchema = "https://"
		port = 443
		break
	case false:
		uriSchema = "http://"
		port = 80
	default:
		uriSchema = "https://"
		port = 443
		break
	}
	// since the host is the only required field, it may be set directly
	host = *apiConfiguration.Host

	// since the path is optional check before setting it from the configuration
	if apiConfiguration.Path != nil {
		path = strings.TrimPrefix(*apiConfiguration.Path, `/`)
	} else {
		path = ""
	}

	// the same with the port
	if apiConfiguration.Port != nil {
		port = *apiConfiguration.Port
	}

	// now get the path that has been requested via the incoming url and will
	// be added to the api configuration
	extraPath := chi.URLParam(r, "*")

	// now build the uri for the request to the data source
	uri := fmt.Sprintf("%s%s:%d/%s/%s", uriSchema, host, port, path, extraPath)

	// now build the request using the same verb used to access the service
	// and reusing the request body
	request, err := http.NewRequest(r.Method, uri, r.Body)
	if err != nil {
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}

	// now extract the query parameters
	queryParameters := request.URL.Query()
	// now copy the query parameters from the incoming request to the outgoing
	for parameter, values := range r.URL.Query() {
		for _, value := range values {
			queryParameters.Add(parameter, value)
		}
	}
	request.URL.RawQuery = queryParameters.Encode()

	// now check if there are any additional headers to be set and add them
	// if needed.
	if apiConfiguration.AdditionalHeaders != nil {
		for _, header := range *apiConfiguration.AdditionalHeaders {
			request.Header.Add(header.X, header.Y)
		}
	}

	// now get the transformations from the database
	transformationRows, err := globals.SqlQueries.Query(globals.Db, "get-transformations-for-ds", uuid)
	if err != nil {
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}
	// parse the available transformations
	var transformations []types.TransformationDefinition
	err = scan.Rows(&transformations, transformationRows)

	// now set the proxy mode and assume a simple proxy and split the
	// transformations into the appropriate categories
	var preRequestTransformations, postRequestTransformations []types.TransformationDefinition
	proxyMode := enums.PROXY_SIMPLE
	for _, transformation := range transformations {
		// unset the simple proxy since at least one transformation is set
		proxyMode = proxyMode &^ enums.PROXY_SIMPLE
		if transformation.Action.IsPre() {
			proxyMode = proxyMode | enums.PROXY_TRANSFORM_REQUEST
			preRequestTransformations = append(preRequestTransformations, transformation)
		}
		if transformation.Action.IsPost() {
			proxyMode = proxyMode | enums.PROXY_TRANSFORM_RESPONSE
			postRequestTransformations = append(postRequestTransformations, transformation)
		}
	}

	// this checks if the proxy mode is not zero
	if proxyMode == 0 {
		wisdomErrorChannel <- "UNABLE_TO_DETERMINE_PROXY_MODE"
		<-wisdomErrorHandled
		return
	}

	// now execute the pre-request transformations if needed
	if proxyMode&enums.PROXY_TRANSFORM_REQUEST != 0 {
		for _, tf := range preRequestTransformations {
			err = tf.ApplyBefore(request)
			if err != nil {
				nativeErrorChannel <- err
				<-nativeErrorHandled
				return
			}
		}
	}

	requestStart := time.Now()
	res, err := globals.HttpClient.Do(request)
	if err != nil {
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}
	// now calculate the time that has passed between sending the request and
	// completely getting a response
	requestDuration := time.Since(requestStart)
	serverTimingValue := fmt.Sprintf("remote;dur=%f", requestDuration.Seconds())
	// now set the request duration as header
	w.Header().Add("Server-Timing", serverTimingValue)

	// retrieve the content type header from the response and set it in the
	// proxy stream back
	contentType := res.Header.Get("Content-Type")
	w.Header().Set("Content-Type", contentType)

	// now execute the post-request transformations if needed
	if proxyMode&enums.PROXY_TRANSFORM_RESPONSE != 0 {
		for _, tf := range postRequestTransformations {
			err = tf.ApplyAfter(res)
			if err != nil {
				nativeErrorChannel <- err
				<-nativeErrorHandled
				return
			}
		}
	}

	// now copy the response from the datasource to the writer
	_, err = io.Copy(w, res.Body)
	if err != nil {
		nativeErrorChannel <- err
		<-nativeErrorHandled
		return
	}
}
