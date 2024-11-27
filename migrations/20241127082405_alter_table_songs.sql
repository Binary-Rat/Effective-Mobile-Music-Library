-- +goose Up
-- +goose StatementBegin
ALTER TABLE songs RENAME COLUMN "text" TO lyrics;
ALTER TABLE songs RENAME COLUMN "group" TO band;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE songs RENAME COLUMN lyrics TO "tetx";
ALTER TABLE songs RENAME COLUMN band TO "group";
-- +goose StatementEnd
