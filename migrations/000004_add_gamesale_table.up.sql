ALTER TABLE Sales
DROP COLUMN discount_percent,
DROP COLUMN game_id;

CREATE TABLE Sales_games (
	game_id int not null references Games(id) on delete cascade,
	sale_id int not null references Sales(id) on delete cascade,
	discount_percent smallint null,
	PRIMARY KEY(game_id, sale_id)
)