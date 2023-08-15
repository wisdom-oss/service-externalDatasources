package interfaces

import "net/http"

type TransformationFunction interface {
	ApplyBefore(r *http.Request, data []byte) error
	ApplyAfter(r *http.Response, data []byte) error
}
