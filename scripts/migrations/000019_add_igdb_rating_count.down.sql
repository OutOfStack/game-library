ALTER TABLE games
    DROP COLUMN igdb_rating_count,
    ALTER COLUMN rating DROP NOT NULL; -- fix of migration 000011_igdb_tables.down

DELETE FROM background_tasks
WHERE name = 'update_game_info';