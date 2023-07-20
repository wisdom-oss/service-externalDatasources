package types

import (
	"external-api-service/types/database"
)

type BaseData struct {
	ID          string  `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Description *string `json:"description" db:"description"`
}

type API struct {
	IsSecure          *bool                 `json:"isSecure" db:"is_secure"`
	Host              *string               `json:"host" db:"host"`
	Port              *int                  `json:"port" db:"port"`
	Path              *string               `json:"path" db:"path"`
	AdditionalHeaders *databaseTypes.Tuples `json:"additionalHeaders"`
}

type ExternalDataSource struct {
	BaseData
	Metadata Metadata `json:"metadata"`
	API      API      `json:"api"`
}
