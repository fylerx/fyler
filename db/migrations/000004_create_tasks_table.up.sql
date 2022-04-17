CREATE TABLE tasks (
  id BIGSERIAL PRIMARY KEY,
  project_id BIGINT REFERENCES projects NOT NULL,
  task_type task_type NOT NULL,
  file_path TEXT,
  status status NOT NULL DEFAULT 'queued',
  error TEXT,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE
);
