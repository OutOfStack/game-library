package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/OutOfStack/game-library/pkg/types"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
)

// Storage provides required dependencies for repository
type Storage struct {
	DB *sqlx.DB
}

// ErrNotFound is used when a requested entity with id does not exist
type ErrNotFound struct {
	Entity string
	ID     int64
}

var tracer = otel.Tracer("")

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("%v with id %v was not found", e.Entity, e.ID)
}

// GetInfos returns list of games with extended properties. Limited by pageSize and starting Id
func (s *Storage) GetInfos(ctx context.Context, pageSize int, lastID int64) ([]GameExt, error) {
	ctx, span := tracer.Start(ctx, "db.game.getinfos")
	defer span.End()

	list := []GameExt{}

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

	if err := s.DB.SelectContext(ctx, &list, q, lastID, pageSize); err != nil {
		return nil, err
	}

	return list, nil
}

// RetrieveInfo returns a single game with extended properties
// If such entity does not exist returns error ErrNotFound{}
func (s *Storage) RetrieveInfo(ctx context.Context, id int64) (*GameExt, error) {
	ctx, span := tracer.Start(ctx, "db.game.retrieveinfo")
	defer span.End()

	var g GameExt

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

	if err := s.DB.GetContext(ctx, &g, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound{"game", id}
		}
		return nil, err
	}

	return &g, nil
}

// SearchInfos returns list of games with extended properties limited by search query
func (s *Storage) SearchInfos(ctx context.Context, search string) ([]GameExt, error) {
	ctx, span := tracer.Start(ctx, "db.game.searchinfos")
	defer span.End()

	list := []GameExt{}

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

	if err := s.DB.SelectContext(ctx, &list, q, strings.ToLower(search)+"%"); err != nil {
		return nil, err
	}

	return list, nil
}

// Retrieve returns a single game
// If such entity does not exist returns error ErrNotFound{}
func (s *Storage) Retrieve(ctx context.Context, id int64) (*Game, error) {
	ctx, span := tracer.Start(ctx, "db.game.retrieve")
	defer span.End()

	var g Game

	const q = `select g.id, g.name, g.developer, g.publisher, g.release_date, g.genre, g.price, g.logo_url
	from games g
	where g.id = $1`

	if err := s.DB.GetContext(ctx, &g, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound{"game", id}
		}
		return nil, err
	}

	return &g, nil
}

// Create creates a new game
func (s *Storage) Create(ctx context.Context, cg CreateGame) (int64, error) {
	ctx, span := tracer.Start(ctx, "db.game.create")
	defer span.End()

	const q = `insert into games
	(name, developer, publisher, release_date, price, genre, logo_url)
	values ($1, $2, $3, $4, $5, $6, $7)
	returning id`

	var lastInsertID int64
	err := s.DB.QueryRowContext(ctx, q, cg.Name, cg.Developer, cg.Publisher, cg.ReleaseDate, cg.Price, pq.StringArray(cg.Genre), cg.LogoURL).Scan(&lastInsertID)
	if err != nil {
		return 0, fmt.Errorf("inserting game %v: %w", cg, err)
	}

	return lastInsertID, nil
}

// Update modifes information about a game
// If such entity does not exist returns error ErrNotFound{}
func (s *Storage) Update(ctx context.Context, ug UpdateGame) error {
	ctx, span := tracer.Start(ctx, "db.game.update")
	defer span.End()

	const q = `update games set 
	 name = $1,
	 developer = $2,
	 publisher = $3,
	 release_date = $4,
	 price = $5,
	 genre = $6,
	 logo_url = $7
	 where id = $8`

	releaseDate, err := types.ParseDate(ug.ReleaseDate)
	if err != nil {
		return errors.Wrapf(err, "invalid date: %s", releaseDate.String())
	}
	res, err := s.DB.ExecContext(ctx, q, ug.Name, ug.Developer, ug.Publisher, releaseDate.String(), ug.Price, pq.StringArray(ug.Genre), ug.LogoURL, ug.ID)
	if err != nil {
		return errors.Wrap(err, "updating game")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrapf(err, "getting rows affected on game %d update", ug.ID)
	}
	if rowsAffected == 0 {
		return ErrNotFound{"game", ug.ID}
	}

	return nil
}

// Delete deletes specified game
// If such entity does not exist returns error ErrNotFound{}
func (s *Storage) Delete(ctx context.Context, id int64) error {
	ctx, span := tracer.Start(ctx, "db.game.delete")
	defer span.End()

	const q = `delete from games where id = $1`
	res, err := s.DB.ExecContext(ctx, q, id)
	if err != nil {
		return errors.Wrap(err, "deleting game")
	}
	count, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "deleting game. count")
	}
	if count == 0 {
		return ErrNotFound{"game", id}
	}
	return nil
}

// AddGameOnSale connects game with a sale
// If such game or sale does not exist returns error ErrNotFound{}
func (s *Storage) AddGameOnSale(ctx context.Context, cgs CreateGameSale) error {
	ctx, span := tracer.Start(ctx, "db.game.addgameonsale")
	defer span.End()

	const q = `insert into sales_games
	(game_id, sale_id, discount_percent)
	values ($1, $2, $3)
	on conflict (game_id, sale_id) do update set discount_percent = $3`

	_, err := s.DB.ExecContext(ctx, q, cgs.GameID, cgs.SaleID, cgs.DiscountPercent)
	if err != nil {
		return errors.Wrapf(err, "adding game %v on sale %v", cgs.GameID, cgs.SaleID)
	}

	return nil
}

// ListGameSales returns all sales for specified game
func (s *Storage) ListGameSales(ctx context.Context, gameID int64) ([]GameSale, error) {
	ctx, span := tracer.Start(ctx, "db.game.listgamesales")
	defer span.End()

	gameSales := []GameSale{}

	const q = `select sg.game_id, sg.sale_id, s.name as sale, s.begin_date, s.end_date, sg.discount_percent
	from sales_games sg
	left join sales s on s.id = sg.sale_id
	left join games g on g.id = sg.game_id
	where sg.game_id = $1`
	if err := s.DB.SelectContext(ctx, &gameSales, q, gameID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound{"game sale", gameID}
		}
		return nil, errors.Wrapf(err, "selecting sales of game %q", gameID)
	}

	return gameSales, nil
}
