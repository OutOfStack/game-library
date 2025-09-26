ALTER TABLE games
    DROP COLUMN IF EXISTS moderation_id;

DROP TABLE IF EXISTS game_moderation;