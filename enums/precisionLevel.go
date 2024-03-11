package enums

// PrecisionLevel is a custom type defined as an unsigned 8-bit integer.
// It is used to represent the precision level of data in a data source.
type PrecisionLevel uint8

const (
	PrecisionImprecise PrecisionLevel = iota + 1
	PrecisionUnusual
	PrecisionUsual
	PrecisionFine
)

func (p PrecisionLevel) String() string {
	switch p {
	case PrecisionImprecise:
		return "imprecise"
	case PrecisionUnusual:
		return "unusual"
	case PrecisionUsual:
		return "usual"
	case PrecisionFine:
		return "fine"
	default:
		return ""
	}
}

func (p *PrecisionLevel) Parse(src interface{}) {
	var precisionString string
	switch src.(type) {
	case string:
		precisionString = src.(string)
		break
	case []byte:
		precisionString = string(src.([]byte))
		break
	}

	switch precisionString {
	case "imprecise":
		*p = PrecisionImprecise
	case "unusual":
		*p = PrecisionUnusual
	case "usual":
		*p = PrecisionUsual
	case "fine":
		*p = PrecisionFine
	default:
		*p = 0
	}
}
