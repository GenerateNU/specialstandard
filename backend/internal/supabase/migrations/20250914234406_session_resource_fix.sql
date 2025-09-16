DROP TABLE session_resource;

CREATE TABLE session_resource (
    session_id UUID,
    resource_id UUID,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    PRIMARY KEY (session_id, resource_id),
    FOREIGN KEY (session_id) REFERENCES session(id) ON DELETE CASCADE,
    FOREIGN KEY (resource_id) REFERENCES resource(id) ON DELETE CASCADE
);

CREATE TRIGGER update_session_resource_updated_at BEFORE UPDATE ON session_resource
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();