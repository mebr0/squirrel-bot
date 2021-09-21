CREATE TABLE IF NOT EXISTS games(
    id VARCHAR PRIMARY KEY,
    score_1 SMALLINT NOT NULL,
    score_2 SMALLINT NOT NULL,
    rounds SMALLINT NOT NULL
);

CREATE TABLE IF NOT EXISTS players(
    id BIGINT PRIMARY KEY,
    username VARCHAR,
    first_name VARCHAR,
    last_name VARCHAR
);

CREATE TABLE IF NOT EXISTS game_players(
    id BIGSERIAL PRIMARY KEY,
    game_id VARCHAR NOT NULL,
    player_id BIGINT NOT NULL,
    team SMALLINT NOT NULL,
    CONSTRAINT mtm_games_to_players FOREIGN KEY(game_id) REFERENCES games(id) ON DELETE SET NULL,
    CONSTRAINT mtm_players_to_games FOREIGN KEY(player_id) REFERENCES players(id) ON DELETE SET NULL
);