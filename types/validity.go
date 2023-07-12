package types

type Validity string

const (
	VALIDITY_FULLY   = "fully"
	VALIDITY_PARTIAL = "partially"
	VALIDITY_NONE    = "none"
)

func (v Validity) String() string {
	return string(v)
}
