package enums

import (
	"errors"
	"regexp"
)

type Reputation string

const (
	REPUTATION_INDEPENDENT_AND_EXTERNAL = "independent_and_external"
	REPUTATION_INDEPENDENT_OR_EXTERNAL  = "independent_or_external"
	REPUTATION_SUSPECTED_HIGH           = "suspected_high"
	REPUTATION_SUSPECTED_LOW            = "suspected_low"
)

func (r Reputation) String() string {
	return string(r)
}

func (r *Reputation) Scan(src interface{}) error {
	var rowString string
	switch src.(type) {
	case []byte:
		rowString = string(src.([]byte))
	case string:
		rowString = src.(string)
	default:
		return errors.New("unsupported scan input")
	}
	regex := regexp.MustCompile(`^(independent_and_external|independent_or_external|suspected_high|suspected_low)$`)
	matches := regex.FindStringSubmatch(rowString)
	if len(matches) != 2 {
		return errors.New("unsupported count of matches found")
	}
	*r = Reputation(matches[1])
	return nil
}
