-- +goose Up
-- +goose StatementBegin
CREATE TABLE containers (
  id BIGSERIAL PRIMARY KEY,
  name text,
  image text,
  status text DEFAULT 'starting'::text
);

CREATE TABLE flows (
  id BIGSERIAL PRIMARY KEY,
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone,
  name text,
  status text,
  container_id bigint REFERENCES containers(id)
);

CREATE TABLE tasks (
  id BIGSERIAL PRIMARY KEY,
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone,
  type text,
  status text,
  args jsonb DEFAULT '{}'::jsonb,
  results text DEFAULT '{}'::jsonb,
  flow_id bigint,
  message text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tasks;
DROP TABLE flows;
DROP TABLE containers;
-- +goose StatementEnd
