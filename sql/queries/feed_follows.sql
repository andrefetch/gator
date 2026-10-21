-- name: CreateFeedFollow :one
WITH inserted AS (
    INSERT INTO feed_follows (
        id,
        created_at,
        updated_at,
        user_id,
        feed_id
    )
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT
    inserted.id,
    inserted.created_at,
    inserted.updated_at,
    inserted.user_id,
    inserted.feed_id,
    users.name AS user_name,
    feeds.name AS feed_name
FROM inserted
JOIN users ON inserted.user_id = users.id
JOIN feeds ON inserted.feed_id = feeds.id;