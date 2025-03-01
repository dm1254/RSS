-- name: CreateFeedFollow :one
WITH inserted_feed_follows AS(
INSERT INTO feed_follows(id,created_at,updated_at,user_id,feed_id)
values(
	$1,
	$2,
	$3,
	$4,
	$5

	)
RETURNING *
)
	
SELECT inserted_feed_follows.*, feeds.name AS feeds_name, users.name AS users_name 
FROM inserted_feed_follows
JOIN feeds ON inserted_feed_follows.feed_id = feeds.id
JOIN users ON inserted_feed_follows.user_id = users.id;

-- name: GetFeedFollowForUser :many

SELECT feed_follows.id,
       feed_follows.created_at,
       feed_follows.updated_at,
       feed_follows.user_id,
       feed_follows.feed_id,
       feeds.name AS feed_name,
       users.name AS user_name
FROM feed_follows 
JOIN feeds ON feed_follows.feed_id = feeds.id
JOIN users ON feed_follows.user_id = users.id
WHERE feed_follows.user_id =  $1;

-- name: UnfollowFeed :exec
DELETE FROM feed_follows
WHERE feed_follows.user_id = $1
AND feed_id IN(
	SELECT id FROM feeds WHERE url = $2
);
