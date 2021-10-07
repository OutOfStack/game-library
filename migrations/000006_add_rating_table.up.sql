CREATE TABLE Ratings (
	game_id int NOT NULL references Games(id) ON DELETE CASCADE,
	user_id varchar(40) NOT NULL,
	rating smallint NULL,
	PRIMARY KEY(game_id, user_id)
);

CREATE INDEX ratings_user_id_idx on Ratings(user_id); 