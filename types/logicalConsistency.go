package types

import "github.com/jackc/pgtype"

type LogicalConsistency struct {
	Checked                  bool
	ContradictionsExaminable bool
	Range                    NoneHighRange
	Status                   pgtype.Status
}

// DecodeBinary implements a function to allow the decoding of the CheckedRange type in database situations
func (lc *LogicalConsistency) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	if src == nil {
		*lc = LogicalConsistency{}
		return nil
	}
	if err := (pgtype.CompositeFields{&lc.Checked, &lc.ContradictionsExaminable, &lc.Status}).DecodeBinary(ci, src); err != nil {
		return err
	}
	return nil
}

// EncodeBinary implements a function to allow the encoding of the AuthorizationInformation type in database situations
func (lc *LogicalConsistency) EncodeBinary(ci *pgtype.ConnInfo, src []byte) (newBuf []byte, err error) {
	checked := pgtype.Bool{
		Bool:   lc.Checked,
		Status: pgtype.Present,
	}
	examinable := pgtype.Bool{
		Bool:   lc.ContradictionsExaminable,
		Status: pgtype.Present,
	}
	return (pgtype.CompositeFields{checked, examinable, lc.Range}).EncodeBinary(ci, src)
}
