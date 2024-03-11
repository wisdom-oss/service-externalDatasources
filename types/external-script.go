package types

import "github.com/jackc/pgx/v5/pgtype"

type ExternalTransformationScript struct {
	ID               pgtype.UUID `json:"id" db:"id"`
	Name             pgtype.Text `json:"name" db:"name"`
	Description      pgtype.Text `json:"description" db:"description"`
	RequiredPackages []string    `json:"requiredPackages" db:"requiredPackages"`
	ScriptContents   pgtype.Text `json:"contents" db:"contents"`
}
