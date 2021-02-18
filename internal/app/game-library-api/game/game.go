package game

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

// List returns all games
func List(db *sqlx.DB) ([]Game, error) {
	list := []Game{}

	const q = `select id, name, developer, releasedate, genre from games`

	if err := db.Select(&list, q); err != nil {
		return nil, err
	}

	return list, nil
}

// Retrieve returns a single game
func Retrieve(db *sqlx.DB, id uint64) (*Game, error) {
	var g Game

	const q = `select id, name, developer, releasedate, genre 
		from games
		where id = $1`

	if err := db.Get(&g, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &g, nil
}
