CREATE TABLE storages (
  project_id BIGINT REFERENCES projects NOT NULL PRIMARY KEY,
  access_key_id BYTEA NOT NULL,
  secret_access_key BYTEA NOT NULL,
  bucket TEXT NOT NULL,
  endpoint TEXT,
  region TEXT,
  disable_ssl BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE
);
