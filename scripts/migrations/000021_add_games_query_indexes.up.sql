-- GIN index for genres array filtering (genre_id = ANY(genres))
CREATE INDEX IF NOT EXISTS games_genres_gin_idx
    ON games USING GIN (genres);

-- GIN index for developers array filtering (developer_id = ANY(developers))
CREATE INDEX IF NOT EXISTS games_developers_gin_idx
    ON games USING GIN (developers);

-- GIN index for publishers array filtering (publisher_id = ANY(publishers))
CREATE INDEX IF NOT EXISTS games_publishers_gin_idx
    ON games USING GIN (publishers);

-- Composite index for moderation status filtering with trending index ordering
-- This covers the most common query pattern: WHERE moderation_status = 'ready' ORDER BY trending_index DESC
CREATE INDEX IF NOT EXISTS games_moderation_trending_idx
    ON games (moderation_status, trending_index DESC)
    WHERE moderation_status = 'ready';

-- rename game status
ALTER TABLE games
    ALTER COLUMN moderation_status SET DEFAULT 'pending';

UPDATE games
SET moderation_status = 'pending' WHERE moderation_status = 'check';

ALTER TABLE companies
    ADD CONSTRAINT unique_name UNIQUE (name);