insert into games (name, developer, publisher, release_date, genre, price) values
('Red Dead Redemption 2', 'Rockstar Games', 'Rockstar Games', '2019-12-05', '{"Action", "Western", "Adventure"}', 60),
('Ori and the Will of the Wisps', 'Moon Studios GmbH', 'Xbox Game Studios', '2020-03-11', '{"Action", "Platformer"}', 20),
('The Wolf Among Us', 'Telltale', 'Telltale', '2013-10-11', '{"Adventure", "Episodic", "Detective"}', 20),
('Call of Duty 2', 'Infinity Ward', 'Activision', '2005-10-25', '{"Action", "FPS", "Shooter"}', 10),
('FlatOut 2', 'Bugbear Entertainment', 'Strategy First', '2006-08-01', '{"Racing", "Automobile Sim"}', 10),
('Team Fortress 2', 'Valve', 'Valve', '2007-10-10', '{"Hero shooter", "Multiplayer", "FPS"}', 10),
('Portal', 'Valve', 'Valve', '2007-10-10', '{"Puzzle", "Sci-Fi"}', 10),
('BioShock', '2K', '2K', '2007-08-21', '{"FPS", "Action", "Horror"}', 10),
('S.T.A.L.K.E.R.: Shadow of Chernobyl', 'GSC Game World', 'GSC Game World', '2007-03-20', '{"Post-Apocalyptic", "FPS", "Open world"}', 15),
('Call of Duty 4: Modern Warfare', 'Infinity Ward', 'Activision', '2007-11-12', '{"FPS", "Action", "Shooter"}', 15),
('Grand Theft Auto IV', 'Rockstar North', 'Rockstar Games', '2008-12-02', '{"Open world", "Action"}', 15),
('Mass Effect', 'BioWare', 'Electronic Arts', '2008-05-28', '{"RPG", "Action", "Space"}', 15),
('Fallout 3', 'Bethesda Game Studios', 'Bethesda Softworks', '2008-10-28', '{"Open world", "RPG", "Post-apocalyptic"}', 20),
('Assassin''s Creed', 'Ubisoft Montreal', 'Ubisoft', '2008-04-09', '{"Action", "Adventure", "Stealth"}', 15),
('Crysis', 'Crytek', 'Electronic Arts', '2007-11-13', '{"FPS", "Action", "Shooter"}', 10),
('The Elder Scrolls IV: Oblivion', 'Bethesda Game Studios', 'Bethesda Softworks', '2007-09-11', '{"RPG", "Open World", "Fantasy"}', 15),
('Fallout: New Vegas', 'Obsidian Entertainment', 'Bethesda Softworks', '2010-10-19', '{"Action", "RPG", "Open world"}', 15),
('Metro 2033', '4A Games', 'Deep Silver', '2010-03-16', '{"FPS", "Action", "Post-apocalyptic"}', 10),
('Counter-Strike: Global Offensive', 'Valve', 'Valve', '2012-08-21', '{"FPS", "Shooter", "Multiplayer"}', 10),
('The Walking Dead', 'Telltale', 'Skybound Games', '2012-04-24', '{"Zombies", "Adventure"}', 10),
('Dota 2', 'Valve', 'Valve', '2013-07-13', '{"PvP", "Strategy", "Multiplayer"}', 0),
('BioShock Infinite', 'Irrational Games', '2K', '2013-03-25', '{"FPS", "Action", "Shooter"}', 15),
('Middle-earth: Shadow of Mordor', 'Monolith Productions', 'Warner Bros. Games', '2014-09-30', '{"Open world", "Fantasy", "RPG"}', 15),
('Wolfenstein: The New Order', 'MachineGames', 'Bethesda Softworks', '2014-05-20', '{"FPS", "Action", "Shooter"}', 15),
('The Witcher 3: Wild Hunt', 'CD PROJEKT RED', 'CD PROJEKT RED', '2015-05-18', '{"Open world", "RPG"}', 30),
('Dying Light', 'Techland', 'Techland', '2015-01-26', '{"Zombies", "Horror", "Survival"}', 15),
('DARK SOULS III', 'FromSoftware, Inc.', 'FromSoftware, Inc.', '2016-04-11', '{"Dark fantasy", "RPG", "Souls-like"}', 50),
('DOOM', 'id Software', 'Bethesda Softworks', '2016-05-13', '{"Action", "FPS", "Violent"}', 15),
('Resident Evil 7 Biohazard', 'CAPCOM Co., Ltd.', 'CAPCOM Co., Ltd.', '2017-01-24', '{"Horror", "First-Person"}', 15),
('Subnautica', 'Unknown Worlds Entertainment', 'Unknown Worlds Entertainment', '2018-01-23', '{"Open world", "Survival"}', 15),
('Sekiro: Shadows Die Twice', 'FromSoftware', 'Activision', '2019-03-21', '{"Souls-like", "Action", "Ninja"}', 30),
('STAR WARS Jedi: Fallen Order', 'Respawn Entertainment', 'Electronic Arts', '2019-11-15', '{"Third-person", "RPG", "Action"}', 60),
('DOOM Eternal', 'id Software', 'Bethesda Softworks', '2020-03-20', '{"Action", "FPS", "Violent"}', 40),
('Titanfall 2', 'Respawn Entertainment', 'Electronic Arts', '2016-10-28', '{"FPS", "Multiplayer", "Action"}', 25),
('DEATH STRANDING', 'KOJIMA PRODUCTIONS', '505 Games', '2020-07-14', '{"Open world", "Sci-fi", "Walking simulator"}', 70),
('Valheim', 'Iron Gate AB', 'Coffee Stain Publishing', '2021-02-02', '{"Openworld", "Survival", "Craft"}', 10),
('Resident Evil Village', 'CAPCOM Co., Ltd.', 'CAPCOM Co., Ltd.', '2021-09-07', '{"Horror", "First-Person"}', 60)
on conflict do nothing;

insert into sales (name, begin_date, end_date) values
('Winter sale 2021', '2021-01-01',  '2021-01-15'),
('Summer sale 2021', '2021-06-15',  '2021-06-30'),
('Developer Week', current_date, current_date + 7)
on conflict do nothing;

insert into sales_games(game_id, sale_id, discount_percent)
values (1, 1, 70),
(2, 2, 25)
on conflict do nothing;