-- Change IEP field from TEXT to TEXT[] (array of strings)
-- This allows students to have multiple IEP goals instead of just one

-- First, update existing data to convert single string values to arrays
-- NULL values will remain NULL, and non-null values will become single-element arrays
ALTER TABLE student 
ALTER COLUMN iep TYPE TEXT[] 
USING CASE 
    WHEN iep IS NULL THEN NULL 
    ELSE ARRAY[iep]::TEXT[] 
END;

-- Update the default to be NULL (empty array would be {} but NULL is better for optional fields)
ALTER TABLE student ALTER COLUMN iep SET DEFAULT NULL;

