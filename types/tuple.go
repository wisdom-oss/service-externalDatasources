package types

import (
	"encoding/json"
	"errors"
	"github.com/jackc/pgtype"
)

// Tuple allows associating two strings to another without requiring a unique
// key.
//
// The Tuple is marshalled into an array containing X as the first value and
// Y as the second value. When unmarshalling json into this type, only the first
// two values will be used, the other ones will be discarded
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

// DecodeBinary handles the conversion of the database entry into the tuple
func (t *Tuple) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	if src == nil {
		return errors.New("NULL values cannot be decoded. Scan into a &*Tuple to handle NULLs")
	}
	if err := (pgtype.CompositeFields{&t.X, &t.Y}).DecodeBinary(ci, src); err != nil {
		return err
	}
	return nil
}

// EncodeBinary handles the conversion of the struct into a database entry
func (t Tuple) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) (newBuy []byte, err error) {
	left := pgtype.Text{String: t.X, Status: pgtype.Present}
	right := pgtype.Text{String: t.Y, Status: pgtype.Present}
	return (pgtype.CompositeFields{&left, &right}).EncodeBinary(ci, buf)
}
