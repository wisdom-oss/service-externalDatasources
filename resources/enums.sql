/* This file contains all prepared sql queries related to the creation of
   the enumerations used for this microservice on the database side
 */

-- name: create-enum-noneHighRange
-- this enumeration allows the usage of a range from none to high
-- for some of the composite types or columns
CREATE TYPE noneHighRange AS ENUM ('none', 'low', 'medium', 'high');

-- name: create-enum-validity
-- this enumeration contains the three levels of validity used in the service
CREATE TYPE validityLevel AS ENUM ('none', 'partially', 'fully');

-- name: create-enum-pricingModel
-- this enumeration contains the five possible types of pricing used
-- for the external data source
CREATE TYPE pricingModel AS ENUM ('openSource', 'free', 'singlePurchase', 'perAccess', 'perTimeUnit', 'perDataAmount');

-- name: create-precisionLevel
-- the enumeration contains the three precision levels available for a
-- external data source
CREATE TYPE precisionLevel AS ENUM ('imprecise', 'unusual', 'unusual', 'fine');