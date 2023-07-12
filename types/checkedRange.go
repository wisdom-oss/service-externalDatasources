package types

import (
	"errors"
	"regexp"
	"strconv"
)

type CheckedRange struct {
	Checked bool
	Range   NoneHighRange
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
	cr.Range = NoneHighRange(matches[2])
	return nil
}
