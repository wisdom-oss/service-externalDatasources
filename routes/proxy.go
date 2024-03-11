package routes

import (
	"fmt"
	"io"
	"net/http"
	"path"
	"reflect"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	wisdomMiddleware "github.com/wisdom-oss/microservice-middlewares/v4"

	"external-api-service/enums"
	"external-api-service/globals"
	"external-api-service/interfaces"
	"external-api-service/transformations"
	"external-api-service/types"
)

// ProxySwitch is the first handler for a new proxy request.
// It validates the provided data source id and checks if an api endpoint is
// currently configured for the data source.
// Afterward, the switch checks if any transformations are to be executed, and
// uses the appropriate request function to proxy the request.
func ProxySwitch(w http.ResponseWriter, r *http.Request) {
	// get the error handlers from the request context
	errorHandler := r.Context().Value(wisdomMiddleware.ErrorChannelName).(chan<- interface{})
	statusChannel := r.Context().Value(wisdomMiddleware.StatusChannelName).(<-chan bool)

	// now extract the data source id from the url
	uuid := chi.URLParam(r, "datasourceUUID")
	datasourceID := pgtype.UUID{}
	err := datasourceID.Scan(uuid)
	if err != nil {
		errorHandler <- ErrInvalidUUID
		<-statusChannel
		return
	}

	query, err := globals.SqlQueries.Raw("data-source-exists")
	if err != nil {
		errorHandler <- err
		<-statusChannel
		return
	}

	var datasourceExists []bool
	err = pgxscan.Select(r.Context(), globals.Db, &datasourceExists, query, datasourceID)
	if err != nil {
		errorHandler <- err
		<-statusChannel
		return
	}

	if datasourceExists == nil || !datasourceExists[0] {
		errorHandler <- ErrDatasourceNotFound
		<-statusChannel
		return
	}

	query, err = globals.SqlQueries.Raw("get-api")
	if err != nil {
		errorHandler <- err
		<-statusChannel
		return
	}

	var apiConfiguration types.API
	err = pgxscan.Get(r.Context(), globals.Db, &apiConfiguration, query, datasourceID)
	if err != nil {
		errorHandler <- err
		<-statusChannel
		return
	}

	if !apiConfiguration.Valid() {
		errorHandler <- ErrInvalidDatasourceAPI
		<-statusChannel
		return
	}

	var (
		scheme            = "https://"
		serverPath        = "/"
		port       uint16 = 443
	)

	if apiConfiguration.IsSecure == nil || *apiConfiguration.IsSecure {
		scheme = "https://"
		port = 443
	} else {
		scheme = "http://"
		port = 80
	}

	if apiConfiguration.Port != nil {
		port = *apiConfiguration.Port
	}

	if apiConfiguration.Path != nil {
		serverPath = *apiConfiguration.Path
	}

	extraPath := chi.URLParam(r, "*")

	requestPath := serverPath + extraPath
	requestPath = path.Clean(requestPath)

	uri := fmt.Sprintf("%s%s:%d%s", scheme, *apiConfiguration.Host, port, requestPath)
	proxyRequest, err := http.NewRequest(r.Method, uri, r.Body)
	if err != nil {
		errorHandler <- err
		<-statusChannel
		return
	}

	// set some headers for proxy requests to allow the identification of the
	// request's origin
	proxyRequest.Header.Set("X-Forwarded-Host", r.Host)
	proxyRequest.Header.Set("X-Forwarded-Proto", r.URL.Scheme)
	proxyRequest.Header.Set("X-Forwarded-For", r.RemoteAddr)
	proxyRequest.Header.Set("X-Real-IP", r.RemoteAddr)

	// now check if any pre-request transformations are configured
	query, err = globals.SqlQueries.Raw("get-transformations-for-ds")
	if err != nil {
		errorHandler <- err
		<-statusChannel
		return
	}

	var transformationConfigurations []types.TransformationConfiguration
	err = pgxscan.Select(r.Context(), globals.Db, &transformationConfigurations, query, datasourceID)
	if err != nil {
		errorHandler <- err
		<-statusChannel
		return
	}

	var preProxyTransformations map[string]map[string]interface{}
	preProxyInterface := reflect.TypeOf((*interfaces.PreProxyTransformation)(nil)).Elem()
	for _, transformationConfiguration := range transformationConfigurations {
		transformation := transformations.GetTransformation(transformationConfiguration.Action.String)
		if reflect.TypeOf(transformation).Implements(preProxyInterface) {
			preProxyTransformations[transformationConfiguration.Action.String] = transformationConfiguration.Options
		}
	}

	var postProxyTransformations map[string]map[string]interface{}
	postProxyInterface := reflect.TypeOf((*interfaces.PostProxyTransformation)(nil)).Elem()
	for _, transformationConfiguration := range transformationConfigurations {
		transformation := transformations.GetTransformation(transformationConfiguration.Action.String)
		if reflect.TypeOf(transformation).Implements(postProxyInterface) {
			postProxyTransformations[transformationConfiguration.Action.String] = transformationConfiguration.Options
		}
	}

	proxyMode := enums.SimpleProxy
	if len(preProxyTransformations) > 0 {
		proxyMode |= enums.TransformRequest
	}
	if len(postProxyTransformations) > 0 {
		proxyMode |= enums.TransformResponse
	}

	if proxyMode&enums.TransformRequest != 0 {
		requestTransformationStart := time.Now()
		for transformationKey, transformationOptions := range preProxyTransformations {
			transformation := transformations.GetTransformation(transformationKey).(interfaces.PreProxyTransformation)
			err = transformation.ApplyBefore(r, transformationOptions)
			if err != nil {
				errorHandler <- err
				<-statusChannel
				return
			}
		}
		requestTransformationDuration := time.Since(requestTransformationStart)
		timingHeaderValue := fmt.Sprintf("preProxyTransformations;dur=%f", requestTransformationDuration.Seconds())
		w.Header().Add("Server-Timing", timingHeaderValue)
	}

	requestStart := time.Now()
	proxyResponse, err := globals.HttpClient.Do(proxyRequest)
	if err != nil {
		errorHandler <- err
		<-statusChannel
		return
	}
	requestDuration := time.Since(requestStart)
	timingHeaderValue := fmt.Sprintf("proxyRequest;dur=%f", requestDuration.Seconds())
	w.Header().Add("Server-Timing", timingHeaderValue)

	if proxyMode&enums.TransformResponse != 0 {
		responseTransformationStart := time.Now()
		for transformationKey, transformationOptions := range postProxyTransformations {
			transformation := transformations.GetTransformation(transformationKey).(interfaces.PostProxyTransformation)
			err = transformation.ApplyAfter(proxyResponse, transformationOptions)
			if err != nil {
				errorHandler <- err
				<-statusChannel
				return
			}
		}
		responseTransformationDuration := time.Since(responseTransformationStart)
		timingHeaderValue := fmt.Sprintf("postProxyTransformations;dur=%f", responseTransformationDuration.Seconds())
		w.Header().Add("Server-Timing", timingHeaderValue)
	}

	for headerKey, headerValues := range proxyResponse.Header {
		for _, headerValue := range headerValues {
			w.Header().Add(headerKey, headerValue)
		}
	}

	_, err = io.Copy(w, proxyResponse.Body)
	if err != nil {
		errorHandler <- err
		<-statusChannel
		return
	}
}
