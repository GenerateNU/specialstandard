-- Create SessionParent table
CREATE TABLE IF NOT EXISTS public.session_parent (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    start_date date NOT NULL,
    end_date date NOT NULL,
    therapist_id uuid NOT NULL,
    days smallint[], -- stores days of week, e.g., [1,3]
    every_n_weeks int,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT session_parent_pkey PRIMARY KEY (id),
    CONSTRAINT session_parent_therapist_id_fkey FOREIGN KEY (therapist_id) REFERENCES therapist (id) ON DELETE RESTRICT,
    CONSTRAINT session_parent_check CHECK (end_date >= start_date)
) TABLESPACE pg_default;

CREATE INDEX IF NOT EXISTS idx_session_parent_therapist ON public.session_parent USING btree (therapist_id);

-- Add session_parent_id to session table
ALTER TABLE public.session
ADD COLUMN IF NOT EXISTS session_parent_id uuid NOT NULL;

-- Add foreign key constraint only if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'session_parent_fkey'
          AND table_name = 'session'
    ) THEN
        ALTER TABLE public.session
        ADD CONSTRAINT session_parent_fkey FOREIGN KEY (session_parent_id) REFERENCES session_parent(id) ON DELETE CASCADE;
    END IF;
END $$;

-- Remove therapist_id from session table
ALTER TABLE public.session
DROP CONSTRAINT IF EXISTS session_therapist_id_fkey;

ALTER TABLE public.session
DROP COLUMN IF EXISTS therapist_id;

-- Ensure triggers exist on session table
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM pg_trigger 
        WHERE tgname = 'table_accessed_session'
    ) THEN
        CREATE TRIGGER table_accessed_session
        AFTER INSERT OR DELETE OR UPDATE ON session
        FOR EACH ROW
        EXECUTE FUNCTION log_table_access();
    END IF;
    
    IF NOT EXISTS (
        SELECT 1 
        FROM pg_trigger 
        WHERE tgname = 'update_session_updated_at'
    ) THEN
        CREATE TRIGGER update_session_updated_at
        BEFORE UPDATE ON session
        FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;

-- Add indexes for faster queries
CREATE INDEX IF NOT EXISTS idx_session_datetime ON public.session USING btree (start_datetime, end_datetime);
CREATE INDEX IF NOT EXISTS idx_session_parent_id ON public.session USING btree (session_parent_id);



