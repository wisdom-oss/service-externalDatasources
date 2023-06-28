package types

import (
	"fmt"
	"github.com/jackc/pgtype"
)

type AuthorizationDataLocation string

const (
	LOCATION_HEADER AuthorizationDataLocation = "header"
	LOCATION_QUERY  AuthorizationDataLocation = "query"
)

func (loc AuthorizationDataLocation) ToString() string {
	return fmt.Sprintf("%s", loc)
}

type AuthorizationInformation struct {
	Location AuthorizationDataLocation
	Key      string
	Value    string
}

// DecodeBinary implements a function to allow the decoding of the AuthorizationInformation type in database situations
func (ai AuthorizationInformation) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	if src == nil {
		return ErrNilDecoding
	}
	if err := (pgtype.CompositeFields{&ai.Location, &ai.Key, &ai.Value}).DecodeBinary(ci, src); err != nil {
		return err
	}
	return nil
}

// EncodeBinary implements a function to allow the encoding of the AuthorizationInformation type in database situations
func (ai AuthorizationInformation) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) (newBuf []byte, err error) {
	loc := pgtype.Text{
		String: ai.Location.ToString(),
		Status: pgtype.Present,
	}
	key := pgtype.Text{
		String: ai.Key,
		Status: pgtype.Present,
	}
	val := pgtype.Text{
		String: ai.Value,
		Status: pgtype.Present,
	}
	return (pgtype.CompositeFields{&loc, &key, &val}).EncodeBinary(ci, buf)
}
