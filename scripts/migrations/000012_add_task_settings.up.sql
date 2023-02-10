ALTER TABLE background_tasks
    ADD COLUMN settings jsonb;

ALTER TABLE background_tasks ALTER COLUMN settings SET DEFAULT '{}';
UPDATE background_tasks set settings = '{}';
--nolint:set-not-null
ALTER TABLE background_tasks ALTER COLUMN settings SET NOT NULL;