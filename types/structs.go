package types

type BaseData struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type Metadata struct {
	ID                  string   `json:"-" db:"id"`
	DistinctiveFeatures []Tuple  `json:"distinctiveFeatures" db:"distinctive_features"`
	UsageRights         *string  `json:"usageRights" db:"usage_rights"`
	UsageDuties         *string  `json:"usageDuties" db:"usage_duties"`
	RealEntities        []string `json:"realEntities" db:"real_entities"`
}

type API struct {
}

type ExternalDataSource struct {
	BaseData
	Metadata
	API
}
