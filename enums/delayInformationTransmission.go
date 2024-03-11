package enums

import (
	"encoding/json"
	"fmt"
)

type DelayInformationTransmission int8

const DelayTransmissionUnknown DelayInformationTransmission = -1

const (
	DelayTransmissionNone DelayInformationTransmission = iota
	DelayTransmissionManual
	DelayTransmissionAutomatic
	DelayTransmissionDirect
)

func (t DelayInformationTransmission) String() string {
	switch t {
	case DelayTransmissionNone:
		return "none"
	case DelayTransmissionManual:
		return "manual"
	case DelayTransmissionAutomatic:
		return "automatic"
	case DelayTransmissionDirect:
		return "direct"
	default:
		return ""
	}
}

func (t DelayInformationTransmission) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%s", t))
}

func (t *DelayInformationTransmission) UnmarshalJSON(src []byte) error {
	var str string
	err := json.Unmarshal(src, &str)
	if err != nil {
		return err
	}
	switch str {
	case "none":
		*t = DelayTransmissionNone
	case "manual":
		*t = DelayTransmissionManual
	case "automatic":
		*t = DelayTransmissionAutomatic
	case "direct":
		*t = DelayTransmissionDirect
	default:
		*t = DelayTransmissionUnknown
	}
	return nil
}
