package enums

import (
	"errors"
	"regexp"
)

type DelayInformationTransmission string

func (dit DelayInformationTransmission) String() string {
	return string(dit)
}

func (dit *DelayInformationTransmission) Scan(src interface{}) error {
	var rowString string
	switch src.(type) {
	case []byte:
		rowString = string(src.([]byte))
	case string:
		rowString = src.(string)
	default:
		return errors.New("unsupported scan input")
	}
	regex := regexp.MustCompile(`^(direct|automatic|manual|none)$`)
	matches := regex.FindStringSubmatch(rowString)
	if len(matches) != 2 {
		return errors.New("unsupported count of matches found")
	}
	*dit = DelayInformationTransmission(matches[1])
	return nil
}
