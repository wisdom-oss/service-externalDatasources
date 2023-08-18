package databaseTypes

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type TransformationScript struct {
	ID               uuid.UUID      `json:"id" db:"id"`
	Name             string         `json:"name" db:"name"`
	Description      *string        `json:"description" db:"description"`
	Engine           string         `json:"engine" db:"engine"`
	RequiredPackages pq.StringArray `json:"requiredPackages" db:"required_packages"`
	Contents         string         `json:"contents" db:"contents"`
}
