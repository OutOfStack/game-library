ALTER TABLE game_moderation
    ADD COLUMN attempts int NOT NULL DEFAULT 0;

-- no new moderation data record is created on every attempt, so this field would be just overwritten on every service fail,
-- therefore, field is redundant - look into logs for fails
ALTER TABLE game_moderation
    DROP COLUMN error;

CREATE INDEX IF NOT EXISTS idx_game_moderation_status_id ON game_moderation(status, id);

CREATE INDEX IF NOT EXISTS idx_games_moderation_id ON games(moderation_id);