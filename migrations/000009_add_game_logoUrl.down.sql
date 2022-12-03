CREATE TABLE sales (
	id serial PRIMARY KEY,
	name varchar NOT NULL,
	begin_date date NULL,
	end_date date NULL
);

CREATE TABLE sales_games (
	game_id int NOT NULL references games(id) ON DELETE CASCADE,
	sale_id int NOT NULL references sales(id) ON DELETE CASCADE,
	discount_percent smallint NULL,
	PRIMARY KEY(game_id, sale_id)
);

ALTER TABLE games 
DROP COLUMN rating;

ALTER TABLE games 
ADD COLUMN price numeric(6, 2) NULL;