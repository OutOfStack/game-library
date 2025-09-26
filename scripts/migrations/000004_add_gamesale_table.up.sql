ALTER TABLE sales
DROP COLUMN discount_percent,
DROP COLUMN game_id;

CREATE TABLE sales_games (
	game_id int NOT NULL references games(id) ON DELETE CASCADE,
	sale_id int NOT NULL references sales(id) ON DELETE CASCADE,
	discount_percent smallint NULL,
	PRIMARY KEY(game_id, sale_id)
)