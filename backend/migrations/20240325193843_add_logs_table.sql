-- +goose Up
-- +goose StatementBegin
CREATE TABLE logs (
  id SERIAL PRIMARY KEY,
  message TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  flow_id bigint REFERENCES flows(id) ON DELETE CASCADE,
  type TEXT NOT NULL -- "input" or "output"
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE logs;
-- +goose StatementEnd
