ALTER TABLE companies
    DROP COLUMN created_at;

ALTER TABLE genres
  DROP COLUMN created_at;

ALTER TABLE games
  DROP COLUMN created_at,
  DROP COLUMN updated_at;

ALTER TABLE ratings
  DROP COLUMN created_at,
  DROP COLUMN updated_at;
