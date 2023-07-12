package types

import (
	"errors"
	"regexp"
)

type DataDelay struct {
	Ingress NoneHighRange `json:"ingress"`
	Egress  NoneHighRange `json:"egress"`
}

func (dd *DataDelay) Scan(src interface{}) error {
	var rowString string
	switch src.(type) {
	case []byte:
		rowString = string(src.([]byte))
	case string:
		rowString = src.(string)
	default:
		return errors.New("unsupported scan input")
	}
	regex := regexp.MustCompile(`^\((none|low|medium|high),(none|low|medium|high)\)$`)
	matches := regex.FindStringSubmatch(rowString)
	if len(matches) != 3 {
		return errors.New("unexpected match count")
	}
	dd.Ingress = NoneHighRange(matches[1])
	dd.Egress = NoneHighRange(matches[2])
	return nil
}
