package game

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/OutOfStack/game-library/pkg/types"
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

// ListInfo returns all games with extended properties
func ListInfo(ctx context.Context, db *sqlx.DB) ([]GameInfo, error) {
	list := []GameInfo{}

	// here we unite two groups of games - games which are currently on sale and the rest ones
	// if there are currently more than one sale (maybe should be restricted) and a game is on both sales we choose max discount
	// after uniting we count rating for all games
	const q = `
	select all_g.*, coalesce(avg(r.rating), 0) as rating from (
		with games_on_sale as (
			select g.id, g.name, g.developer, g.publisher, g.release_date, g.genre, g.price, max(sg.discount_percent) as discount
			from sales s
			inner join sales_games sg on sg.sale_id = s.id
			inner join games g on g.id = sg.game_id
			where CURRENT_DATE >= s.begin_date and CURRENT_DATE <= s.end_date
			group by g.id
		)
		select gs.id, gs.name, gs.developer, gs.publisher, gs.release_date, gs.genre, gs.price, gs.price * ((100 - gs.discount) / 100.0) as current_price 
		from games_on_sale gs
		union
		select g.id, g.name, g.developer, g.publisher, g.release_date, g.genre, g.price, g.price
		from games g
		where g.id not in (select id from games_on_sale)
	) all_g
	left join ratings r on r.game_id = all_g.id
	group by all_g.id, all_g.name, all_g.developer, all_g.publisher, all_g.release_date, all_g.genre, all_g.price, all_g.current_price
	order by all_g.id`

	if err := db.SelectContext(ctx, &list, q); err != nil {
		return nil, err
	}

	return list, nil
}

// RetrieveInfo returns a single game with extended properties
func RetrieveInfo(ctx context.Context, db *sqlx.DB, id int64) (*GameInfo, error) {
	var g GameInfo

	const q = `
	select g.id, g.name, g.developer, g.publisher, g.release_date, g.genre, g.price, 
		g.price * (100 - coalesce(
			(select discount_percent
			from sales_games sg 
			inner join sales s on s.id = sg.sale_id 
			where sg.game_id = g.id and CURRENT_DATE >= s.begin_date and CURRENT_DATE <= s.end_date
			order by discount_percent desc
			limit 1),
			0)) / 100.0 as current_price,
			coalesce(avg(r.rating), 0) as rating
	from games g
	left join ratings r on r.game_id = g.id
	where g.id = $1
	group by g.id`

	if err := db.GetContext(ctx, &g, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound{"game", id}
		}
		return nil, err
	}

	return &g, nil
}

// Retrieve returns a single game
func Retrieve(ctx context.Context, db *sqlx.DB, id int64) (*Game, error) {
	var g Game

	const q = `select g.id, g.name, g.developer, g.publisher, g.release_date, g.genre, g.price
	from games g
	where g.id = $1`

	if err := db.GetContext(ctx, &g, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound{"game", id}
		}
		return nil, err
	}

	return &g, nil
}

// Create creates a new game
func Create(ctx context.Context, db *sqlx.DB, cg CreateGame) (int64, error) {
	const q = `insert into games
	(name, developer, publisher, release_date, price, genre)
	values ($1, $2, $3, $4, $5, $6)
	returning id`

	var lastInsertID int64
	err := db.QueryRowContext(ctx, q, cg.Name, cg.Developer, cg.Publisher, cg.ReleaseDate, cg.Price, pq.StringArray(cg.Genre)).Scan(&lastInsertID)
	if err != nil {
		return 0, fmt.Errorf("inserting game %v: %w", cg, err)
	}

	return lastInsertID, nil
}

// Update modifes information about a game
func Update(ctx context.Context, db *sqlx.DB, id int64, update UpdateGame) (*Game, error) {
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
	if update.Publisher != nil {
		g.Publisher = *update.Publisher
	}
	if update.ReleaseDate != nil {
		releaseDate, _ := types.ParseDate(*update.ReleaseDate)
		g.ReleaseDate = releaseDate
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
	 publisher = $3,
	 release_date = $4,
	 price = $5,
	 genre = $6
	 where id = $7`
	_, err = db.ExecContext(ctx, q, g.Name, g.Developer, g.Publisher, g.ReleaseDate.String(), g.Price, pq.StringArray(g.Genre), id)
	if err != nil {
		return nil, errors.Wrap(err, "updating game")
	}
	return g, nil
}

// Delete deletes specified game
func Delete(ctx context.Context, db *sqlx.DB, id int64) error {
	const q = `delete from games where id = $1`
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
func AddGameOnSale(ctx context.Context, db *sqlx.DB, gameID int64, cgs CreateGameSale) (*GameSale, error) {
	g, err := Retrieve(ctx, db, gameID)
	if err != nil {
		return nil, err
	}
	s, err := RetrieveSale(ctx, db, cgs.SaleID)
	if err != nil {
		return nil, err
	}
	const q = `insert into sales_games
	(game_id, sale_id, discount_percent)
	values ($1, $2, $3)
	on conflict (game_id, sale_id) do update set discount_percent = $3`

	_, err = db.ExecContext(ctx, q, g.ID, s.ID, cgs.DiscountPercent)
	if err != nil {
		return nil, fmt.Errorf("adding game with id %v on sale  with id %v: %w", gameID, cgs.SaleID, err)
	}

	gameSale := cgs.MapToGameSale(s, gameID)

	return gameSale, nil
}

// ListGameSales returns all sales for specified game
func ListGameSales(ctx context.Context, db *sqlx.DB, gameID int64) ([]GameSale, error) {
	_, err := Retrieve(ctx, db, gameID)
	if err != nil {
		return nil, err
	}

	gameSales := []GameSale{}

	const q = `select sg.game_id, sg.sale_id, s.name as sale, s.begin_date, s.end_date, sg.discount_percent
	from sales_games sg
	left join sales s on s.id = sg.sale_id
	left join games g on g.id = sg.game_id
	where sg.game_id = $1`
	if err := db.SelectContext(ctx, &gameSales, q, gameID); err != nil {
		return nil, errors.Wrapf(err, "selecting sales of game with id %q", gameID)
	}

	return gameSales, nil
}
