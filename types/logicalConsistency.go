package types

import (
	"errors"
	"regexp"
	"strconv"
)

type LogicalConsistency struct {
	Checked                  bool
	ContradictionsExaminable bool
	Range                    NoneHighRange
}

func (lc *LogicalConsistency) Scan(src interface{}) error {
	var rowString string
	switch src.(type) {
	case []byte:
		rowString = string(src.([]byte))
	case string:
		rowString = src.(string)
	default:
		return errors.New("unsupported scan input")
	}
	regex := regexp.MustCompile(`^\(([tf]),([tf]),(none|low|medium|high)\)$`)
	matches := regex.FindStringSubmatch(rowString)
	if len(matches) != 4 {
		return errors.New("unexpected count of matches in regex")
	}
	var err error
	lc.Checked, err = strconv.ParseBool(matches[1])
	if err != nil {
		return err
	}
	lc.ContradictionsExaminable, err = strconv.ParseBool(matches[2])
	return nil
}
