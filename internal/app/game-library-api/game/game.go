package game

import (
	"database/sql"
	"errors"
	"fmt"

	"cloud.google.com/go/civil"
	"github.com/OutOfStack/game-library/pkg/types"
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

// Create creates a new game
func Create(db *sqlx.DB, pm PostModel) (*Game, error) {
	date, err := civil.ParseDate(pm.ReleaseDate)
	if err != nil {
		return nil, fmt.Errorf("parsing releaseDate: %w", err)
	}
	const q = `insert into games
	(name, developer, releasedate, genre)
	values ($1, $2, $3, $4)
	returning id`

	var lastInsertID uint64
	err = db.QueryRow(q, pm.Name, pm.Developer, pm.ReleaseDate, pm.Genre).Scan(&lastInsertID)
	if err != nil {
		return nil, fmt.Errorf("inserting game %v: %w", pm, err)
	}

	g := Game{
		ID:          lastInsertID,
		Name:        pm.Name,
		Developer:   pm.Developer,
		ReleaseDate: types.Date(date),
		Genre:       pm.Genre,
	}

	return &g, nil
}
