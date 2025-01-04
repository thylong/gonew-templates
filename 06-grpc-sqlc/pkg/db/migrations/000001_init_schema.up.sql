CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE events (
    event_id UUID NOT NULL DEFAULT (uuid_generate_v4()),
    event_privacy INT NOT NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    department TEXT NOT NULL,
    regions TEXT NOT NULL,
    tags TEXT NOT NULL,
    start_at TIMESTAMP NOT NULL
);
