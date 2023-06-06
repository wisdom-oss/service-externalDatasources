package types

import (
	"encoding/json"
	"fmt"
	"net/http"

	"external-api-service/globals"
)

// this file contains functions extending some types from the http section

// WrapError allows using a native golang error to be wrapped into a WISdoMError
// to allow the population of an internal error.
// This function always uses HTTP error 500 to indicate that the issue was
// internal and not the client's fault.
func (e *WISdoMError) WrapError(err error) {
	// get the error message from the error
	msg := err.Error()

	// now build set the error message to the WISdoMError
	e.ErrorDescription = fmt.Sprintf("an internal error occurred while processing your request: %s", msg)

	// now set the other parts of the WISdoMError
	e.ErrorCode = fmt.Sprintf("%s.%s", globals.ServiceName, "INTERNAL_ERROR")
	e.ErrorTitle = "Internal Error in Microservice"
	e.HttpStatusCode = http.StatusInternalServerError
	e.HttpStatusText = http.StatusText(e.HttpStatusCode)
}

// SendError allows the direct sending of the WISdoMError as a response. If an
// error occurs while sending the error, it will be returned.
func (e *WISdoMError) SendError(w http.ResponseWriter) error {
	// start the process by setting the content type of the response to
	// `text/json`.
	w.Header().Set("Content-Type", "text/json")
	// then start the sending by sending the response code first
	w.WriteHeader(e.HttpStatusCode)
	// then start converting the message and print it to the response using the
	// json format and return the error if one occurred or nil otherwise.
	return json.NewEncoder(w).Encode(e)
}
