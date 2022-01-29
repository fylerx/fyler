CREATE TABLE conversions (
  id BIGSERIAL PRIMARY KEY,
  task_id INTEGER REFERENCES tasks,
  download_time INTEGER,
  upload_time INTEGER,
  file_size INTEGER,
  file_path TEXT,
  result_path TEXT,
  time_spent INTEGER,
  error TEXT,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE
);
