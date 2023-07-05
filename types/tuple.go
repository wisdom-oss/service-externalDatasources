package types

import (
	"encoding/json"
	"errors"
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
