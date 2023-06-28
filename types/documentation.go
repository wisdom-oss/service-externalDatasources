package types

import "github.com/jackc/pgtype"

type Documentation struct {
	Type      string
	Location  []Tuple
	Verbosity NoneHighRange
	Status    pgtype.Status
}

func (doc *Documentation) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	if src == nil {
		*doc = Documentation{Status: pgtype.Null}
		return nil
	}

	var fields pgtype.CompositeFields
	fields = append(
		fields,
		&doc.Type,
		&doc.Location,
		&doc.Verbosity,
	)
	if err := fields.DecodeBinary(ci, src); err != nil {
		*doc = Documentation{Status: pgtype.Null}
		return err
	}
	doc.Status = pgtype.Present
	return nil
}

func (doc *Documentation) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) (newBuf []byte, err error) {
	fields := pgtype.CompositeFields{
		&doc.Type,
		&doc.Location,
		&doc.Verbosity,
	}
	return fields.EncodeBinary(ci, buf)
}
