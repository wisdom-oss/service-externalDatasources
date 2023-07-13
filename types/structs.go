package types

import (
	"external-api-service/enums"
	"github.com/lib/pq"
)

type BaseData struct {
	ID          string  `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Description *string `json:"description" db:"description"`
}

type Metadata struct {
	// Reference contains information about the
	Reference *DataReference `json:"reference" db:"reference"`
	// Origin
	Origin *DataOrigin `json:"origin" db:"origin"`
	// DistinctiveFeatures
	DistinctiveFeatures          *Tuples                             `json:"distinctiveFeatures" db:"distinctive_features"`
	UsageRights                  *string                             `json:"usageRights" db:"usage_rights"`
	UsageDuties                  *string                             `json:"usageDuties" db:"usage_duties"`
	RealEntities                 pq.StringArray                      `json:"realEntities" db:"real_entities"`
	LocalExpert                  *Tuples                             `json:"localExpert" db:"local_expert"`
	Documentation                *Documentations                     `json:"externalDocumentation" db:"external_documentation"`
	UpdateRate                   *Duration                           `json:"updateRate" db:"update_rate"`
	Languages                    pq.StringArray                      `json:"languages" db:"languages"`
	Billing                      *PricingInformation                 `json:"pricingInformation" db:"billing"`
	Provision                    *DataProvisioning                   `json:"provision" db:"provisioning"`
	DerivedFrom                  *string                             `json:"derivedFrom" db:"derived_from"`
	Recent                       *bool                               `json:"isRecent" db:"is_recent"`
	Validity                     *Validity                           `json:"validity" db:"validity"`
	Duplicates                   *CheckedRange                       `json:"duplicates" db:"duplicates"`
	Errors                       *CheckedRange                       `json:"errors" db:"errors"`
	Precision                    *enums.PrecisionLevel               `json:"precision" db:"precision"`
	Reputation                   *enums.Reputation                   `json:"reputation" db:"reputation"`
	DataObjectivity              *DataObjectivity                    `json:"dataObjectivity" db:"objectivity"`
	UsualSurveyMethod            *bool                               `json:"usualSurveyMethod" db:"usual_survey_method"`
	Density                      *CheckedRange                       `json:"density" db:"density"`
	Coverage                     *enums.NoneHighRange                `json:"coverage" db:"coverage"`
	RepresentationConsistency    *enums.NoneHighRange                `json:"representationConsistency" db:"representation_consistency"`
	LogicalConsistency           *LogicalConsistency                 `json:"logicalConsistency" db:"logical_consistency"`
	DataDelay                    *DataDelay                          `json:"dataDelay" db:"delay"`
	DelayInformationTransmission *enums.DelayInformationTransmission `json:"delayInformationTransmission" db:"delay_information"`
	PerformanceLimitations       *enums.NoneHighRange                `json:"performanceLimitations" db:"performancelimitations"`
	Availability                 *enums.NoneHighRange                `json:"availability" db:"availability"`
	GDPRCompliant                *bool                               `json:"gdprCompliant" db:"gdpr_compliant"`
}

type API struct {
	IsSecure          *bool   `json:"isSecure" db:"is_secure"`
	Host              *string `json:"host" db:"host"`
	Port              *int    `json:"port" db:"port"`
	Path              *string `json:"path" db:"path"`
	AdditionalHeaders []Tuple `json:"additionalHeaders" `
}

type ExternalDataSource struct {
	BaseData
	Metadata Metadata `json:"metadata"`
	API      API      `json:"api"`
}
