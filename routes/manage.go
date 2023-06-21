package routes

import (
	"encoding/json"
	"external-api-service/globals"
	"external-api-service/types"
	wisdomType "github.com/wisdom-oss/commonTypes"
	"net/http"
	"net/url"
)

// NewExternalAPI allows the creation of a new forward-proxy for an external api
func NewExternalAPI(w http.ResponseWriter, r *http.Request) {
	// since the data is expected as form-data parse the form to start handling it
	err := r.ParseMultipartForm(8589934592)
	if err != nil {
		// handle the error that occurred
		e := wisdomType.WISdoMError{}
		e.WrapError(err, globals.ServiceName)
		_ = e.Send(w)
		return
	}

	// now check if the name is set
	nameList, isNameSet := r.MultipartForm.Value["name"]
	if !isNameSet {
		_ = globals.Errors["MISSING_NAME"].Send(w)
		return
	}

	// now check if the base uri is set
	baseUriList, isBaseUriSet := r.MultipartForm.Value["baseUri"]
	if !isBaseUriSet {
		_ = globals.Errors["MISSING_BASE_URI"].Send(w)
		return
	}

	// now use only the first entries of the baseUri and the names
	name := nameList[0]
	baseUriStr := baseUriList[0]

	// now parse the base uri
	baseUri, err := url.ParseRequestURI(baseUriStr)
	if err != nil {
		_ = globals.Errors["INVALID_BASE_URI"].Send(w)
		return
	}

	// now check if the protocol is either http or https
	if baseUri.Scheme != "http" && baseUri.Scheme != "https" {
		_ = globals.Errors["INVALID_URI_SCHEME"].Send(w)
		return
	}

	// since the uri and all other required variables are present, now get the additional headers

	// fixme: remove this
	api := types.ExternalAPI{
		Name:    name,
		BaseURI: baseUriStr,
	}
	_ = json.NewEncoder(w).Encode(api)
}
