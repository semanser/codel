-- +goose Up
-- +goose StatementBegin
ALTER TABLE flows ALTER COLUMN created_at TYPE timestamptz USING created_at AT TIME ZONE 'UTC';
ALTER TABLE flows ALTER COLUMN updated_at TYPE timestamptz USING updated_at AT TIME ZONE 'UTC';
ALTER TABLE tasks ALTER COLUMN created_at TYPE timestamptz USING created_at AT TIME ZONE 'UTC';
ALTER TABLE tasks ALTER COLUMN updated_at TYPE timestamptz USING updated_at AT TIME ZONE 'UTC';
ALTER TABLE logs ALTER COLUMN created_at TYPE timestamptz USING created_at AT TIME ZONE 'UTC';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE flows ALTER COLUMN created_at TYPE timestamp USING created_at AT TIME ZONE 'UTC';
ALTER TABLE flows ALTER COLUMN updated_at TYPE timestamp USING updated_at AT TIME ZONE 'UTC';
ALTER TABLE tasks ALTER COLUMN created_at TYPE timestamp USING created_at AT TIME ZONE 'UTC';
ALTER TABLE tasks ALTER COLUMN updated_at TYPE timestamptz USING updated_at AT TIME ZONE 'UTC';
ALTER TABLE logs ALTER COLUMN created_at TYPE timestamp USING created_at AT TIME ZONE 'UTC';
-- +goose StatementEnd
