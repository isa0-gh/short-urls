-- name: CreateUser :one
INSERT INTO
  users (username, email, password_hash)
VALUES
  ($1, $2, $3) RETURNING *;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1
  AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1
  AND deleted_at IS NULL;

-- name: GetUserIDBySessionHash :one
SELECT s.owner_id
FROM sessions AS s
JOIN users AS u ON u.id = s.owner_id
WHERE s.session_hash = $1
  AND s.deleted_at IS NULL
  AND s.expires_at > CURRENT_TIMESTAMP
  AND u.deleted_at IS NULL;

-- name: CreateSession :one
INSERT INTO sessions (session_hash, owner_id, expires_at)
SELECT $1, u.id, $3
FROM users AS u
WHERE u.id = $2
  AND u.deleted_at IS NULL
RETURNING *;

-- name: LogoutSession :exec
UPDATE sessions
SET deleted_at = CURRENT_TIMESTAMP
WHERE session_hash = $1
  AND deleted_at IS NULL;

-- name: LogoutAllSessions :exec
UPDATE sessions
SET deleted_at = CURRENT_TIMESTAMP
WHERE owner_id = $1
  AND deleted_at IS NULL;

-- name: GetDestinationURL :one
SELECT
  destination_url
FROM
  links
WHERE
  short_code = $1 AND deleted_at IS NULL;

-- name: CreateTelemetry :exec
INSERT INTO telemetry (
  link_id,
  referrer_url,
  user_agent,
  ip_address
)
VALUES
  ($1, $2, $3, $4);

-- name: CreateShortLink :one
INSERT INTO links (
  owner_id,short_code,destination_url
)
VALUES
  ($1,$2,$3)
RETURNING *;

-- name: ListShortLinks :many
SELECT
  *
FROM
  links
WHERE
  owner_id = $1
  AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2;

-- name: UpdateDestination :exec
UPDATE
  links
SET destination_url = $1
WHERE id = $2 AND owner_id = $3;

-- name: DeleteShortLink :exec
UPDATE links
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND owner_id = $2;