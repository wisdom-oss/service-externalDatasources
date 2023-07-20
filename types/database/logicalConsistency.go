package databaseTypes

import (
	"database/sql/driver"
	"errors"
	"external-api-service/enums"
	"fmt"
	"regexp"
	"strconv"
)

type LogicalConsistency struct {
	Checked                  bool                `json:"checked"`
	ContradictionsExaminable bool                `json:"contradictionsExaminable"`
	Range                    enums.NoneHighRange `json:"range"`
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
	lc.Range = enums.NoneHighRange(matches[3])
	return nil
}

func (lc LogicalConsistency) Value() (driver.Value, error) {
	return fmt.Sprintf("(%t,%t,%s)", lc.Checked, lc.ContradictionsExaminable, lc.Range), nil
}
