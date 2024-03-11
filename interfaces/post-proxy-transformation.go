package interfaces

import "net/http"

type PostProxyTransformation interface {
	ApplyAfter(r *http.Response, options map[string]interface{}) error
}
