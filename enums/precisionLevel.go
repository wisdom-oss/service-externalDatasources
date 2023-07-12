package enums

import (
	"errors"
	"regexp"
)

type PrecisionLevel string

const (
	PRECISION_LEVEL_FINE      = "fine"
	PRECISION_LEVEL_USUAL     = "usual"
	PRECISION_LEVEL_UNUSUAL   = "unusual"
	PRECISION_LEVEL_IMPRECISE = "imprecise"
)

func (pl PrecisionLevel) String() string {
	return string(pl)
}

func (pl *PrecisionLevel) Scan(src interface{}) error {
	var rowString string
	switch src.(type) {
	case []byte:
		rowString = string(src.([]byte))
	case string:
		rowString = src.(string)
	default:
		return errors.New("unsupported scan input")
	}
	regex := regexp.MustCompile(`^(fine|usual|unusual|imprecise)$`)
	matches := regex.FindStringSubmatch(rowString)
	if len(matches) != 2 {
		return errors.New("unsupported count of matches found")
	}
	*pl = PrecisionLevel(matches[1])
	return nil
}
