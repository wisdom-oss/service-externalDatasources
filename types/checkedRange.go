package types

import "github.com/jackc/pgtype"

type CheckedRange struct {
	Checked bool
	Range   NoneHighRangeType
}

// DecodeBinary implements a function to allow the decoding of the CheckedRange type in database situations
func (cr *CheckedRange) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	if src == nil {
		*cr = CheckedRange{}
		return nil
	}
	if err := (pgtype.CompositeFields{&cr.Checked, &cr.Range}).DecodeBinary(ci, src); err != nil {
		return err
	}
	return nil
}

// EncodeBinary implements a function to allow the encoding of the AuthorizationInformation type in database situations
func (cr *CheckedRange) EncodeBinary(ci *pgtype.ConnInfo, src []byte) (newBuf []byte, err error) {
	checked := pgtype.Bool{
		Bool:   cr.Checked,
		Status: pgtype.Present,
	}
	return (pgtype.CompositeFields{&checked, &cr.Range}).EncodeBinary(ci, src)
}
