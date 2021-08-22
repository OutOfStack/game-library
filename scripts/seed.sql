insert into games (name, developer, release_date, genre, price)
values ('Red Dead Redemption 2', 'Rockstar Games', '2019-12-05', '{"Action", "Western", "Adventure"}', 60),
('Ori and the Will of the Wisps', 'Moon Studios GmbH', '2020-03-11', '{"Action", "Platformer"}', 20),
('The Wolf Among Us', 'Telltale', '2013-10-11', '{"Adventure", "Episodic", "Detective"}', 15)
on conflict do nothing;

insert into sales (name, begin_date, end_date)
values ('Winter sale 2020', '2020-12-22',  '2021-01-10'),
('Developer Week', current_date, current_date + 7)
on conflict do nothing;

insert into sales_games(game_id, sale_id, discount_percent)
values (1, 1, 70),
(2, 2, 25)
on conflict do nothing;