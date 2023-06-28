package types

import "github.com/jackc/pgtype"

type NoneHighRange string

const (
	RANGE_NONE   NoneHighRange = "none"
	RANGE_LOW    NoneHighRange = "low"
	RANGE_MEDIUM NoneHighRange = "medium"
	RANGE_HIGH   NoneHighRange = "high"
)

type NoneHighRangeType struct {
	Value  NoneHighRange
	Status pgtype.Status
	buf    []byte
}

func (nhrt *NoneHighRangeType) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	if src == nil {
		*nhrt = NoneHighRangeType{Status: pgtype.Null}
		return nil
	}
	var s pgtype.Varchar
	if err := s.DecodeBinary(ci, src); err != nil {
		return err
	}
	*nhrt = NoneHighRangeType{Value: NoneHighRange(s.String), Status: pgtype.Present}
	return nil
}

func (nhrt *NoneHighRangeType) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	switch nhrt.Status {
	case pgtype.Present:
		str := string(nhrt.Value)
		strVal := pgtype.Varchar{String: str, Status: pgtype.Present}
		return strVal.EncodeBinary(ci, buf)
	case pgtype.Null:
		strVal := pgtype.Varchar{Status: pgtype.Null}
		return strVal.EncodeBinary(ci, buf)
	default:
		return nil, nil
	}
}
