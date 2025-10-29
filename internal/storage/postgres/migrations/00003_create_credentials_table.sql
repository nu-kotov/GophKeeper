-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS credentials (
    user_id  UUID  NOT NULL,
    data_id  TEXT  NOT NULL,
    login    TEXT  DEFAULT NULL,
    password TEXT  DEFAULT NULL,
    PRIMARY KEY (user_id, data_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS credentials;
-- +goose StatementEnd
