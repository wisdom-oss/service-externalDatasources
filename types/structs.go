package types

type ExternalAPI struct {
	ID                string            `json:"id" db:"id"`
	Name              string            `json:"name" db:"name"`
	BaseURI           string            `json:"baseUri" db:"baseUri"`
	Description       *string           `json:"description" db:"description"`
	AdditionalHeaders map[string]string `json:"additionalHeaders" db:"additionalHeaders"`
}
