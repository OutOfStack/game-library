insert into games (name, developer, release_date, genre, price)
values ('Red Dead Redemption 2', 'Rockstar Games', '2019-12-05', '{"Action", "Western", "Adventure"}', 60),
('Ori and the Will of the Wisps', 'Moon Studios GmbH', '2020-03-11', '{"Action", "Platformer"}', 20),
('The Wolf Among Us', 'Telltale', '2013-10-11', '{"Adventure", "Episodic", "Detective"}', 15)
on conflict do nothing;

insert into sales (name, game_id, begin_date, end_date, discount_percent)
values ('Winter sale 2020', 1, '2020-12-22',  '2021-01-10', 70),
('Developer Week', 2, '2021-02-11', '2021-02-17', 25)
on conflict do nothing;