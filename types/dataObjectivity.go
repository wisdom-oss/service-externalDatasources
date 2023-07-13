package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type DataObjectivity struct {
	ConflictOfInterest bool `json:"conflictOfInterest"`
	RawData            bool `json:"rawData"`
	AutomaticCapture   bool `json:"automaticCapture"`
}

func (do *DataObjectivity) Scan(src interface{}) error {
	var rowString string
	switch src.(type) {
	case []byte:
		rowString = string(src.([]byte))
	case string:
		rowString = src.(string)
	default:
		return errors.New("unsupported scan input")
	}
	regex := regexp.MustCompile(`^\(([tf]),([tf]),([tf])\)$`)
	matches := regex.FindStringSubmatch(rowString)
	if len(matches) < 4 {
		return errors.New("not enough entries in data type")
	}
	var obj DataObjectivity
	var err error
	obj.ConflictOfInterest, err = strconv.ParseBool(matches[1])
	if err != nil {
		return err
	}
	obj.RawData, err = strconv.ParseBool(matches[2])
	if err != nil {
		return err
	}
	obj.AutomaticCapture, err = strconv.ParseBool(matches[3])
	if err != nil {
		return err
	}
	*do = obj
	return nil
}

func (do DataObjectivity) Value() (driver.Value, error) {
	return fmt.Sprintf("(%t,%t,%t)", do.ConflictOfInterest, do.RawData, do.AutomaticCapture), nil
}
