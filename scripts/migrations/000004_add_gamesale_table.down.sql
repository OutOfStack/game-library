ALTER TABLE sales 
ADD COLUMN discount_percent smallint NULL,
ADD COLUMN game_id int NULL references Games(id) ON DELETE CASCADE;

DROP TABLE Sales_games;