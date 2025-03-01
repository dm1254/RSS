-- name: CreateFeed :one 
INSERT INTO feeds(id,created_at, updated_at,last_fetched_at,name,url,user_id)
values(
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7
)
RETURNING *;

-- name: GetFeeds :many
SELECT feeds.name AS feed_name,feeds.url,users.name AS user_name FROM feeds
JOIN users
ON feeds.user_id = users.id;


-- name: GetFeed :one
SELECT * FROM feeds
WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = NOW(),
updated_at = NOW()
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT id,name,url,user_id FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;
