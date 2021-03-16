package game

import (
	"context"
	"database/sql"
	"fmt"

	"cloud.google.com/go/civil"
	"github.com/OutOfStack/game-library/pkg/types"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var (
	// ErrNotFound is used when a requested entity with id does not exist
	ErrNotFound = errors.New("game not found")
)

// List returns all games
func List(ctx context.Context, db *sqlx.DB) ([]Game, error) {
	list := []Game{}

	const q = `select g.id, g.name, g.developer, g.release_date, g.genre,
	case 
		when CURRENT_DATE >= max(s.begin_date) and CURRENT_DATE <= max(s.end_date) then g.price*((100 - max(s.discount_percent))/100.0)
		else price
	end as price
	from games g
	left join sales s on s.game_id = g.id
	group by g.id, g.name`

	if err := db.SelectContext(ctx, &list, q); err != nil {
		return nil, err
	}

	return list, nil
}

// Retrieve returns a single game
func Retrieve(ctx context.Context, db *sqlx.DB, id int64) (*Game, error) {
	var g Game

	const q = `select g.id, g.name, g.developer, g.release_date, g.genre,
		case 
			when CURRENT_DATE >= s.begin_date and CURRENT_DATE <= s.end_date then g.price*((100 - s.discount_percent)/100.0)
			else price
		end as price
		from games g
		left join sales s on s.game_id = g.id
		where g.id = $1
		order by price asc
		limit 1`

	if err := db.GetContext(ctx, &g, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &g, nil
}

// Create creates a new game
func Create(ctx context.Context, db *sqlx.DB, ng NewGame) (*Game, error) {
	date, err := civil.ParseDate(ng.ReleaseDate)
	if err != nil {
		return nil, fmt.Errorf("parsing releaseDate: %w", err)
	}
	const q = `insert into games
	(name, developer, release_date, price, genre)
	values ($1, $2, $3, $4, $5)
	returning id`

	var lastInsertID int64
	err = db.QueryRowContext(ctx, q, ng.Name, ng.Developer, ng.ReleaseDate, ng.Price, ng.Genre).Scan(&lastInsertID)
	if err != nil {
		return nil, fmt.Errorf("inserting game %v: %w", ng, err)
	}

	g := Game{
		ID:          lastInsertID,
		Name:        ng.Name,
		Developer:   ng.Developer,
		ReleaseDate: types.Date(date),
		Genre:       ng.Genre,
	}

	return &g, nil
}

// Update modifes information about a game
func Update(ctx context.Context, db *sqlx.DB, id int64, update UpdateGame) error {
	g, err := Retrieve(ctx, db, id)
	if err != nil {
		return err
	}

	if update.Name != nil {
		g.Name = *update.Name
	}
	if update.Developer != nil {
		g.Developer = *update.Developer
	}
	if update.ReleaseDate != nil {
		date, err := civil.ParseDate(*update.ReleaseDate)
		if err != nil {
			return errors.Wrap(err, "parsing releaseDate")
		}
		g.ReleaseDate = types.Date(date)
	}
	if update.Price != nil {
		g.Price = *update.Price
	}
	if update.Genre != nil {
		g.Genre = *update.Genre
	}
	const q = `update games set 
	 name = $1,
	 developer = $2,
	 release_date = $3,
	 price = $4,
	 genre = $5
	 where id = $6;`
	_, err = db.ExecContext(ctx, q, g.Name, g.Developer, g.ReleaseDate.String(), g.Price, g.Genre, id)
	if err != nil {
		return errors.Wrap(err, "updating product")
	}
	return nil
}
