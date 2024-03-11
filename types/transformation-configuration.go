package types

import "github.com/jackc/pgx/v5/pgtype"

type TransformationConfiguration struct {
	ID         int                    `json:"id" db:"id"`
	Datasource pgtype.UUID            `json:"dataSource" db:"datasource"`
	Action     pgtype.Text            `json:"action" db:"action"`
	Options    map[string]interface{} `json:"options" db:"data"`
	Priority   int                    `json:"priority" db:"priority"`
}
