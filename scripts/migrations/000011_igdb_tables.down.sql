DROP TABLE platforms;

DROP TABLE companies;

DROP TABLE genres;

ALTER TABLE games
    DROP COLUMN summary,
    DROP COLUMN genres,
    DROP COLUMN platforms,
    DROP COLUMN screenshots,
    DROP COLUMN developers,
    DROP COLUMN publishers,
    DROP COLUMN slug,
    DROP COLUMN websites,
    DROP COLUMN igdb_rating,
    DROP COLUMN igdb_id;

ALTER TABLE games ALTER COLUMN logo_url DROP DEFAULT;
UPDATE games set logo_url = null WHERE logo_url = '';
ALTER TABLE games ALTER COLUMN logo_url DROP NOT NULL;

ALTER TABLE games ALTER COLUMN rating DROP DEFAULT;
UPDATE games set rating = null WHERE rating = 0.0;
ALTER TABLE games ALTER COLUMN logo_url DROP NOT NULL;
