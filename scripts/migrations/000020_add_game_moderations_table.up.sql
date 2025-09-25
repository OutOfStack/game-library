CREATE TABLE IF NOT EXISTS game_moderation (
    id              int         GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    game_id         int         NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    result_status   text        NOT NULL DEFAULT '',
    -- details of moderation result
    details         text        NOT NULL DEFAULT '',
    -- error of moderation task
    error           text,
    -- payload of game data sent to moderation
    game_data       jsonb       NOT NULL DEFAULT '{}'::jsonb,
    created_at      timestamptz,
    updated_at      timestamptz
);

CREATE INDEX IF NOT EXISTS game_moderation_created_idx
    ON game_moderation (game_id, id DESC);

ALTER TABLE games
    ADD COLUMN IF NOT EXISTS moderation_id int REFERENCES game_moderation(id) ON DELETE SET NULL;

-- Index for GetGamesByPublisherID query optimization
CREATE INDEX IF NOT EXISTS games_publishers_gin_idx
    ON games USING GIN (publishers);

