-- +goose Up
-- +goose StatementBegin
CREATE TABLE logs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  message TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  flow_id INTEGER REFERENCES flows(id) ON DELETE CASCADE,
  type TEXT NOT NULL -- "input" or "output"
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE logs;
-- +goose StatementEnd
