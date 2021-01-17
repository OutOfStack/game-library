CREATE TABLE Games (
	id serial PRIMARY key,
	name varchar NOT NULL,
	developer varchar NULL,
	releasedate date NULL,
	genre text[] NULL
);