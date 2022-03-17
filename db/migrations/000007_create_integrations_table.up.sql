CREATE TABLE integrations (
  id BIGSERIAL PRIMARY KEY,
  project_id INTEGER REFERENCES projects,
  service service NOT NULL,
  credentials BYTEA,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE
);
