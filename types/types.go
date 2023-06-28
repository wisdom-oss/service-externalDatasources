package types

import "errors"

var ErrNilDecoding = errors.New("NULL values cannot be decoded. scan into a &* of the type to handle NULL values")
