CREATE TABLE tasks (
  id BIGSERIAL PRIMARY KEY,
  project_id INTEGER REFERENCES projects,
  task_type task_type NOT NULL,
  status status NOT NULL DEFAULT 'queued',
  download_time INTEGER,
  upload_time INTEGER,
  file_size INTEGER,
  file_path TEXT,
  result_path TEXT,
  time_spent INTEGER,
  url TEXT,
  error TEXT,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);
