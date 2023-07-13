package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"time"
)

type DataReference struct {
	Topic             string       `json:"topic"`
	LocalReference    string       `json:"localReference"`
	TemporalReference [2]time.Time `json:"temporalReference"`
}

func (dr *DataReference) Scan(src interface{}) error {
	var rowString string
	switch src.(type) {
	case []byte:
		rowString = string(src.([]byte))
	case string:
		rowString = src.(string)
	default:
		return errors.New("unsupported scan input")
	}
	regex := regexp.MustCompile(`^\(([^[:cntrl:],]+),([^[:cntrl:],]+),"\[([0-9]{4}-[0-9]{2}-[0-9]{2}),([0-9]{4}-[0-9]{2}-[0-9]{2})\)"\)$`)
	matches := regex.FindStringSubmatch(rowString)
	if len(matches) != 5 {
		return errors.New("unexpected count of matches")
	}
	var obj DataReference
	obj.Topic = matches[1]
	obj.LocalReference = matches[2]
	startTime, err := time.Parse("2006-01-02", matches[3])
	if err != nil {
		return err
	}
	endTime, err := time.Parse("2006-01-02", matches[4])
	if err != nil {
		return err
	}
	obj.TemporalReference = [2]time.Time{startTime, endTime}
	*dr = obj
	return nil
}

func (dr DataReference) Value() (driver.Value, error) {
	// build the time strings
	rangeStart := dr.TemporalReference[0].Format("2006-01-02")
	rangeEnd := dr.TemporalReference[1].Format("2006-01-02")
	dateRange := fmt.Sprintf("[%s,%s)", rangeStart, rangeEnd)
	val := fmt.Sprintf("(%s,%s,\"%s\")", dr.Topic, dr.LocalReference, dateRange)
	return val, nil
}
