ALTER TABLE game_moderation
    ADD COLUMN error text;

ALTER TABLE game_moderation
    DROP COLUMN attempts;

DROP INDEX IF EXISTS idx_game_moderation_status_id;