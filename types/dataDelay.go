package types

import "github.com/jackc/pgtype"

type DataDelay struct {
	Ingress NoneHighRange
	Egress  NoneHighRange
	Status  pgtype.Status
}

// DecodeBinary implements a function to allow the decoding of the DataDelay type in database situations
func (dd *DataDelay) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	if src == nil {
		*dd = DataDelay{
			Status: pgtype.Null,
		}
		return nil
	}
	return (pgtype.CompositeFields{&dd.Ingress, &dd.Egress}).DecodeBinary(ci, src)
}

// EncodeBinary implements a function to allow the encoding of the DataDelay type in database situations
func (dd *DataDelay) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) (newBuf []byte, err error) {
	fields := pgtype.CompositeFields{
		&dd.Ingress,
		&dd.Egress,
	}
	return fields.EncodeBinary(ci, buf)
}
