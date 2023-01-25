CREATE TYPE task_status AS ENUM ('idle', 'running', 'error');

CREATE TABLE background_tasks (
     name       varchar(100)    PRIMARY KEY,
     run_count  bigint          NOT NULL    DEFAULT 0,
     last_run   timestamptz,
     status     task_status     NOT NULL    DEFAULT 'idle',
     updated_at timestamptz     NOT NULL    DEFAULT now()
);

INSERT INTO background_tasks(name, last_run)
values ('fetch_igdb_games', null);