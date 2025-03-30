-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING feed_follows.*,
    (SELECT name FROM feeds WHERE feeds.id = feed_follows.feed_id) AS feed_name,
    (SELECT name FROM users WHERE users.id = feed_follows.user_id) AS user_name;

-- name: GetFeedByUrl :one
SELECT * FROM feeds WHERE url = $1;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*,
    (SELECT name FROM feeds WHERE feeds.id = feed_follows.feed_id) as feed_name,
    (SELECT name FROM users WHERE users.id = feed_follows.user_id) as user_name
FROM feed_follows
WHERE feed_follows.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE feed_follows.user_id = $1 AND feed_id = (SELECT id FROM feeds WHERE url = $2);