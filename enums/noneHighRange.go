package enums

import (
	"errors"
	"regexp"
)

type NoneHighRange string

const (
	RANGE_NONE   NoneHighRange = "none"
	RANGE_LOW    NoneHighRange = "low"
	RANGE_MEDIUM NoneHighRange = "medium"
	RANGE_HIGH   NoneHighRange = "high"
)

func (nhr NoneHighRange) String() string {
	return string(nhr)
}

func (nhr *NoneHighRange) Scan(src interface{}) error {
	var rowString string
	switch src.(type) {
	case []byte:
		rowString = string(src.([]byte))
	case string:
		rowString = src.(string)
	default:
		return errors.New("unsupported scan input")
	}
	regex := regexp.MustCompile(`^(none|low|medium|high)$`)
	matches := regex.FindStringSubmatch(rowString)
	if len(matches) != 2 {
		return errors.New("unexpected match count")
	}
	*nhr = NoneHighRange(matches[1])
	return nil
}
