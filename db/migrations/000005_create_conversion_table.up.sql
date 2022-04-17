CREATE TABLE conversions (
  id BIGSERIAL PRIMARY KEY,
  task_id BIGINT REFERENCES tasks NOT NULL,
  job_id TEXT,
  download_time INTEGER,
  upload_time INTEGER,
  file_size BIGINT,
  result_path TEXT,
  time_spent INTEGER,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE
);
