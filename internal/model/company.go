package model

import "database/sql"

// company types
const (
	CompanyTypeDeveloper = "dev"
	CompanyTypePublisher = "pub"
)

// Company represents company entity
type Company struct {
	ID     int32         `db:"id"`
	Name   string        `db:"name"`
	IGDBID sql.NullInt64 `db:"igdb_id"`
}
