ALTER TABLE public.resource
ADD COLUMN week integer NOT NULL DEFAULT 1 CHECK (week >= 1 AND week <= 4);

-- 2. Drop the old 'date' column
ALTER TABLE public.resource
DROP COLUMN date;