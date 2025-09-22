-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS text_data (
    user_id   UUID  NOT NULL,
    data_id   TEXT  NOT NULL
    text_data TEXT  DEFAULT NULL,
    PRIMARY KEY (user_id, data_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS text_data;
-- +goose StatementEnd
