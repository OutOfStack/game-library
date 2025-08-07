ALTER TABLE games
    ADD COLUMN trending_index numeric(10, 4) NOT NULL DEFAULT 0.0,
    ALTER COLUMN release_date SET NOT NULL;

UPDATE games 
    SET trending_index = (EXTRACT(year FROM release_date)/2.5 + COALESCE(igdb_rating, 0) + COALESCE(rating, 0)/2.0);

CREATE INDEX idx_games_trending_index ON games (trending_index DESC);