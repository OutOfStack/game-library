DROP TABLE IF EXISTS sales_games;

DROP TABLE IF EXISTS sales;

ALTER TABLE games 
DROP COLUMN price;

ALTER TABLE games 
ADD COLUMN rating numeric(5, 2) NULL;