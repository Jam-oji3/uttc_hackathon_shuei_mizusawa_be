package mysql

import (
	"context"
	"database/sql"
	"hackathon/model"
	"hackathon/repository"
)

type NotificationsRepository struct{}

var _ repository.NotificationsRepository = (*NotificationsRepository)(nil)

func NewNotificationsRepository() *NotificationsRepository {
	return &NotificationsRepository{}
}

func (r *NotificationsRepository) GetNotifications(ctx context.Context, dbtx repository.DBTX, userId string, limit int) ([]*model.Notification, error) {
	query := `
    SELECT n.type, n.target_id, u.username, u.icon_url, n.created_at
    FROM (
       SELECT 'like' AS type, CAST(l.post_id AS CHAR(128)) AS target_id, l.user_id AS actor_id, l.created_at
       FROM likes l
       JOIN posts p ON l.post_id = p.id
       WHERE p.user_id = ?

       UNION ALL

       SELECT 'repost', CAST(r.post_id AS CHAR(128)), r.user_id, r.created_at
       FROM reposts r
       JOIN posts p ON r.post_id = p.id
       WHERE p.user_id = ?

       UNION ALL

       SELECT 'follow', CAST(f.followed_id AS CHAR(128)), f.follower_id, f.created_at
       FROM follows f
       WHERE f.followed_id = ?
       
       UNION ALL

       SELECT 'reply', CAST(p.reply_to AS CHAR(128)), p.user_id, p.created_at
       FROM posts p
       JOIN posts parent_post ON p.reply_to = parent_post.id
       WHERE parent_post.user_id = ?
    ) AS n
    JOIN users u ON n.actor_id = u.id
    ORDER BY n.created_at DESC
    LIMIT ?
`

	rows, err := dbtx.QueryContext(ctx, query, userId, userId, userId, userId, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*model.Notification

	for rows.Next() {
		var n model.Notification
		var targetID sql.NullString
		if err := rows.Scan(
			&n.Type, &targetID,
			&n.Actor.Username, &n.Actor.IconUrl, &n.CreatedAt,
		); err != nil {
			return nil, err
		}
		if targetID.Valid {
			n.TargetId = &targetID.String
		} else {
			n.TargetId = nil
		}
		notifications = append(notifications, &n)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}
