-- +goose Up
-- +goose StatementBegin
CREATE TABLE game_stats (
    id SERIAL PRIMARY KEY,
    player_name VARCHAR(50) NOT NULL,
    difficulty VARCHAR(20) NOT NULL,
    time_taken_seconds INT NOT NULL,
    solved_at TIMESTAMPTZ DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS game_stats;
-- +goose StatementEnd