package types

import "github.com/jackc/pgtype"

type PricingModel string

const (
	PRICING_MODEL_OPEN_SOURCE     PricingModel = "openSource"
	PRICING_MODEL_FREE            PricingModel = "free"
	PRICING_MODEL_SINGLE_PURCHASE PricingModel = "singlePurchase"
	PRICING_MODEL_PER_ACCESS      PricingModel = "perAccess"
	PRICING_MODEL_PER_TIME        PricingModel = "perTimeUnit"
	PRICING_MODEL_PER_DATA_AMOUNT PricingModel = "perDataAmount"
)

type PricingModelType struct {
	Value  PricingModel
	Status pgtype.Status
	buf    []byte
}

func (pmt *PricingModelType) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	if src == nil {
		*pmt = PricingModelType{Status: pgtype.Null}
		return nil
	}
	var s pgtype.Varchar
	if err := s.DecodeBinary(ci, src); err != nil {
		return err
	}
	*pmt = PricingModelType{Value: PricingModel(s.String), Status: pgtype.Present}
	return nil
}

func (pmt *PricingModelType) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	switch pmt.Status {
	case pgtype.Present:
		str := string(pmt.Value)
		strVal := pgtype.Varchar{String: str, Status: pgtype.Present}
		return strVal.EncodeBinary(ci, buf)
	case pgtype.Null:
		strVal := pgtype.Varchar{Status: pgtype.Null}
		return strVal.EncodeBinary(ci, buf)
	default:
		return nil, nil
	}
}

type PricingInformation struct {
	Model        PricingModel
	PricePerUnit float64
}

func (pi *PricingInformation) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	if src == nil {
		*pi = PricingInformation{}
	}
	return (pgtype.CompositeFields{&pi.Model, &pi.PricePerUnit}).DecodeBinary(ci, src)
}

func (pi *PricingInformation) EncodeBinary(ci *pgtype.ConnInfo, src []byte) (newBuf []byte, err error) {
	return (pgtype.CompositeFields{&pi.Model, &pi.PricePerUnit}).EncodeBinary(ci, src)
}
