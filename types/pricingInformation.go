package types

type PricingModel string

const (
	PRICING_MODEL_OPEN_SOURCE     PricingModel = "openSource"
	PRICING_MODEL_FREE            PricingModel = "free"
	PRICING_MODEL_SINGLE_PURCHASE PricingModel = "singlePurchase"
	PRICING_MODEL_PER_ACCESS      PricingModel = "perAccess"
	PRICING_MODEL_PER_TIME        PricingModel = "perTimeUnit"
	PRICING_MODEL_PER_DATA_AMOUNT PricingModel = "perDataAmount"
)

type PricingInformation struct {
	Model        PricingModel
	PricePerUnit float64
}

// TODO: implement conversion/parsing functions
