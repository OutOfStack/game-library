ALTER TABLE games
    ADD COLUMN trailers text[] NOT NULL DEFAULT '{}',
    ADD COLUMN game_mode_ids int[] NOT NULL DEFAULT '{}',
    ADD COLUMN similar_game_ids int[] NOT NULL DEFAULT '{}',
    ADD COLUMN franchise_game_ids int[] NOT NULL DEFAULT '{}';
