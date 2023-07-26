package types

type BaseData struct {
	ID          string  `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Description *string `json:"description" db:"description"`
}

type ExternalDataSource struct {
	BaseData
	Metadata Metadata `json:"metadata"`
	API      API      `json:"api"`
}
