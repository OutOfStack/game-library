DROP INDEX IF EXISTS idx_games_trending_index;

ALTER TABLE games
    DROP COLUMN IF EXISTS trending_index,
    ALTER COLUMN release_date DROP NOT NULL;