ALTER TABLE companies
    DROP CONSTRAINT unique_name;

UPDATE games
    SET moderation_status = 'check' WHERE moderation_status = 'pending';

ALTER TABLE games
    ALTER COLUMN moderation_status SET DEFAULT 'pending';

DROP INDEX IF EXISTS games_moderation_trending_idx;
DROP INDEX IF EXISTS games_developers_gin_idx;
DROP INDEX IF EXISTS games_genres_gin_idx;
DROP INDEX IF EXISTS games_publishers_gin_idx;