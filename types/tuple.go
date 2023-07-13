package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Tuple allows associating two strings to another without requiring a unique
// key.
//
// The Tuple is marshaled into an array containing X as the first value and
// Y as the second value.
// While unmarshalling json into this type, only the first two values will be
// used, the other ones will be discarded
type Tuple struct {
	// The left-hand value
	X string
	// The right-hand value
	Y string
}

// Tuples is a type alias for []Tuple since a direct array is not allowed in
// a struct scanned by the blockloop/scan package
type Tuples []Tuple

// MarshalJSON handles the custom conversion of the Tuple to a json entity
func (t *Tuple) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string{t.X, t.Y})
}

// UnmarshalJSON handles the custom conversion of json into the Tuple
func (t *Tuple) UnmarshalJSON(data []byte) error {
	var aux []string
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if len(aux) < 2 {
		return errors.New("tuple requires at least 2 entries in array")
	}
	t.X = aux[0]
	t.Y = aux[1]
	return nil
}

func (t *Tuple) Scan(src interface{}) error {
	var rowString string
	switch src.(type) {
	case []byte:
		rowString = string(src.([]byte))
	case string:
		rowString = src.(string)
	default:
		return errors.New("unsupported scan input")
	}
	var tuple Tuple
	rowString = strings.ReplaceAll(rowString, ",", " ")
	rowString = strings.Trim(rowString, `{"()}`)
	_, err := fmt.Sscanf(strings.TrimSpace(rowString), "%s %s", &tuple.X, &tuple.Y)
	if err != nil {
		return err
	}
	*t = tuple
	return nil
}

func (t *Tuple) Value() (driver.Value, error) {
	return fmt.Sprintf("\"(%s,%s)\"", t.X, t.Y), nil
}

// Scan implements the db.Scanner interface on Tuples
func (t *Tuples) Scan(value interface{}) error {
	rowString := string(value.([]byte))
	entries := strings.Split(rowString, `","`)
	var tuples Tuples
	for _, entry := range entries {
		var tuple Tuple
		if err := tuple.Scan(entry); err != nil {
			return err
		}
		tuples = append(tuples, tuple)
	}
	*t = tuples
	return nil
}

func (t Tuples) Value() (driver.Value, error) {
	var singleTuples []string

	for _, tuple := range t {
		val, err := tuple.Value()
		if err != nil {
			return nil, err
		}
		singleTuples = append(singleTuples, val.(string))
	}
	tupleArray := strings.Join(singleTuples, ",")
	tupleArray = fmt.Sprintf("{%s}", tupleArray)
	return tupleArray, nil
}
