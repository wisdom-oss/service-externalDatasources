package enums

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Reputation uint16

const (
	Independent Reputation = 1 << iota
	External
	SuspectedLow
	SuspectedHigh
)

func (r Reputation) String() (reputation string) {
	if Independent&r != 0 {
		reputation += fmt.Sprintf("%s,", "independent")
	}
	if External&r != 0 {
		reputation += fmt.Sprintf("%s,", "external")
	}
	if SuspectedLow&r != 0 {
		reputation += fmt.Sprintf("%s,", "suspectedLow")
	}
	if SuspectedHigh&r != 0 {
		reputation += fmt.Sprintf("%s,", "suspectedHigh")
	}
	return strings.Trim(reputation, `,`)
}

func (r Reputation) MarshalJSON() ([]byte, error) {
	reputation := r.String()
	return json.Marshal(strings.Split(reputation, `,`))
}

func (r *Reputation) UnmarshalJSON(src []byte) error {
	var reputationSlice []string
	err := json.Unmarshal(src, &reputationSlice)
	if err != nil {
		return err
	}
	var reputation Reputation
	for _, rep := range reputationSlice {
		switch rep {
		case "independent":
			reputation |= Independent
		case "external":
			reputation |= External
		case "suspectedLow":
			reputation |= SuspectedLow
		case "suspectedHigh":
			reputation |= SuspectedHigh
		}
	}
	*r = reputation
	return nil
}
