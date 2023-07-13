package types

import (
	"database/sql/driver"
	"errors"
	"external-api-service/enums"
	"fmt"
	"regexp"
)

type DataDelay struct {
	Ingress enums.NoneHighRange `json:"ingress"`
	Egress  enums.NoneHighRange `json:"egress"`
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
	dd.Ingress = enums.NoneHighRange(matches[1])
	dd.Egress = enums.NoneHighRange(matches[2])
	return nil
}

func (dd DataDelay) Value() (driver.Value, error) {
	return fmt.Sprintf("(%s,%s)", dd.Ingress, dd.Egress), nil
}
