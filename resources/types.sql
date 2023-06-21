CREATE TYPE dataReference AS (
    -- the topic of the data
    topic text,
    -- the local reference area
    local_reference text,
    -- the temporal reference
    temporal_reference daterange
);

-- name: create-type-dataOrigin
CREATE TYPE dataOrigin AS (
    -- the name of the data provider
    provider text,
    -- the name of the data creator
    creator text,
    -- the name of the data owner
    owner text
);

-- name: create-type-tuple
CREATE TYPE tuple AS (
    -- a non-unique "key"
    x text,
    -- the value for the "key"
    y text
);

-- name: create-type-documentation
CREATE TYPE documentation AS (
    -- the type of the documentation (a book, a pdf, a website, etc.)
    type text,
    -- the location of the documentation (a shelf number, url, etc.)
    location text
);

-- name: create-enum-costModel
CREATE TYPE billingModel AS ENUM ('openSource', 'free', 'singlePurchase', 'byTime', 'byAccess', 'byData');

-- name: create-type-accessCosts
CREATE TYPE billingInformation AS (
    -- the model used for calculating the costs
    model billingModel,
    -- the price per access to the
    pricePerUnit numeric
);