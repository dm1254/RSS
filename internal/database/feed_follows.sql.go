// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: feed_follows.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFeedFollow = `-- name: CreateFeedFollow :one
WITH inserted_feed_follows AS(
INSERT INTO feed_follows(id,created_at,updated_at,user_id,feed_id)
values(
	$1,
	$2,
	$3,
	$4,
	$5

	)
RETURNING id, created_at, updated_at, user_id, feed_id
)
	
SELECT inserted_feed_follows.id, inserted_feed_follows.created_at, inserted_feed_follows.updated_at, inserted_feed_follows.user_id, inserted_feed_follows.feed_id, feeds.name AS feeds_name, users.name AS users_name 
FROM inserted_feed_follows
JOIN feeds ON inserted_feed_follows.feed_id = feeds.id
JOIN users ON inserted_feed_follows.user_id = users.id
`

type CreateFeedFollowParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
}

type CreateFeedFollowRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
	FeedsName string
	UsersName string
}

func (q *Queries) CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) (CreateFeedFollowRow, error) {
	row := q.db.QueryRowContext(ctx, createFeedFollow,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.FeedID,
	)
	var i CreateFeedFollowRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
		&i.FeedsName,
		&i.UsersName,
	)
	return i, err
}

const getFeedFollowForUser = `-- name: GetFeedFollowForUser :many

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
WHERE feed_follows.user_id =  $1
`

type GetFeedFollowForUserRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
	FeedName  string
	UserName  string
}

func (q *Queries) GetFeedFollowForUser(ctx context.Context, userID uuid.UUID) ([]GetFeedFollowForUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedFollowForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedFollowForUserRow
	for rows.Next() {
		var i GetFeedFollowForUserRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.FeedID,
			&i.FeedName,
			&i.UserName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const unfollowFeed = `-- name: UnfollowFeed :exec
DELETE FROM feed_follows
WHERE feed_follows.user_id = $1
AND feed_id IN(
	SELECT id FROM feeds WHERE url = $2
)
`

type UnfollowFeedParams struct {
	UserID uuid.UUID
	Url    string
}

func (q *Queries) UnfollowFeed(ctx context.Context, arg UnfollowFeedParams) error {
	_, err := q.db.ExecContext(ctx, unfollowFeed, arg.UserID, arg.Url)
	return err
}
