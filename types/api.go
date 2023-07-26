package types

import databaseTypes "external-api-service/types/database"

type API struct {
	IsSecure          *bool                 `json:"isSecure" db:"is_secure"`
	Host              *string               `json:"host" db:"host"`
	Port              *uint16               `json:"port" db:"port"`
	Path              *string               `json:"path" db:"path"`
	AdditionalHeaders *databaseTypes.Tuples `json:"additionalHeaders"`
}

// IsValid checks if at least the host was set in the configuration
func (a API) IsValid() bool {
	return a.Host != nil
}
