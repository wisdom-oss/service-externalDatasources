package types

// This file contains structs used on the http side of this microservice (e.g.,
// response layouts)

// WISdoMError is the default response layout for an error that occurred during
// request handling.
type WISdoMError struct {
	// ErrorCode contains a short error code identifying the error and the
	// service
	ErrorCode string `json:"code"`
	// ErrorTitle contains a human-readable error title
	ErrorTitle string `json:"title"`
	// ErrorDescription contains a human-readable description of the error that
	// occurred while handling the request
	ErrorDescription string `json:"description"`
	// HttpStatusCode contains the numeric http code that is associated with the
	// error. this value should be used to write out the error
	HttpStatusCode int `json:"httpCode"`
	// HttpStatusText contains the description of the HttpStatusCode to allow
	// humans to understand the numeric http code
	HttpStatusText string `json:"httpError"`
}
