-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS card (
    user_id   UUID  NOT NULL,
    data_id   TEXT  NOT NULL,
	card_data TEXT  DEFAULT NULL,
    PRIMARY KEY (user_id, data_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS card;
-- +goose StatementEnd
