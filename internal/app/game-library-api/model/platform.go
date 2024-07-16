package model

// Platform represents platform entity
type Platform struct {
	ID           int32  `db:"id"`
	Name         string `db:"name"`
	Abbreviation string `db:"abbreviation"`
	IGDBID       int64  `db:"igdb_id"`
}
