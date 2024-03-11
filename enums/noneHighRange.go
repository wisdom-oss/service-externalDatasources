package enums

import (
	"encoding/json"
	"fmt"
)

type NoneHighRange int8

const (
	RangeNone NoneHighRange = iota
	RangeLow
	RangeMedium
	RangeHigh
)

func (r NoneHighRange) String() string {
	switch r {
	case RangeNone:
		return "none"
	case RangeLow:
		return "low"
	case RangeMedium:
		return "medium"
	case RangeHigh:
		return "high"
	default:
		return ""
	}
}

func (r NoneHighRange) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%s", r))
}

func (r *NoneHighRange) UnmarshalJSON(src []byte) error {
	var str string
	err := json.Unmarshal(src, &str)
	if err != nil {
		return err
	}
	switch str {
	case "none":
		*r = RangeNone
	case "low":
		*r = RangeLow
	case "medium":
		*r = RangeMedium
	case "high":
		*r = RangeHigh
	default:
		*r = NoneHighRange(-1)
	}
	return nil
}
