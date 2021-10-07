CREATE TABLE Games (
	id serial PRIMARY KEY,
	name varchar NOT NULL,
	developer varchar NOT NULL,
	release_date date NULL,
	genre text[] NULL
);