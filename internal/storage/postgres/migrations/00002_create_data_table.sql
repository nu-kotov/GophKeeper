-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS private_data (
    user_id        UUID  NOT NULL,
    data_name      TEXT  NOT NULL
    login          TEXT  DEFAULT NULL,
    password       TEXT  DEFAULT NULL,
    text_data      TEXT  DEFAULT NULL,
	binary_content BYTEA DEFAULT NULL,
	card           TEXT  DEFAULT NULL,
    PRIMARY KEY (user_id, data_name)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS private_data;
-- +goose StatementEnd
