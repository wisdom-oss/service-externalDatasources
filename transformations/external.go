package transformations

import (
	"encoding/json"
	"net/http"
)

type ExternalTransformation struct {
	ExecuteBeforeProxy bool        `json:"beforeProxy"`
	ExecuteAfterProxy  bool        `json:"afterProxy"`
	Executable         string      `json:"executable"`
	Configuration      interface{} `json:"configuration"`
}

func (et ExternalTransformation) ApplyBefore(r *http.Request, data []byte) error {
	// try to parse the configuration for the transformation
	if err := json.Unmarshal(data, &et); err != nil {
		return err
	}
	// check if the transformation shall be executed before proxying
	if !et.ExecuteBeforeProxy {
		return nil
	}
	// TODO: Discuss and implement external transformation execution possibilities
	return nil
}

func (et ExternalTransformation) ApplyAfter(r *http.Response, data []byte) error {
	// try to parse the configuration for the transformation
	if err := json.Unmarshal(data, &et); err != nil {
		return err
	}
	// check if the transformation shall be executed before proxying
	if !et.ExecuteBeforeProxy {
		return nil
	}
	// TODO: Discuss and implement external transformation execution
	// 	possibilities (only response body)
	return nil
}
