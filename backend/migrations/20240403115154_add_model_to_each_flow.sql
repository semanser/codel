-- +goose Up
-- +goose StatementBegin
ALTER TABLE flows
ADD COLUMN model TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE flows
DROP COLUMN model;
-- +goose StatementEnd
