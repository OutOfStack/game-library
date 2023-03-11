ALTER TABLE companies ALTER COLUMN igdb_id SET DEFAULT 0;
UPDATE companies set igdb_id = 0 WHERE igdb_id IS NULL;
ALTER TABLE companies ALTER COLUMN igdb_id SET NOT NULL;
