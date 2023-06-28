package types

import (
	"github.com/jackc/pgtype"
)

type Header struct {
	Key    string
	Value  string
	Status pgtype.Status
}

// DecodeBinary is a method on Header to decode a pgtype binary to a Header.
func (hdr *Header) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	if src == nil {
		*hdr = Header{Status: pgtype.Null}
		return nil
	}

	var t pgtype.Text
	if err := t.DecodeBinary(ci, src); err != nil {
		return err
	}

	hdr.Key = t.String
	hdr.Value = t.String
	hdr.Status = pgtype.Present

	return nil
}

// EncodeBinary is a method on Header to encode a Header to a pgtype binary.
func (hdr *Header) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	if hdr.Status == pgtype.Null {
		return nil, nil
	}

	t := pgtype.Text{String: hdr.Key, Status: pgtype.Present}
	t2 := pgtype.Text{String: hdr.Value, Status: pgtype.Present}
	buf, err := t.EncodeBinary(ci, buf)
	if err != nil {
		return nil, err
	}

	return t2.EncodeBinary(ci, buf)
}

// DecodeText is a method on Header to decode a string to a Header.
func (hdr *Header) DecodeText(ci *pgtype.ConnInfo, src string) error {
	if src == "" {
		*hdr = Header{Status: pgtype.Null}
		return nil
	}

	hdr.Key = src
	hdr.Value = src
	hdr.Status = pgtype.Present

	return nil
}

// EncodeText is a method on Header to encode a Header to a string.
func (hdr *Header) EncodeText(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	if hdr.Status == pgtype.Null {
		return nil, nil
	}

	t := pgtype.Text{String: hdr.Key, Status: pgtype.Present}
	t2 := pgtype.Text{String: hdr.Value, Status: pgtype.Present}
	buf, err := t.EncodeText(ci, buf)
	if err != nil {
		return nil, err
	}

	return t2.EncodeText(ci, buf)
}
