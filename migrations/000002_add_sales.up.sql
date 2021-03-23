CREATE TABLE Sales (
	id serial PRIMARY key,
	name varchar NOT NULL,
	game_id int NOT NULL references Games(id) ON DELETE CASCADE,
	begin_date date NULL,
	end_date date NULL,
	discount_percent smallint NULL
);