CREATE TYPE external_data_sources.dataReference AS
(
    -- the topic of the data
    topic              text,
    -- the local reference area
    local_reference    text,
    -- the temporal reference
    temporal_reference daterange
);

-- name: create-type-dataOrigin
CREATE TYPE external_data_sources.dataOrigin AS
(
    -- the name of the data provider
    provider text,
    -- the name of the data creator
    creator  text,
    -- the name of the data owner
    owner    text
);

-- name: create-type-tuple
CREATE TYPE external_data_sources.tuple AS
(
    -- a non-unique "key"
    x text,
    -- the value for the "key"
    y text
);

-- name: create-type-documentation
CREATE TYPE external_data_sources.documentation AS
(
    -- the type of the documentation (a book, a pdf, a website, etc.)
    type      text,
    -- the location of the documentation (a shelf number, url, etc.)
    location  text,
    -- the verbosity of this documentation entry
    verbosity noneHighRange
);

-- name: create-enum-costModel
CREATE TYPE external_data_sources.billingModel AS ENUM ('openSource', 'free', 'singlePurchase', 'byTime', 'byAccess', 'byData');

-- name: create-type-accessCosts
CREATE TYPE external_data_sources.billingInformation AS
(
    -- the model used for calculating the costs
    model        billingModel,
    -- the price per access to the
    pricePerUnit numeric
);

-- name: create-type-provisionInformation
CREATE TYPE external_data_sources.provisionInformation AS
(
    -- the type of the data source
    type   text,
    -- the format of the data source, should be a MIME-type at best
    format text
);

-- name: create-enum-validity
CREATE TYPE external_data_sources.validity AS ENUM ('fully', 'partially', 'none');

-- name: create-enum-none-high-range
CREATE TYPE external_data_sources.noneHighRange AS ENUM ('none', 'low', 'medium', 'high');

-- name: create-enum-precision
CREATE TYPE external_data_sources.precisionLevel AS ENUM ('fine', 'usual', 'unusual', 'imprecise');

-- name: create-enum-reputation
CREATE TYPE external_data_sources.reputation AS ENUM ('independent_and_external', 'independent_or_external', 'suspected_high', 'suspected_low');

-- name: create-type-checked-range
CREATE TYPE external_data_sources.checkedRange AS
(
    checked bool,
    range   noneHighRange
);

CREATE TYPE external_data_sources.objectivity AS
(
    conflict_of_interest bool,
    raw_data             bool,
    automatic_capture    bool
);

CREATE TYPE external_data_sources.logicalConsistency AS
(
    checked                   bool,
    contradictions_examinable bool,
    range                     noneHighRange
);

CREATE TYPE external_data_sources.delay AS
(
    source    noneHighRange,
    recording noneHighRange
);

CREATE TYPE external_data_sources.header AS
(
    key    text,
    value  text,
    secure bool
);

CREATE TYPE external_data_sources.delayInformation AS ENUM ('direct', 'automatic', 'manual', 'none');

-- name: create-table-sources
-- This table contains the following basic data of external sources:
--   * ID
--   * Name
--   * Description
--   * hasExternalAPI
--   * hasMetadata
CREATE TABLE external_data_sources.sources
(
    -- the unique identifier of the data source
    id          uuid primary key default gen_random_uuid(),
    -- the name of the data source
    name        text not null,
    -- the optional description of the data source
    description text             default null
);

-- name: create-table-metadata
CREATE TABLE external_data_sources.metadata
(
    -- the id for the metadata which needs to be the id of the external source
    id                         uuid primary key references sources (id),
    -- additional reference data describing the external source
    reference                  dataReference        default null,
    -- a description of the data origin
    origin                     dataOrigin           default null,
    -- a list of distinctive features as tuples
    distinctive_features       tuple[]              default null,
    -- rights for using the external data
    usage_rights               text                 default null,
    -- duties for using the external data
    usage_duties               text                 default null,
    -- real entities that are represented in the data source
    real_entities              text[]               default null,
    -- contact data for a local expert
    local_expert               tuple[]              default null,
    -- a array of documentation entries for this source
    external_documentation     documentation[]      default null,
    -- a interval at which the data source is updated
    update_rate                interval             default null,
    -- a array of ISO 639-1 language codes
    languages                  text[]               default null,
    -- a object containing the billing information
    billing                    billingInformation   default null,
    -- a object containing information about how the data is provisioned
    provisioning               provisionInformation default null,
    -- a id referencing a service from which the data is derived from
    derived_from               uuid                 default null,
    -- the recency of the data
    is_recent                  bool                 default false,
    -- the validity of the data
    validity                   validity             default null,
    -- indicator if the about duplicates in the data source
    duplicates                 checkedRange         default null,
    -- indicator if the about the errors in the data source
    errors                     checkedRange         default null,
    -- indicator for how precise the dataset is
    precision                  precisionLevel       default null,
    -- indicator for the reputation
    reputation                 reputation           default null,
    -- indicator for the objectivity
    objectivity                objectivity          default null,
    -- indicator if a usual survey method was used
    usual_survey_method        bool                 default false,
    -- indicator about how dense the data is
    density                    checkedRange         default null,
    -- indicator about how good the data matches the real entities
    coverage                   noneHighRange        default null,
    --
    representation_consistency noneHighRange        default null,
    --
    logical_consistency        logicalConsistency   default null,
    --
    delay                      delay                default null,
    --
    delay_information          delayInformation     default null,
    --
    performanceLimitations     noneHighRange        default null,
    --
    availability               noneHighRange        default null,
    --
    gdpr_compliant             bool                 default false
);

CREATE TABLE external_data_sources.apis
(
    -- the id for the metadata which needs to be the id of the external source
    id                 uuid primary key references sources (id),
    -- is https access
    is_secure          bool          default true,
    -- the host on which the api resides
    host               text not null,
    -- the port of the  api
    port               int  not null default 443,
    -- the path of the main api endpoint
    path               text not null,
    -- additional headers that need to be set
    additional_headers header[]      default null
);

CREATE VIEW external_data_sources.info AS
(
SELECT
    s.*,
    reference,
    origin,
    distinctive_features,
    usage_rights,
    usage_duties,
    real_entities,
    local_expert,
    external_documentation,
    update_rate,
    languages,
    billing,
    provisioning,
    derived_from,
    is_recent,
    validity,
    duplicates,
    errors,
    precision,
    reputation,
    objectivity,
    usual_survey_method,
    density,
    coverage,
    representation_consistency,
    logical_consistency,
    delay,
    delay_information,
    performancelimitations,
    availability,
    gdpr_compliant,
    is_secure,
    host,
    port,
    path,
    additional_headers
FROM
    external_data_sources.sources s
        LEFT JOIN external_data_sources.metadata m on s.id = m.id
        LEFT JOIN external_data_sources.apis a on s.id = a.id
    );

--name: get-all-information
SELECT *
FROM
    external_data_sources.info;

--name: add-base-data
INSERT INTO
    external_data_sources.sources (name, description)
VALUES
    ($1, $2)
RETURNING id;

-- name: add-metadata
INSERT INTO
    external_data_sources.metadata
(
    id,
    reference,
    origin,
    distinctive_features,
    usage_rights,
    usage_duties,
    real_entities,
    local_expert,
    external_documentation,
    update_rate,
    languages,
    billing,
    provisioning,
    derived_from,
    is_recent,
    validity,
    duplicates,
    errors,
    precision,
    reputation,
    objectivity,
    usual_survey_method,
    density,
    coverage,
    representation_consistency,
    logical_consistency,
    delay,
    delay_information,
    performancelimitations,
    availability,
    gdpr_compliant)
VALUES
    (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10,
        $11,
        $12,
        $13,
        $14,
        $15,
        $16,
        $17,
        $18,
        $19,
        $20,
        $21,
        $22,
        $23,
        $24,
        $25,
        $26,
        $27,
        $28,
        $29,
        $30,
        $31);

-- name: add-api
INSERT INTO
    external_data_sources.apis (id, is_secure, host, port, path, additional_headers)
VALUES
($1, $2, $3, $4, $5, $6);

-- name: get-single-source
SELECT * FROM external_data_sources.info
WHERE id = $1::uuid OR name = $1::text;