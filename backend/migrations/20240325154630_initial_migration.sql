-- +goose Up
-- +goose StatementBegin
CREATE TABLE containers (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT,
  local_id TEXT,
  image TEXT,
  status TEXT DEFAULT 'starting'
);

CREATE TABLE flows (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  name TEXT,
  status TEXT,
  container_id INTEGER,
  FOREIGN KEY (container_id) REFERENCES containers (id)
);

CREATE TABLE tasks (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  type TEXT,
  status TEXT,
  args TEXT DEFAULT '{}',
  results TEXT DEFAULT '{}',
  message TEXT,
  flow_id INTEGER,
  FOREIGN KEY (flow_id) REFERENCES flows (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tasks;
DROP TABLE flows;
DROP TABLE containers;
-- +goose StatementEnd
```
