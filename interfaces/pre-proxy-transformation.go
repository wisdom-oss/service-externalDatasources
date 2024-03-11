package interfaces

import "net/http"

type PreProxyTransformation interface {
	ApplyBefore(r *http.Request, options map[string]interface{}) error
}
