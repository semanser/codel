-- +goose Up
-- +goose StatementBegin
ALTER TABLE flows
ADD COLUMN model_provider TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE flows
DROP COLUMN model_provider;
-- +goose StatementEnd
