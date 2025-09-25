ALTER TABLE games
    DROP COLUMN IF EXISTS moderation_id;

DROP INDEX IF EXISTS games_publishers_gin_idx;

DROP TABLE IF EXISTS game_moderation;