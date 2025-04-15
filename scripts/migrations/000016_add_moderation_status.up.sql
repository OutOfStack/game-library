-- Add moderation_status column to games table
ALTER TABLE games
    ADD COLUMN moderation_status text NOT NULL DEFAULT '';

-- Set 'ready' status for games fetched from igdb
UPDATE games SET moderation_status = 'ready'
WHERE igdb_id > 0;

-- Set 'check' status for games added by publishers
UPDATE games SET moderation_status = 'check'
WHERE igdb_id = 0;