package game

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// ErrNotFound is used when a requested entity with id does not exist
type ErrNotFound struct {
	Entity string
	ID     int64
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("%v with id %v was not found", e.Entity, e.ID)
}

// List returns all games
func List(ctx context.Context, db *sqlx.DB) ([]GetGame, error) {
	list := []Game{}

	const q = `select g.id, g.name, g.developer, g.release_date, g.genre,
	case 
		when CURRENT_DATE >= max(s.begin_date) and CURRENT_DATE <= max(s.end_date) then g.price*((100 - max(sg.discount_percent))/100.0)
		else price
	end as price
	from games g
	left join sales_games sg on sg.game_id = g.id
	inner join sales s on s.id = sg.sale_id
	group by g.id, g.name`

	if err := db.SelectContext(ctx, &list, q); err != nil {
		return nil, err
	}

	getGames := []GetGame{}
	for _, g := range list {
		getGames = append(getGames, *g.mapToGetGame())
	}

	return getGames, nil
}

// Retrieve returns a single game
func Retrieve(ctx context.Context, db *sqlx.DB, id int64) (*GetGame, error) {
	var g Game

	const q = `select g.id, g.name, g.developer, g.release_date, g.genre,
		case
			when CURRENT_DATE >= s.begin_date and CURRENT_DATE <= s.end_date then g.price*((100 - sg.discount_percent)/100.0)
			else price
		end as price
		from games g
		left join sales_games sg on sg.game_id = g.id
		left join sales s on s.id = sg.sale_id
		where g.id = $1
		order by price asc
		limit 1`

	if err := db.GetContext(ctx, &g, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound{"game", id}
		}
		return nil, err
	}

	getGame := g.mapToGetGame()

	return getGame, nil
}

// Create creates a new game
func Create(ctx context.Context, db *sqlx.DB, ng NewGame) (*GetGame, error) {
	const q = `insert into games
	(name, developer, release_date, price, genre)
	values ($1, $2, $3, $4, $5)
	returning id`

	var lastInsertID int64
	err := db.QueryRowContext(ctx, q, ng.Name, ng.Developer, ng.ReleaseDate, ng.Price, pq.StringArray(ng.Genre)).Scan(&lastInsertID)
	if err != nil {
		return nil, fmt.Errorf("inserting game %v: %w", ng, err)
	}

	getGame := ng.mapToGetGame(lastInsertID)

	return getGame, nil
}

// Update modifes information about a game
func Update(ctx context.Context, db *sqlx.DB, id int64, update UpdateGame) (*GetGame, error) {
	g, err := Retrieve(ctx, db, id)
	if err != nil {
		return nil, err
	}

	if update.Name != nil {
		g.Name = *update.Name
	}
	if update.Developer != nil {
		g.Developer = *update.Developer
	}
	if update.ReleaseDate != nil {
		g.ReleaseDate = *update.ReleaseDate
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
	_, err = db.ExecContext(ctx, q, g.Name, g.Developer, g.ReleaseDate, g.Price, pq.StringArray(g.Genre), id)
	if err != nil {
		return nil, errors.Wrap(err, "updating game")
	}
	return g, nil
}

// Delete deletes specified game
func Delete(ctx context.Context, db *sqlx.DB, id int64) error {
	const q = `delete from games where id = $1;`
	res, err := db.ExecContext(ctx, q, id)
	if err != nil {
		return errors.Wrap(err, "deleting game")
	} else {
		count, err := res.RowsAffected()
		if err != nil {
			return errors.Wrap(err, "deleting game. count")
		} else if count == 0 {
			return ErrNotFound{"game", id}
		}
	}
	return nil
}

// AddGameOnSale connects game with a sale
func AddGameOnSale(ctx context.Context, db *sqlx.DB, gameID int64, ngs NewGameSale) (*GetGameSale, error) {
	g, err := Retrieve(ctx, db, gameID)
	if err != nil {
		return nil, err
	}
	s, err := RetrieveSale(ctx, db, ngs.SaleID)
	if err != nil {
		return nil, err
	}
	const q = `insert into sales_games
	(game_id, sale_id, discount_percent)
	values ($1, $2, $3)`

	_, err = db.ExecContext(ctx, q, g.ID, s.ID, ngs.DiscountPercent)
	if err != nil {
		return nil, fmt.Errorf("adding game with id %v on sale  with id %v: %w", gameID, ngs.SaleID, err)
	}

	getGameSale := ngs.mapToGetGameSale(s, gameID)

	return getGameSale, nil
}
