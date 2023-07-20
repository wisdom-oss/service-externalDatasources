package databaseTypes

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Duration time.Duration

// Scan implements the driver.Scan interface used to read custom
// data types and type aliases from the database
func (d *Duration) Scan(src interface{}) error {
	var rowString string
	switch src.(type) {
	case []byte:
		rowString = string(src.([]byte))
	case string:
		rowString = src.(string)
	default:
		return errors.New("unsupported scan input")
	}
	durationString := fmt.Sprintf("%ss", rowString)
	dur, err := time.ParseDuration(durationString)
	if err != nil {
		return err
	}
	*d = Duration(dur)
	return nil
}

// MarshalJSON handles the conversion of the type alias into a json
// entity.
//
// In this particular case, the duration is converted into microseconds
// since it is the smallest duration possible to be stored by the postgres
// database used in the backend
func (d *Duration) MarshalJSON() ([]byte, error) {
	nativeDuration := time.Duration(*d)
	return json.Marshal(nativeDuration.Microseconds())
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	var aux int64
	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}
	durationString := fmt.Sprintf("%dus", aux)
	dur, err := time.ParseDuration(durationString)
	*d = Duration(dur)
	return nil
}

func (d Duration) Value() (driver.Value, error) {
	dur := time.Duration(d)
	years := int64(dur / (365 * 24 * time.Hour))
	months := int64((dur % (365 * 24 * time.Hour)) / (30 * 24 * time.Hour))
	days := int64((dur % (30 * 24 * time.Hour)) / (24 * time.Hour))
	hours := int64((dur % (24 * time.Hour)) / time.Hour)
	minutes := int64((dur % time.Hour) / time.Minute)
	seconds := int64((dur % time.Minute) / time.Second)
	microseconds := int64((dur % time.Second) / time.Millisecond)
	return fmt.Sprintf("%d years %d mons %d days %02d:%02d:%02d.%d\n", years, months, days, hours, minutes, seconds, microseconds), nil
}
