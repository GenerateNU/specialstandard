-- Change IEP field from TEXT to TEXT[] (array of strings)
-- This allows students to have multiple IEP goals instead of just one

-- First, drop any existing default value
ALTER TABLE student ALTER COLUMN iep DROP DEFAULT;

-- Then, convert the column type from TEXT to TEXT[]
-- NULL values will remain NULL, and non-null values will become single-element arrays
ALTER TABLE student 
ALTER COLUMN iep TYPE TEXT[] 
USING CASE 
    WHEN iep IS NULL THEN NULL 
    ELSE ARRAY[iep]::TEXT[] 
END;

