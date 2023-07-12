package types

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type PricingModel string

const (
	PRICING_MODEL_OPEN_SOURCE     PricingModel = "openSource"
	PRICING_MODEL_FREE            PricingModel = "free"
	PRICING_MODEL_SINGLE_PURCHASE PricingModel = "singlePurchase"
	PRICING_MODEL_PER_ACCESS      PricingModel = "perAccess"
	PRICING_MODEL_PER_TIME        PricingModel = "perTimeUnit"
	PRICING_MODEL_PER_DATA_AMOUNT PricingModel = "perDataAmount"
)

func (pm PricingModel) String() string {
	return string(pm)
}

var pricingModels = []string{
	PRICING_MODEL_OPEN_SOURCE.String(),
	PRICING_MODEL_FREE.String(),
	PRICING_MODEL_SINGLE_PURCHASE.String(),
	PRICING_MODEL_PER_ACCESS.String(),
	PRICING_MODEL_PER_TIME.String(),
	PRICING_MODEL_PER_DATA_AMOUNT.String(),
}

type PricingInformation struct {
	Model        PricingModel `json:"model"`
	PricePerUnit float64      `json:"pricePerUnit"`
}

func (pi *PricingInformation) Scan(src interface{}) error {
	var rowString string
	switch src.(type) {
	case []byte:
		rowString = string(src.([]byte))
	case string:
		rowString = src.(string)
	default:
		return errors.New("unsupported scan input")
	}
	// create a regular expression matching the needed two values
	pricingModelGroup := strings.Join(pricingModels, "|")
	regexString := `^\((` + pricingModelGroup + `),([0-9]+[.]?[0-9]*)\)$`
	regex := regexp.MustCompile(regexString)
	matches := regex.FindStringSubmatch(rowString)
	if len(matches) != 3 {
		return errors.New("unexpected match count")
	}
	values := matches[1:]
	pi.Model = PricingModel(values[0])
	cost, err := strconv.ParseFloat(values[1], 64)
	if err != nil {
		return err
	}
	pi.PricePerUnit = cost
	return nil
}
