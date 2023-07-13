package types

import (
	"database/sql/driver"
	"errors"
	"external-api-service/enums"
	"fmt"
	"regexp"
	"strconv"
)

type CheckedRange struct {
	Checked bool                `json:"checked"`
	Range   enums.NoneHighRange `json:"range"`
}

func (cr *CheckedRange) Scan(src interface{}) error {
	var rowString string
	switch src.(type) {
	case []byte:
		rowString = string(src.([]byte))
	case string:
		rowString = src.(string)
	default:
		return errors.New("unsupported scan input")
	}
	regex := regexp.MustCompile(`^\(([tf]),(none|low|medium|high)\)$`)
	matches := regex.FindStringSubmatch(rowString)
	if len(matches) != 3 {
		return errors.New("unsupported count of matches found")
	}
	var err error
	cr.Checked, err = strconv.ParseBool(matches[1])
	if err != nil {
		return err
	}
	cr.Range = enums.NoneHighRange(matches[2])
	return nil
}

func (cr CheckedRange) Value() (driver.Value, error) {
	return fmt.Sprintf("(%b,%s)", cr.Checked, cr.Range), nil
}
