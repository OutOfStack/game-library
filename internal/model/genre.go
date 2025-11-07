package model

// Genre represents genre entity
type Genre struct {
	ID     int32  `db:"id"`
	Name   string `db:"name"`
	IGDBID int64  `db:"igdb_id"`
}
