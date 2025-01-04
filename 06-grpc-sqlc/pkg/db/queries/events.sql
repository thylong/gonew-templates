-- queries.sql

-- Get events with pagination and optional search
-- name: GetEvents :many
SELECT * FROM events
WHERE (LOWER(name) LIKE LOWER('%' || $1 || '%') OR $1 = '')
LIMIT $2 OFFSET $3;

-- Get total count for pagination
-- name: GetEventsCount :one
SELECT COUNT(*) FROM events
WHERE (LOWER(name) LIKE LOWER('%' || $1 || '%') OR $1 = '');

-- Get a single event by ID
-- name: GetEventByID :one
SELECT * FROM events WHERE event_id = $1;

-- Insert a new event
-- name: InsertEvent :one
INSERT INTO events (
    event_privacy, name, type, department, regions, tags, start_at
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- Delete an event by ID
-- name: DeleteEvent :exec
DELETE FROM events WHERE event_id = $1;
