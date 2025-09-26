ALTER TABLE sales 
ADD COLUMN discount_percent smallint NULL,
ADD COLUMN game_id int NULL references games(id) ON DELETE CASCADE;

DROP TABLE sales_games;