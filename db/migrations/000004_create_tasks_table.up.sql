CREATE TABLE tasks (
  id BIGSERIAL PRIMARY KEY,
  project_id INTEGER REFERENCES projects,
  task_type task_type NOT NULL,
  status status NOT NULL DEFAULT 'queued',
  url TEXT,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE
);
