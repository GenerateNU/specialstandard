CREATE TYPE access AS ENUM ('GET', 'DELETE', 'UPDATE', 'CREATE');

CREATE TABLE Audit_Log (
    access_type access,
    previous_data json,
    new_data json,
    time_of_operation TIMESTAMPTZ DEFAULT now(),
    user NULL,
);