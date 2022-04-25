package game

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/OutOfStack/game-library/pkg/types"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/api/trace"
)

// ErrNotFound is used when a requested entity with id does not exist
type ErrNotFound struct {
	Entity string
	ID     int64
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("%v with id %v was not found", e.Entity, e.ID)
}

// GetInfos returns list of games with extended properties. Limited by pageSize and starting Id
func GetInfos(ctx context.Context, db *sqlx.DB, pageSize int, lastId int64) ([]GameInfo, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "sql.game.getinfos")
	defer span.End()

	list := []GameInfo{}

	// first, data set is limited by page size and last id,
	// then we unite two groups of games - games which are currently on sale and the rest ones,
	// if there are currently more than one sale (maybe should be restricted) and a game is on both sales we choose max discount
	// after uniting we count rating for selected games
	const q = `
	select all_g.*, coalesce(avg(r.rating), 0) as rating from (
		with page as (select id, name, developer, publisher, release_date, genre, price, logo_url
				from games
				where id > $1
				order by id
				fetch first $2 rows only),
		on_sale as (
			select g.id, g.name, g.developer, g.publisher, g.release_date, g.genre, g.price, max(sg.discount_percent) as discount, g.logo_url
			from page g
			inner join sales_games sg on sg.game_id = g.id
			inner join sales s on s.id = sg.sale_id
			where CURRENT_DATE >= s.begin_date and CURRENT_DATE <= s.end_date
			group by g.id, g.name, g.developer, g.publisher, g.release_date, g.genre, g.price, g.logo_url
		)
		select os.id, os.name, os.developer, os.publisher, os.release_date, os.genre, os.price, os.price * ((100 - os.discount) / 100.0) as current_price, os.logo_url
		from on_sale os
		union all
		select g.id, g.name, g.developer, g.publisher, g.release_date, g.genre, g.price, g.price as current_price, g.logo_url
		from page g
		where g.id not in (select id from on_sale)
	) all_g
	left join ratings r on r.game_id = all_g.id
	group by all_g.id, all_g.name, all_g.developer, all_g.publisher, all_g.release_date, all_g.genre, all_g.price, all_g.current_price, all_g.logo_url
	order by all_g.id`

	if err := db.SelectContext(ctx, &list, q, lastId, pageSize); err != nil {
		return nil, err
	}

	return list, nil
}

// RetrieveInfo returns a single game with extended properties
// If such entity does not exist returns error ErrNotFound{}
func RetrieveInfo(ctx context.Context, db *sqlx.DB, id int64) (*GameInfo, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "sql.game.retrieveinfo")
	defer span.End()

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
		g.logo_url,
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

// SearchInfos returns list of games with extended properties limited by search query
func SearchInfos(ctx context.Context, db *sqlx.DB, search string) ([]GameInfo, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "sql.game.searchinfos")
	defer span.End()

	list := []GameInfo{}

	// first, data set is limited by search query,
	// then we unite two groups of games - games which are currently on sale and the rest ones,
	// if there are currently more than one sale (maybe should be restricted) and a game is on both sales we choose max discount
	// after uniting we count rating for selected games
	const q = `
	select all_g.*, coalesce(avg(r.rating), 0) as rating from (
		with filtered as (select id, name, developer, publisher, release_date, genre, price, logo_url
				from games
				where lower(name) like $1),
		on_sale as (
			select f.id, f.name, f.developer, f.publisher, f.release_date, f.genre, f.price, max(sg.discount_percent) as discount, f.logo_url
			from filtered f
			inner join sales_games sg on sg.game_id = f.id
			inner join sales s on s.id = sg.sale_id
			where CURRENT_DATE >= s.begin_date and CURRENT_DATE <= s.end_date
			group by f.id, f.name, f.developer, f.publisher, f.release_date, f.genre, f.price, f.logo_url
		)
		select os.id, os.name, os.developer, os.publisher, os.release_date, os.genre, os.price, os.price * ((100 - os.discount) / 100.0) as current_price, os.logo_url
		from on_sale os
		union all
		select f.id, f.name, f.developer, f.publisher, f.release_date, f.genre, f.price, f.price as current_price, f.logo_url
		from filtered f
		where f.id not in (select id from on_sale)
	) all_g
	left join ratings r on r.game_id = all_g.id
	group by all_g.id, all_g.name, all_g.developer, all_g.publisher, all_g.release_date, all_g.genre, all_g.price, all_g.current_price, all_g.logo_url`

	if err := db.SelectContext(ctx, &list, q, strings.ToLower(search)+"%"); err != nil {
		return nil, err
	}

	return list, nil
}

// Retrieve returns a single game
// If such entity does not exist returns error ErrNotFound{}
func Retrieve(ctx context.Context, db *sqlx.DB, id int64) (*Game, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "sql.game.retrieve")
	defer span.End()

	var g Game

	const q = `select g.id, g.name, g.developer, g.publisher, g.release_date, g.genre, g.price, g.logo_url
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
func Create(ctx context.Context, db *sqlx.DB, cg CreateGameReq) (int64, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "sql.game.create")
	defer span.End()

	const q = `insert into games
	(name, developer, publisher, release_date, price, genre, logo_url)
	values ($1, $2, $3, $4, $5, $6, $7)
	returning id`

	var lastInsertID int64
	err := db.QueryRowContext(ctx, q, cg.Name, cg.Developer, cg.Publisher, cg.ReleaseDate, cg.Price, pq.StringArray(cg.Genre), cg.LogoUrl).Scan(&lastInsertID)
	if err != nil {
		return 0, fmt.Errorf("inserting game %v: %w", cg, err)
	}

	return lastInsertID, nil
}

// Update modifes information about a game
// If such entity does not exist returns error ErrNotFound{}
func Update(ctx context.Context, db *sqlx.DB, id int64, update UpdateGameReq) (*Game, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "sql.game.update")
	defer span.End()

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
	if update.LogoUrl != nil {
		g.LogoUrl = sql.NullString{
			String: *update.LogoUrl,
			Valid:  true,
		}
	}
	const q = `update games set 
	 name = $1,
	 developer = $2,
	 publisher = $3,
	 release_date = $4,
	 price = $5,
	 genre = $6,
	 logo_url = $7
	 where id = $8`
	_, err = db.ExecContext(ctx, q, g.Name, g.Developer, g.Publisher, g.ReleaseDate.String(), g.Price, pq.StringArray(g.Genre), g.LogoUrl, id)
	if err != nil {
		return nil, errors.Wrap(err, "updating game")
	}
	return g, nil
}

// Delete deletes specified game
// If such entity does not exist returns error ErrNotFound{}
func Delete(ctx context.Context, db *sqlx.DB, id int64) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "sql.game.delete")
	defer span.End()

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
// If such game or sale does not exist returns error ErrNotFound{}
func AddGameOnSale(ctx context.Context, db *sqlx.DB, gameID int64, cgs CreateGameSaleReq) (*GameSale, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "sql.game.addgameonsale")
	defer span.End()

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
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "sql.game.listgamesales")
	defer span.End()

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
