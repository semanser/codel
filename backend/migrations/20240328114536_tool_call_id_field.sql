-- +goose Up
-- +goose StatementBegin
ALTER TABLE tasks ADD COLUMN tool_call_id TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tasks DROP COLUMN tool_call_id;
-- +goose StatementEnd
