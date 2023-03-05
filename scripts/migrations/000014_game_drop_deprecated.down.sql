ALTER TABLE games 
    ADD COLUMN developer varchar(150) NULL,
    ADD COLUMN publisher varchar(150) NULL,
    ADD COLUMN genre text[] NULL;
