-- Migration to normalize grade fields from string to integer
-- Grade values: -1 (graduated), 0, 1, 2, 3, ..., 12
-- For resources, grade_level cannot be -1

BEGIN;

-- Convert student grades: VARCHAR -> INTEGER using USING clause
ALTER TABLE student 
ALTER COLUMN grade TYPE INTEGER USING (
  CASE 
    WHEN grade IS NULL THEN NULL
    WHEN grade ~ '[0-9]+' AND 
         (SUBSTRING(grade FROM '[0-9]+'))::INTEGER BETWEEN 0 AND 12 
         THEN (SUBSTRING(grade FROM '[0-9]+'))::INTEGER
    WHEN LOWER(grade) IN ('kindergarten', 'k') THEN 0
    WHEN LOWER(grade) IN ('graduated', 'graduate') THEN -1
    ELSE NULL  -- Invalid data becomes NULL rather than causing failure
  END
);

-- Add constraint for student grades
ALTER TABLE student ADD CONSTRAINT check_student_grade 
    CHECK (grade IS NULL OR grade = -1 OR (grade >= 0 AND grade <= 12));

-- Convert resource grade_level: VARCHAR -> INTEGER using USING clause
ALTER TABLE resource 
ALTER COLUMN grade_level TYPE INTEGER USING (
  CASE 
    WHEN grade_level IS NULL THEN NULL
    -- Extract first number from any string: "5th Grade" -> 5, "Grade 5" -> 5, "5" -> 5
    WHEN grade_level ~ '[0-9]+' AND 
         (SUBSTRING(grade_level FROM '[0-9]+'))::INTEGER BETWEEN 0 AND 12 
         THEN (SUBSTRING(grade_level FROM '[0-9]+'))::INTEGER
    -- Handle special grade names
    WHEN LOWER(grade_level) IN ('kindergarten', 'k') THEN 0
    ELSE NULL  -- Invalid data becomes NULL (no -1 for resources)
  END
);

-- Add constraint for resource grade_level (no -1 allowed)
ALTER TABLE resource ADD CONSTRAINT check_resource_grade_level 
    CHECK (grade_level IS NULL OR (grade_level >= 0 AND grade_level <= 12));

-- Update indexes
DROP INDEX IF EXISTS idx_resource_grade;
CREATE INDEX idx_resource_grade_level ON resource(grade_level);

COMMIT;