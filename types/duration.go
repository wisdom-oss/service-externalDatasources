package types

import (
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
