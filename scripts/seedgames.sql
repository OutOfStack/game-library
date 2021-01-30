insert into games (name, developer, releasedate, genre)
values ('Red Dead Redemption 2', 'Rockstar Games', '2019-12-05', '{"Action", "Western", "Adventure"}'),
('Ori and the Will of the Wisps', 'Moon Studios GmbH', '2020-03-11', '{"Action", "Platformer"}'),
('The Wolf Among Us', 'Telltale', '2013-10-11', '{"Adventure", "Episodic", "Detective"}')
on conflict do nothing;