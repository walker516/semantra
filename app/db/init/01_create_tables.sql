CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE features (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  type TEXT NOT NULL CHECK (type IN ('point', 'line', 'polygon')),
  geom GEOMETRY(GEOMETRY, 4326) NOT NULL,
  properties JSONB DEFAULT '{}'::jsonb,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);
