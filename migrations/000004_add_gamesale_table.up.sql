ALTER TABLE Sales
DROP COLUMN discount_percent,
DROP COLUMN game_id;

CREATE TABLE Sales_games (
	game_id int NOT NULL references Games(id) ON DELETE CASCADE,
	sale_id int NOT NULL references Sales(id) ON DELETE CASCADE,
	discount_percent smallint NULL,
	PRIMARY KEY(game_id, sale_id)
)