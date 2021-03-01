CREATE TABLE Games (
	id serial PRIMARY key,
	name varchar NOT NULL,
	developer varchar NULL,
	release_date date NULL,
	genre text[] NULL
);