package game

import "github.com/jmoiron/sqlx"

// List returns all games
func List(db *sqlx.DB) ([]Game, error) {
	list := []Game{}

	const q = `select id, name, developer, releasedate, genre from games`

	if err := db.Select(&list, q); err != nil {
		return nil, err
	}

	return list, nil
}
