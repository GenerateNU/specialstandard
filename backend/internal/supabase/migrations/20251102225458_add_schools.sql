CREATE TABLE IF NOT EXISTS district (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS school (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  district_id INTEGER REFERENCES district(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add placeholders
INSERT INTO district (name) VALUES ('Generate School District') ON CONFLICT DO NOTHING;
INSERT INTO school (name, district_id) VALUES ('Generate Elementary', 1) ON CONFLICT DO NOTHING;

-- Backfill
ALTER TABLE student
ADD COLUMN school_id INTEGER NOT NULL DEFAULT 1 REFERENCES school(id);

ALTER TABLE therapist 
ADD COLUMN schools INTEGER[],
ADD COLUMN district_id INTEGER DEFAULT 1 REFERENCES district(id);

-- Remove default (used when backfilling)
ALTER TABLE student
ALTER COLUMN school_id DROP DEFAULT;

ALTER TABLE therapist
ALTER COLUMN district_id DROP DEFAULT;

CREATE INDEX idx_student_school ON student(school_id);
CREATE INDEX idx_therapist_schools ON therapist USING GIN (schools);
CREATE INDEX idx_therapist_district ON therapist(district_id);

