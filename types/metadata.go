package types

import (
	"external-api-service/enums"
	dbType "external-api-service/types/database"
	"github.com/lib/pq"
)

// Metadata represents the metadata stored in the database for a single
// datasource.
// All fields are optional since the metadata of a datasource is
// also optional. Furthermore, it does contain the ID of the datasource.
// However, the ID will not be exported when marshaling it into json, due to
// the encoding/json package not supporting having the same tag used twice and
// ignoring them completely.
type Metadata struct {
	// DatasourceID contains the uuid of the datasource the metadata is
	//associated with
	DatasourceID string `json:"-" db:"id"`

	// Reference contains the topic of the data available in the datasource
	// and if available
	Reference *dbType.DataReference `json:"reference" db:"reference"`

	// DataOrigin contains information about the data provider, data creator
	// and data owner
	DataOrigin *dbType.DataOrigin `json:"dataOrigin" db:"origin"`

	// DistinctiveFeatures contains important properties and restrictions of
	// the datasource
	DistinctiveFeatures *dbType.Tuples `json:"distinctiveFeatures" db:"distinctive_features"`

	// UsageRights contains the rights that have been awarded for the datasource
	// and the usage of the contained data if the UsageDuties are followed
	UsageRights *string `json:"usageRights" db:"usage_rights"`

	// UsageDuties contains the duties that have to be fulfilled to be able to
	// use the datasource according to the UsageRights
	UsageDuties *string `json:"usageDuties" db:"usage_duties"`

	// Entities is an array of strings listing entities in the real world
	// represented by the data source
	Entities pq.StringArray `json:"entities" db:"real_entities"`

	// Expert contains [dbType.Tuples], which represent contact
	// information for the content-related/technical contact person
	Expert *dbType.Tuples `json:"expert" db:"local_expert"`

	// Documentation is an array consisting of [dbType.Documentation]
	// entries hidden by the used type alias.
	// It lists the available documentation about the data source
	Documentation *dbType.Documentations `json:"documentation" db:"external_documentation"`

	// UpdateRate contains a [time.Duration] hidden behind
	// [dbType.Duration] showing how often the data source is updated.
	// If this value is nil, the datasource is either manually updated or
	// updates are event-triggered
	UpdateRate *dbType.Duration `json:"updateRate" db:"update_rate"`

	// Languages contains ISO 639-1 language codes representing the languages
	// used in the datasource
	Languages pq.StringArray `json:"languages" db:"languages"`

	// Billing contains information about how the data source bills the access
	// to its data.
	Billing *dbType.PricingInformation `json:"billing" db:"billing"`

	// DataProvision contains information about the type of datasource and the
	// content type of the data source.
	DataProvision *dbType.DataProvisioning `json:"dataProvision" db:"provisioning"`

	// DerivedFrom contains the uuid of a datasource this datasource was derived
	// from
	DerivedFrom *string `json:"derivedFrom" db:"derived_from"`

	// Recent indicates if the data in the datasource is recent or already old
	Recent *bool `json:"recent" db:"is_recent"`

	// Validity indicates the validity of the data stored in the datasource
	Validity *dbType.Validity `json:"validity" db:"validity"`

	// Duplicates indicates how many duplicates are in the datasource and if
	// this value has been checked
	Duplicates *dbType.CheckedRange `json:"duplicates" db:"duplicates"`

	// Errors indicates how many errors are in the datasource and if the amount
	// of errors has been checked
	Errors *dbType.CheckedRange `json:"errors" db:"errors"`

	// Precision indicates how precise the data contained in the datasource is
	Precision *enums.PrecisionLevel `json:"precision" db:"precision"`

	// Reputation indicates how reputable the data source is
	Reputation *enums.Reputation `json:"reputation" db:"reputation"`

	// Objectivity indicates the measures taken to ensure the objectivity of
	// the datasource
	Objectivity *dbType.DataObjectivity `json:"objectivity" db:"objectivity"`

	// KnownCaptureMethod indicates if a known capture/survey method is used
	// by the data provider/owner to get data for the datasource
	KnownCaptureMethod *bool `json:"knownCaptureMethod" db:"usual_survey_method"`

	// Density describes the density of the data-points in the datasource, and
	// if the density has been checked/evaluated
	Density *dbType.CheckedRange `json:"density" db:"density"`

	// Coverage describes the level of coverage in for the Entities
	Coverage *enums.NoneHighRange `json:"coverage" db:"coverage"`

	// RepresentationalConsistency describes how consistent the data
	// representation is between the datasets and the documentation
	RepresentationalConsistency *enums.NoneHighRange `json:"representationalConsistency" db:"representation_consistency"`

	// LogicalConsistency describes how many contradictions are in the
	// datasource, if the number of contradictions has been checked and if the
	// possible contradictions can be checked with a second datasource
	LogicalConsistency *dbType.LogicalConsistency `json:"logicalConsistency" db:"logical_consistency"`

	// DataDelay contains information about the delay between the data ingress
	// into the datasource and the egress from the datasource
	DataDelay *dbType.DataDelay `json:"dataDelay" db:"data_delay"`

	// DelayInformationTransmission shows how users of the datasource are
	// informed about a delay in the data transmission
	DelayInformationTransmission *enums.DelayInformationTransmission `json:"delayInformationTransmission" db:"delay_information"`

	// PerformanceLimitations shows if there are any performance limitations
	// when accessing the datta source and how much they restrict the usage
	PerformanceLimitations *enums.NoneHighRange `json:"performanceLimitations" db:"performancelimitations"`

	// Availability shows how high the availability of the data source is
	Availability *enums.NoneHighRange `json:"availability" db:"availability"`

	// GDPRCompliant shows if the datasource is GDPR-compliant and data
	// pulled from the datasource may be used without any limitations
	GDPRCompliant *bool `json:"gdprCompliant" db:"gdpr_compliant"`
}
