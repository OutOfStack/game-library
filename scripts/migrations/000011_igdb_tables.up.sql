ALTER TABLE games
    ADD COLUMN summary text,
    ADD COLUMN genres int[],
    ADD COLUMN platforms int[],
    ADD COLUMN screenshots text[],
    ADD COLUMN developers int[],
    ADD COLUMN publishers int[],
    ADD COLUMN websites text[],
    ADD COLUMN slug varchar(50),
    ADD COLUMN igdb_rating numeric(5, 2),
    ADD COLUMN igdb_id bigint;

ALTER TABLE games ALTER COLUMN logo_url SET DEFAULT '';
UPDATE games set logo_url = '' WHERE logo_url IS NULL;
--nolint:set-not-null
ALTER TABLE games ALTER COLUMN logo_url SET NOT NULL;

ALTER TABLE games ALTER COLUMN rating SET DEFAULT 0.0;
UPDATE games set rating = 0.0 WHERE rating IS NULL;
--nolint:set-not-null
ALTER TABLE games ALTER COLUMN logo_url SET NOT NULL;

ALTER TABLE games ALTER COLUMN summary SET DEFAULT '';
UPDATE games set summary = '';
--nolint:set-not-null
ALTER TABLE games ALTER COLUMN summary SET NOT NULL;

ALTER TABLE games ALTER COLUMN slug SET DEFAULT '';
UPDATE games set slug = '';
--nolint:set-not-null
ALTER TABLE games ALTER COLUMN slug SET NOT NULL;

ALTER TABLE games ALTER COLUMN igdb_rating SET DEFAULT 0.0;
UPDATE games set igdb_rating = 0.0;
--nolint:set-not-null
ALTER TABLE games ALTER COLUMN igdb_rating SET NOT NULL;

ALTER TABLE games ALTER COLUMN igdb_id SET DEFAULT 0;
UPDATE games set igdb_id = 0;
--nolint:set-not-null
ALTER TABLE games ALTER COLUMN igdb_id SET NOT NULL;

CREATE TABLE genres (
    id             int             GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name           varchar(100)    NOT NULL    DEFAULT '',
    igdb_id        bigint          NOT NULL    DEFAULT 0   UNIQUE
);

CREATE TABLE companies (
    id             int             GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name           varchar(100)    NOT NULL    DEFAULT '',
    igdb_id        bigint          NOT NULL    DEFAULT 0   UNIQUE
);

CREATE TABLE platforms (
    id             int             GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name           varchar(100)    NOT NULL    DEFAULT '',
    abbreviation   varchar(15)     NOT NULL    DEFAULT '',
    igdb_id        bigint          NOT NULL    DEFAULT 0   UNIQUE
);
    

INSERT INTO platforms(name, abbreviation, igdb_id) values 
('Linux', 'Linux', 3),
('PC (Microsoft Windows)', 'PC', 6),
('Mac', 'Mac', 14),
('PlayStation', 'PS1', 7),
('PlayStation 2', 'PS2', 8),
('PlayStation 3', 'PS3', 9),
('PlayStation 4', 'PS4', 48),
('PlayStation 5', 'PS5', 167),
('Xbox', 'XBOX', 11),
('Xbox 360', 'X360', 12),
('Xbox One', 'XONE', 49),
('Xbox Series X|S', 'Series X', 169);
