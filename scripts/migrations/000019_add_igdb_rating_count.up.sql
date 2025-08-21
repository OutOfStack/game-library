ALTER TABLE games
    ADD COLUMN igdb_rating_count int NOT NULL DEFAULT 0,
    ALTER COLUMN rating SET NOT NULL; -- fix of migration 000011_igdb_tables.up

INSERT INTO background_tasks(name, last_run)
VALUES ('update_game_info', null);