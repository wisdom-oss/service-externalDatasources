/* This file contains all prepared sql queries related to the creation of
   the composite types used for this microservice on the database side
 */

-- name: create-type-tuple
-- this type allows storing data as tuples
CREATE TYPE Tuple AS
(
    -- a non-unique "key"
    x text,
    -- the value for the "key"
    y text
);

-- name: create-type-referenceData
-- this type organizes the reference data of a data source
CREATE TYPE ReferenceData AS
(
    -- topic contains the topic of the data source
    topic text,
    -- geo_reference contains a reference to a geographic area
    geo_reference text,
    -- temporal_reference contains a daterange containing the temporal range
    -- of the available data
    temporal_reference daterange
);

-- name: create-type-originData
-- this type organizes the data about the provider, creator and owner
-- of the data which are returned by the external data source
CREATE TYPE OriginData AS
(
  -- provider contains the data of the provider for the data
  provider Tuple[],
  -- creator contains the data of the creator of the data source
  creator Tuple[],
  -- owner contains the information about the owner of the data
  -- contained in the data source
  owner Tuple[]
);

-- name: create-type-documentation
CREATE TYPE Documentation AS
(
    -- type contains a string describing the type of documentation
    type text,
    -- location contains information about the location of the documentation
    location Tuple[],
    -- verbosity contains a level of verbosity of the documentation
    verbosity noneHighRange
);

-- name: create-type-dataProvision
-- this type contains information about how the data will be provided
CREATE TYPE dataProvision AS
(
    -- type contains the type of the datasource (e.g. REST-API, Database)
    type text,
    -- format contains a MIME type describing the format in which the
    -- data is returned. in case of databases the format may contain
    -- the database type (e.g. MariaDB, PostgreSQL, etc.)
    format text
);

-- name: create-type-checkedRange
-- this type allows to specify a range and using an indicator for a check for
-- the range value
CREATE TYPE checkedRange AS
(
    checked bool,
    range nonehighrange
);

-- name: create-type-dataObjectivity
CREATE TYPE dataObjectivity AS
(
    conflict_of_interest bool,
    raw_data bool,
    automatic_capture bool
);

-- name: create-type-logicalConsistency
-- this type contains information about the consistency of a data source
CREATE TYPE logicalConsistency AS
(
    checked bool,
    contradictions_examinable bool,
    rang nonehighrange
);

-- name: create-type-dataDelay
-- this type contains information about the delays for data ingress and egress
-- from the external source
CREATE TYPE dataDelay AS
(
    ingress nonehighrange,
    egress nonehighrange
);

-- name: create-type-header
-- this type represents a http header that could be attached to the source
CREATE TYPE header AS
(
    key text,
    value text
);

-- name: create-type-authorization-information
-- this type represents the needed authorization information to access a data
-- source
CREATE TYPE authorizationInformation AS
(
    location authenticationLocation,
    key text,
    value text
)
